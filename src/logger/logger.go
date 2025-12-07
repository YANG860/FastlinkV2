package logger

import (
	"log/slog"
	"os"
)

var l *slog.Logger

func init() {
	l = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	}))

	slog.SetDefault(l)
}
