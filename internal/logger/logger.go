package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
)

var (
	Logger ZapLogger
	cfg zap.Config
)

type ZapLogger struct {
	*zap.Logger
}

func init() {
	var err error
	cfg = zap.Config{
		Level:             zap.NewAtomicLevelAt(zapcore.DebugLevel),
		Encoding:          "json",
		EncoderConfig:     zapcore.EncoderConfig{
			MessageKey:     "message",
			LevelKey:       "level",
			LineEnding:     "\n",
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
		},
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
	}
	Logger.Logger, err = cfg.Build()
	if err != nil {
		log.Fatal("Could not build logger")
	}
}

func SetLevel(l zapcore.Level){
	cfg.Level.SetLevel(l)
}

func (l *ZapLogger) Errorf(format string, args ...interface{}) {
	l.Logger.Error(fmt.Sprintf(format, args...))
}

func (l *ZapLogger) Fatalf(format string, args ...interface{}) {
	l.Logger.Fatal(fmt.Sprintf(format, args...))
}

