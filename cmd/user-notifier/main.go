package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/yolkhovyy/go-userv/internal/logger"
	"github.com/yolkhovyy/go-userv/internal/notifier"
	"github.com/yolkhovyy/go-utilities/osx"
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

	// Congig file.
	configFile := flag.String("config", "config.yml",
		"Path to the configuration file (default: config.yml)")

	flag.Parse()

	// Service configuration.
	config := NewConfig()

	err := config.Load(*configFile, domainName)
	if err != nil {
		log.Error().Err(err).Msg("config load")

		return osx.ExitConfigError
	}

	// Context, Done channel will close when one of the listed signals arrives.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	// Initialize notifier.
	notifier, err := notifier.New(config.Postgres, config.Kafka)
	if err != nil {
		log.Error().Err(err).Msg("notifiern initialization")

		return osx.ExitFailure
	}

	// Listening for user changes and notify consumers.
	err = notifier.Run(ctx)
	if err != nil {
		log.Error().Err(err).Msg("listen and notify")

		return osx.ExitFailure
	}

	log.Info().Msg("exiting")

	return osx.ExitSuccess
}
