package configuration

import (
	"fmt"
	"log/slog"

	"github.com/kelseyhightower/envconfig"
)

type LoggingConfiguration struct {
	Level LoggingLevel `default:"warn"`
}

type LoggingLevel slog.Level

var _ envconfig.Decoder = (*LoggingLevel)(nil)

func (l *LoggingLevel) Decode(value string) error {
	var level slog.Level
	if err := level.UnmarshalText([]byte(value)); err != nil {
		return fmt.Errorf("unable to parse logging level: %w", err)
	}

	*l = LoggingLevel(level)

	return nil
}

func (l *LoggingLevel) ToSlogLevel() slog.Level {
	return slog.Level(*l)
}
