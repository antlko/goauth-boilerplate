package logger

import (
	"io"
	"log/slog"
	"os"
)

type Config struct {
	AppName  string
	Hostname string
	Writer   io.Writer
}

func InitLogger(config Config) {
	writer := config.Writer
	if writer == nil {
		writer = os.Stdout
	}
	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}

	jsonHandler := slog.
		NewJSONHandler(writer, opts).
		WithAttrs([]slog.Attr{
			slog.String("application", config.AppName),
		})

	logger := slog.New(jsonHandler)
	logger = logger.With("hostname", config.Hostname)

	slog.SetDefault(logger)
}
