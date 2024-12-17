package pusher

import (
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

type Logger interface {
	Log(keyvals ...interface{}) error
}

type noopLogger struct {
	log *zap.Logger
}

func (l *noopLogger) Log(keyvals ...interface{}) error {
	l.log.Debug("message", zap.Any("keyvals", keyvals))
	return nil
}

func NewLogger() Logger {
	return &noopLogger{
		log: zap.NewExample(),
	}
}

type PushConfig struct {
	PushGatewayURL string
	InstanceLabel  string
	Interval       int // 间隔时间 单位秒
	TargetExporter []TargetExporter
	logger         Logger
}

type TargetExporter struct {
	ExporterURL string
	JobName     string
	Collector   prometheus.Collector
}

func NewPushConfig(pushGatewayURL, instanceLabel string, targetExporter []TargetExporter, logger Logger) *PushConfig {
	if logger == nil {
		logger = NewLogger()
	}
	return &PushConfig{
		PushGatewayURL: pushGatewayURL,
		InstanceLabel:  instanceLabel,
		TargetExporter: targetExporter,
		logger:         logger,
	}
}
