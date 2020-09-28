package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"squzy/config"
)

var (
	l *zap.Logger
)

func init() {
	var err error
	cfg := zap.Config{
		Level:             zap.NewAtomicLevelAt(getLevel()),
		Development: false,
		Encoding:          "json",
		EncoderConfig:     zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.EpochTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:       []string{"stderr"},
		ErrorOutputPaths:  []string{"stderr"},
	}
	l, err = cfg.Build()
	if err != nil {
		log.Fatal("Could not build logger")
	}
}

func Info(msg string) {
	l.Info(msg)
}

func Error(msg string) {
	l.Error(msg)
}

func Fatal(msg string) {
	l.Fatal(msg)
}

func Panic(msg string) {
	l.Panic(msg)
}

func Infof(format string, args ...interface{}) {
	l.Info(fmt.Sprintf(format, args...))
}

func Errorf(format string, args ...interface{}) {
	l.Error(fmt.Sprintf(format, args...))
}

func Fatalf(format string, args ...interface{}) {
	l.Fatal(fmt.Sprintf(format, args...))
}

func Panicf(format string, args ...interface{}) {
	l.Panic(fmt.Sprintf(format, args...))
}

func getLevel() zapcore.Level {
	c := config.New()
	ll := c.GetLogLevel()
	var level zapcore.Level

	if ll == "" {
		level = zapcore.FatalLevel
	} else {
		err := level.UnmarshalText([]byte(ll))
		if err != nil {
			log.Fatal("Could not get logger level")
		}
	}
	return level
}

