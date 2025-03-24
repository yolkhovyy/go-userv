package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/yolkhovyy/go-userv/internal/domain"
	"github.com/yolkhovyy/go-userv/internal/otelw"
	gqlrouter "github.com/yolkhovyy/go-userv/internal/router/graphql"
	httpserver "github.com/yolkhovyy/go-userv/internal/server/http"
	"github.com/yolkhovyy/go-utilities/buildinfo"
	"github.com/yolkhovyy/go-utilities/osx"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

const (
	domainName  = "user"
	serviceName = "user-graphql"
)

func main() {
	os.Exit(run())
}

//nolint:funlen
func run() int {
	buildInfo := buildinfo.ReadData()

	// Congig file.
	configFile := flag.String("config", "config.yml",
		"Path to the configuration file (default: config.yml)")

	flag.Parse()

	// Service configuration.
	config := NewConfig()

	err := config.Load(*configFile, domainName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "config load: %v", err)

		return osx.ExitConfigError
	}

	// The ctx.Done() channel will close when one of the signals arrives.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	// Telemetry.
	serviceAttributes := []attribute.KeyValue{
		semconv.ServiceNameKey.String(serviceName),
		semconv.ServiceVersionKey.String(buildInfo.Version),
	}

	logger, tracer, metric, err := otelw.Configure(ctx, config.Config, serviceAttributes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "otelw configure: %v", err)

		return osx.ExitFailure
	}

	defer func() {
		err := errors.Join(err,
			metric.Shutdown(ctx),
			tracer.Shutdown(ctx),
			logger.Shutdown(ctx))
		if err != nil {
			fmt.Fprintf(os.Stderr, "otelw shutdown: %v", err)
		}
	}()

	logger.InfoContext(ctx, "build info",
		slog.String("version", buildInfo.Version),
		slog.String("time", buildInfo.Time),
		slog.String("commit", buildInfo.Commit),
	)

	// Initialize user domain.
	domain, err := domain.New(ctx, config.Postgres)
	if err != nil {
		logger.ErrorContext(ctx, "domain",
			slog.String("new", err.Error()),
		)

		return osx.ExitFailure
	}

	defer func() {
		if err := domain.Close(); err != nil {
			logger.ErrorContext(ctx, "domain",
				slog.String("close", err.Error()),
			)
		}
	}()

	// Create graphql router.
	router, err := gqlrouter.New(config.Router, domain)
	if err != nil {
		logger.ErrorContext(ctx, "graphql router",
			slog.String("new", err.Error()),
		)

		return osx.ExitFailure
	}

	// Create and run HTTP server.
	server := httpserver.New(config.HTTP, router.Handler())
	if err := server.Run(ctx); err != nil {
		logger.ErrorContext(ctx, "graphql router",
			slog.String("run", err.Error()),
		)

		return osx.ExitFailure
	}

	logger.InfoContext(ctx, "exiting")

	return osx.ExitSuccess
}
