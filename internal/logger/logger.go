package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"github.com/squzy/squzy/internal/logger/config"
	"time"
)

var (
	l *zap.Logger
)

func init() {
	var encoderConfig = zapcore.EncoderConfig{
		TimeKey:       "ts",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		EncodeTime: zapcore.TimeEncoder(func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.UTC().Format("2006-01-02T15:04:05Z0700"))
			// 2019-08-13T04:39:11Z
		}),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= getLevel() && lvl < zapcore.ErrorLevel
	})

	var core zapcore.Core
	if getLevel() >= zap.ErrorLevel {
		core = zapcore.NewTee(
			zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), zapcore.Lock(os.Stderr), getLevel()),
		)
	} else {
		core = zapcore.NewTee(
			zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), zapcore.Lock(os.Stdout), lowPriority),
			zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), zapcore.Lock(os.Stderr), highPriority),
		)
	}

	l = zap.New(core, zap.AddCaller())
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
