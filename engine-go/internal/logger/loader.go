package logger

import (
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func LoadLogger(config *config.Configuration) (*zap.Logger, error) {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logFormat := "console"

	if config.Log.Format == "json" || config.Log.Format == "console" {
		logFormat = config.Log.Format
	}

	var logLevel zap.AtomicLevel
	switch config.Log.Level {
	case "DEBUG":
		logLevel = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "INFO":
		logLevel = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "WARN":
		logLevel = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "ERROR":
		logLevel = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	default:
		logLevel = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}

	zapConfig := zap.Config{
		Level:            logLevel,
		Development:      false,
		Encoding:         logFormat,
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	logger, err := zapConfig.Build()
	if err != nil {
		return nil, err
	}
	return logger, nil
}
