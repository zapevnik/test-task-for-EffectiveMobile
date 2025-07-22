package logger

import (
	"context"
	"log/slog"
	"os"
)

type prefixHandler struct {
    slog.Handler
    prefix string
}

func (h *prefixHandler) Handle(ctx context.Context, record slog.Record) error {
    record.Message = h.prefix + " " + record.Message
    return h.Handler.Handle(ctx, record)
}

var Log *slog.Logger

func Init(levelCfg string) {
		var level slog.Level

		switch levelCfg {
		case "debug":
				level = slog.LevelDebug
		case "info":
				level = slog.LevelInfo
		case "warn":
				level = slog.LevelWarn
		case "error":
				level = slog.LevelError
		default:
				level = slog.LevelInfo
		}
		
    baseHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
        Level:     level,
        AddSource: true,
    })

    prefixedHandler := &prefixHandler{
        Handler: baseHandler,
        prefix:  "[Effective Mobile test-project]",
    }

    Log = slog.New(prefixedHandler)
}
