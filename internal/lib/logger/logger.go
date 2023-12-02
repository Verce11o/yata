package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	log *zap.SugaredLogger
}

func NewLogger(mode string) *zap.SugaredLogger {

	if mode == "prod" {
		conf := zap.NewProductionConfig()
		conf.DisableStacktrace = true // FIXME
		logger, _ := conf.Build()
		return logger.Sugar()
	}
	logger, err := zap.Config{
		Level:             zap.NewAtomicLevelAt(zapcore.DebugLevel),
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "console",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalColorLevelEncoder,

			TimeKey:    "time",
			EncodeTime: zapcore.ISO8601TimeEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},

		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}.Build()

	if err != nil {
		panic(err)
	}

	return logger.Sugar()
}
