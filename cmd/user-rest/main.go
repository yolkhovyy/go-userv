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

	"github.com/yolkhovyy/go-userv/cmd/user-rest/version"
	"github.com/yolkhovyy/go-userv/internal/domain"
	"github.com/yolkhovyy/go-userv/internal/otelw"
	ginrouter "github.com/yolkhovyy/go-userv/internal/router/gin"
	httpserver "github.com/yolkhovyy/go-userv/internal/server/http"
	"github.com/yolkhovyy/go-utilities/buildinfo"
	"github.com/yolkhovyy/go-utilities/osx"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

const (
	domainName  = "user"
	serviceName = "user-rest"
)

func main() {
	os.Exit(run())
}

//nolint:funlen
func run() int {
	buildInfo := buildinfo.ReadData()

	// Config file.
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
		semconv.ServiceVersionKey.String(version.Tag),
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
		slog.String("version", version.Tag),
		slog.String("time", buildInfo.Time),
		slog.String("commit", buildInfo.Revision),
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

	// Create gin router.
	router := ginrouter.New(config.Router, domain, otelgin.Middleware(serviceName))

	// Create and run HTTP server.
	server := httpserver.New(config.HTTP, router.Handler())
	if err := server.Run(ctx); err != nil {
		logger.ErrorContext(ctx, "http server",
			slog.String("run", err.Error()),
		)

		return osx.ExitFailure
	}

	logger.InfoContext(ctx, "exiting")

	return osx.ExitSuccess
}
