package logger

import (
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Init(name string) zerolog.Logger {
	writer := zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
		w.TimeFormat = time.RFC3339
	})

	// This initializes the global logger,
	// which can be used anywhere in the application.
	log.Logger = zerolog.New(writer).With().Timestamp().Str("service", name).Logger()

	return log.Logger
}
