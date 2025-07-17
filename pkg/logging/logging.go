package logging

import (
	"log/slog"
	"os"
)

func init() {
	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))

	slog.SetDefault(logger)
}
