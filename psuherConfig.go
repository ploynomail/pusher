package pusher

import (
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

type Level int8

type Logger interface {
	Log(level Level, keyvals ...interface{}) error
}

type noopLogger struct {
	log *zap.Logger
}

func (l *noopLogger) Log(level Level, keyvals ...interface{}) error {
	switch level {
	case 0:
		l.log.Debug("", zap.Any("keyvals", keyvals))
	case 1:
		l.log.Info("", zap.Any("keyvals", keyvals))
	case 2:
		l.log.Warn("", zap.Any("keyvals", keyvals))
	case 3:
		l.log.Error("", zap.Any("keyvals", keyvals))
	}
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
