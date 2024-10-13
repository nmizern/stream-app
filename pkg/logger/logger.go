package logger

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	SugaredLogger *zap.SugaredLogger
}

func NewLogger(level string) (*Logger, error) {
	var zapLevel zapcore.Level

	if err := zapLevel.UnmarshalText([]byte(strings.ToUpper(level))); err != nil {
		return nil, err
	}

	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(zapLevel),
		Development: false,
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:       "time",
			LevelKey:      "level",
			NameKey:       "logger",
			CallerKey:     "caller",
			MessageKey:    "msg",
			StacktraceKey: "stacktrace",
			LineEnding:    zapcore.DefaultLineEnding,
			EncodeLevel:   zapcore.LowercaseLevelEncoder,
			EncodeTime:    zapcore.ISO8601TimeEncoder,
			EncodeCaller:  zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	zapLogger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return &Logger{SugaredLogger: zapLogger.Sugar()}, nil
}
