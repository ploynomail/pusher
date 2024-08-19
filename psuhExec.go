package pusher

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus/push"
)

const pushInstanceLabel string = "instance"

type Pusher struct {
	PushConfig
	httpClient push.HTTPDoer
	wg         *sync.WaitGroup
	ticker     *time.Ticker
	sig        chan os.Signal
	exit       chan struct{}
}

func NewPusher(config PushConfig, exitCh chan struct{}) *Pusher {
	if config.Interval <= 0 {
		config.Interval = 10
	}
	if config.InstanceLabel == "" {
		config.InstanceLabel = "localhost"
	}

	if config.PushGatewayURL == "" {
		config.PushGatewayURL = "http://localhost:9091"
	}
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	return &Pusher{
		PushConfig: config,
		httpClient: &http.Client{},
		wg:         &sync.WaitGroup{},
		ticker:     time.NewTicker(time.Duration(config.Interval) * time.Second),
		sig:        interrupt,
		exit:       exitCh,
	}
}

func (p *Pusher) WithHTTPClient(client push.HTTPDoer) *Pusher {
	p.httpClient = client
	return p
}

func (p *Pusher) ExecPush() {
	for {
		select {
		case <-p.ticker.C:
			for _, target := range p.TargetExporter {
				p.wg.Add(1)
				exporterCollector := &exporterCollector{
					Client:    &http.Client{},
					Url:       target.ExporterURL,
					Collector: target.Collector,
				}
				go func(target TargetExporter) {
					defer p.wg.Done()
					if target.ExporterURL != "" {
						pusher := push.New(p.PushConfig.PushGatewayURL, p.PushConfig.InstanceLabel).
							Grouping(pushInstanceLabel, target.JobName).
							Gatherer(exporterCollector)
						if p.httpClient != nil {
							pusher = pusher.Client(p.httpClient)
						}
						if err := pusher.PushContext(context.Background()); err != nil {
							fmt.Printf("Error pushing to Pushgateway: %v\n", err)
						}
					} else if target.Collector != nil {
						pusher := push.New(p.PushConfig.PushGatewayURL, target.JobName).
							Grouping(pushInstanceLabel, target.JobName)
						pusher.Collector(exporterCollector.Collector)
					}
				}(target)
			}
		case <-p.sig:
			p.wg.Wait()
			p.exit <- struct{}{}
		case <-p.exit:
			return
		}
	}
}
