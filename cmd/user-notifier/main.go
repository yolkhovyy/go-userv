package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/yolkhovyy/user/internal/notifier"
)

const (
	domainName  = "user"
	serviceName = "user-notifier"
)

func main() {
	os.Exit(run())
}

func run() int {
	// Logger initialization.
	writer := zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
		w.TimeFormat = time.RFC3339
	})
	log.Logger = zerolog.New(writer).With().Timestamp().Str("service", serviceName).Logger()
	log.Info().Msg("starting")

	// Service configuration.
	configFile := flag.String("config", "config.yml", "Path to the configuration file (default: config.yml)")
	flag.Parse()

	config := Config{}

	err := config.Load(*configFile, domainName)
	if err != nil {
		log.Error().Err(err).Msg("config load")

		return 1
	}

	// Context, Done channel will close when one of the listed signals arrives.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	// Initialize notifier.
	notifier, err := notifier.New(config.Postgres, config.Kafka)
	if err != nil {
		log.Error().Err(err).Msg("notifiern initialization")

		return 1
	}

	// Listening for user changes and notify consumers.
	err = notifier.Run(ctx)
	if err != nil {
		log.Error().Err(err).Msg("listen and notify")

		return 1
	}

	log.Info().Msg("exiting")

	return 0
}
