package logging

import (
	"log/slog"
	"os"
	"time"

	"github.com/durandj/ek/internal/configuration"
	"github.com/lmittmann/tint"
)

func NewStructuredLoggerFromConfig(
	config configuration.Configuration,
) *slog.Logger {
	var handler slog.Handler
	if config.Environment == configuration.EnvironmentDevelopment {
		handler = tint.NewHandler(os.Stdout, &tint.Options{
			AddSource:   true,
			Level:       config.Logging.Level.ToSlogLevel(),
			ReplaceAttr: nil,
			TimeFormat:  time.RFC3339,
			NoColor:     false,
		})
	} else {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource:   true,
			Level:       config.Logging.Level.ToSlogLevel(),
			ReplaceAttr: nil,
		})
	}

	return slog.New(handler)
}

func Err(err error) slog.Attr {
	return slog.Any("error", err)
}
