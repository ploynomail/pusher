package pusher

import "github.com/prometheus/client_golang/prometheus"

type PushConfig struct {
	PushGatewayURL string
	InstanceLabel  string
	Interval       int // 间隔时间 单位秒
	TargetExporter []TargetExporter
}

type TargetExporter struct {
	ExporterURL string
	JobName     string
	Collector   prometheus.Collector
}

func NewPushConfig(pushGatewayURL, instanceLabel string, targetExporter []TargetExporter) *PushConfig {
	return &PushConfig{
		PushGatewayURL: pushGatewayURL,
		InstanceLabel:  instanceLabel,
		TargetExporter: targetExporter,
	}
}
