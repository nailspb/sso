package slogHelper

import (
	"context"
	"log/slog"
	"os"
)

func GetErrAttr(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

func ConfigureForMiddleware(log *slog.Logger, name string) *slog.Logger {
	log.Info("add middleware " + name)
	return log.With(slog.String("middleware", name))
}

func AddOperation(log *slog.Logger, op string) *slog.Logger {
	return log.With(slog.String("op", op))
}

func AddRequestId(log *slog.Logger, context context.Context) *slog.Logger {
	if id, ok := context.Value("X-Request-Id").(string); ok {
		return log.With("requestId", id)
	}
	log.Warn("Not set X-Request-Id")
	return log
}

func NewLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case "local":
		{
			log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		}
	case "dev":
		{
			log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		}
	case "prod":
		{
			log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
		}

	}
	return log
}
