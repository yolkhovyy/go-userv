package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/yolkhovyy/user/internal/logger"
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
	log := logger.Init(serviceName)
	log.Info().Msg("starting")

	// Service configuration.
	config := Config{}

	err := config.Load(domainName)
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
