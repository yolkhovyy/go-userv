package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/yolkhovyy/go-userv/internal/domain"
	"github.com/yolkhovyy/go-userv/internal/logger"
	ginrouter "github.com/yolkhovyy/go-userv/internal/router/gin"
	httpserver "github.com/yolkhovyy/go-userv/internal/server/http"
	"github.com/yolkhovyy/go-utilities/osx"
)

const (
	domainName  = "user"
	serviceName = "user-rest"
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

	// The ctx.Done() channel will close when one of the signals arrives.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	// Initialize user domain.
	domain, err := domain.New(ctx, config.Postgres)
	if err != nil {
		log.Error().Err(err).Msg("domain initialization")

		return osx.ExitFailure
	}

	defer func() {
		if err := domain.Close(); err != nil {
			log.Error().Err(err).Msg("domain close")
		}
	}()

	// Create gin router.
	router := ginrouter.New(config.Router, domain)

	// Create and run HTTP server.
	server := httpserver.New(config.HTTP, router.Handler())
	if err := server.Run(ctx); err != nil {
		log.Error().Err(err).Msg("http server")

		return osx.ExitFailure
	}

	log.Info().Msg("exiting")

	return osx.ExitSuccess
}
