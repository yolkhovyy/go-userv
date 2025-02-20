package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/yolkhovyy/user/internal/domain"
	"github.com/yolkhovyy/user/internal/logger"
	graphqlrouter "github.com/yolkhovyy/user/internal/router/graphql"
	httpserver "github.com/yolkhovyy/user/internal/server/http"
)

const (
	domainName  = "user"
	serviceName = "user-graphql"
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

	// The ctx.Done() channel will close when one of the signals arrives.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	// Initialize user domain.
	domain, err := domain.New(ctx, config.Postgres)
	if err != nil {
		log.Error().Err(err).Msg("domain initialization")

		return 1
	}

	defer func() {
		if err := domain.Close(); err != nil {
			log.Error().Err(err).Msg("domain close")
		}
	}()

	// Create router.
	router, err := graphqlrouter.New(config.Router, domain)
	if err != nil {
		log.Error().Err(err).Msg("router creation")

		return 1
	}

	// Create and run HTTP server.
	server := httpserver.New(config.HTTP, router.Handler())
	if err := server.Run(ctx); err != nil {
		log.Error().Err(err).Msg("http server")

		return 1
	}

	log.Info().Msg("exiting")

	return 0
}
