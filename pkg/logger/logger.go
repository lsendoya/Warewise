package logger

import (
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.SugaredLogger

func init() {
	cfg := zap.NewProductionConfig()

	encoderCFG := zap.NewProductionEncoderConfig()
	encoderCFG.TimeKey = "timestamp"
	encoderCFG.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCFG.StacktraceKey = ""

	cfg.EncoderConfig = encoderCFG

	zapLogger, err := cfg.Build(zap.AddCallerSkip(1))
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}

	logger = zapLogger.Sugar()
}

func Info(msg string, args ...interface{}) {
	if logger != nil {
		logger.Infof(msg, args...)
	}
}

func Errorf(msg string, args ...interface{}) {
	if logger != nil {
		logger.Errorf(msg, args...)

	}
}

func Fatalf(msg string, args ...interface{}) {
	if logger != nil {
		logger.Fatalf(msg, args...)
	}
}

func SyncLogger() {
	if err := logger.Sync(); err != nil {
		log.Printf("can't sync zap logger: %v", err)
	}
}
