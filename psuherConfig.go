package pusher

import (
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

type PushConfig struct {
	PushGatewayURL string
	InstanceLabel  string
	Interval       int // 间隔时间 单位秒
	TargetExporter []TargetExporter
	logger         *zap.Logger
}

type TargetExporter struct {
	ExporterURL string
	JobName     string
	Collector   prometheus.Collector
}

func NewPushConfig(pushGatewayURL, instanceLabel string, targetExporter []TargetExporter, logger *zap.Logger) *PushConfig {
	if logger == nil {
		logger, _ = zap.NewProductionConfig().Build()
	}
	return &PushConfig{
		PushGatewayURL: pushGatewayURL,
		InstanceLabel:  instanceLabel,
		TargetExporter: targetExporter,
		logger:         logger,
	}
}
