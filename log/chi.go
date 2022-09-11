package log

import (
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

type ChiLogger struct {
	middleware.LoggerInterface

	Logger *zap.Logger
}

func (c *ChiLogger) Print(v ...interface{}) {
	c.Logger.Info("http request", zap.Any("args", v))
}
