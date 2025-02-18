package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/yolkhovyy/user/internal/domain"
	"github.com/yolkhovyy/user/internal/logger"
	grpcrouter "github.com/yolkhovyy/user/internal/router/grpc"
	grpcserver "github.com/yolkhovyy/user/internal/server/grpc"
)

const (
	domainName  = "user"
	serviceName = "user-grpc"
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
	router := grpcrouter.New(config.Router, domain)

	// Create and run gRPC grpcServer.
	grpcServer := grpcserver.New(config.GRPC, router, grpcrouter.Interceptors()...)
	if err := grpcServer.Run(ctx); err != nil {
		log.Error().Err(err).Msg("grpc server")

		return 1
	}

	log.Info().Msg("exiting")

	return 0
}
