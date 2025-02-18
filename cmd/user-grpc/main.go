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
	"github.com/yolkhovyy/user/internal/domain"
	router "github.com/yolkhovyy/user/internal/router/grpc"
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
	router := router.New(config.Router, domain)

	// Create and run gRPC grpcServer.
	grpcServer := grpcserver.New(config.GRPC, router)
	if err := grpcServer.Run(ctx); err != nil {
		log.Error().Err(err).Msg("grpc server")

		return 1
	}

	log.Info().Msg("exiting")

	return 0
}
