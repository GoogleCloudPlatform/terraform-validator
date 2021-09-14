package cmd

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func newTestErrorLogger(verbose, useStructuredLogging bool) (*zap.Logger, *observer.ObservedLogs) {
	logger := newErrorLogger(verbose, useStructuredLogging)
	observerCore, logs := observer.New(zap.DebugLevel)
	wrapper := zap.WrapCore(func(zapcore.Core) zapcore.Core {
		return observerCore
	})

	return logger.WithOptions(wrapper), logs
}

func newTestOutputLogger() (*zap.Logger, *observer.ObservedLogs) {
	logger := newOutputLogger()
	observerCore, logs := observer.New(zap.DebugLevel)
	wrapper := zap.WrapCore(func(zapcore.Core) zapcore.Core {
		return observerCore
	})

	return logger.WithOptions(wrapper), logs
}
