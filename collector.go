package pusher

import (
	"net/http"

	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
)

type exporterCollector struct {
	client *http.Client
	url    string
}

// fetchAndDecodeMetrics 从指定的URL获取指标并将其解码为MetricFamily切片。
func (cc *exporterCollector) Gather() ([]*dto.MetricFamily, error) {
	resp, err := cc.client.Get(cc.url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, err
	}
	defer resp.Body.Close()

	p := expfmt.NewDecoder(resp.Body, expfmt.Format("text/plain"))

	var metricFamilies []*dto.MetricFamily
	for {
		mf := &dto.MetricFamily{}
		if err := p.Decode(mf); err != nil {
			break
		}
		metricFamilies = append(metricFamilies, mf)
	}

	return metricFamilies, nil
}
