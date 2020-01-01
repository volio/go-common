package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func L() *zap.Logger {
	if logger == nil {
		return zap.L()
	}
	return logger
}

func InitLog(c Config) {
	var level zapcore.Level
	switch c.Level {
	case "panic":
		level = zapcore.PanicLevel
	case "fatal":
		level = zapcore.FatalLevel
	case "error":
		level = zapcore.ErrorLevel
	case "warn":
		level = zapcore.WarnLevel
	case "info":
		level = zapcore.InfoLevel
	case "debug":
		level = zapcore.DebugLevel
	default:
		level = zapcore.DebugLevel
	}

	var zc zap.Config
	if c.Development {
		zc = zap.NewDevelopmentConfig()
		zc.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		zc = zap.NewProductionConfig()
		zc.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}
	zc.Level = zap.NewAtomicLevelAt(level)
	zc.DisableStacktrace = c.DisableStackTrace

	var err error
	logger, err = zc.Build()
	if err != nil {
		panic("build log failed: " + err.Error())
	}
}
