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

	"github.com/yolkhovyy/go-userv/cmd/user-notifier/version"
	"github.com/yolkhovyy/go-userv/internal/notifier"
	"github.com/yolkhovyy/go-userv/internal/otelw"
	"github.com/yolkhovyy/go-utilities/buildinfo"
	"github.com/yolkhovyy/go-utilities/osx"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

const (
	domainName  = "user"
	serviceName = "user-notifier"
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

	// Context, Done channel will close when one of the listed signals arrives.
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

	// Initialize notifier.
	notifier, err := notifier.New(config.Postgres, config.Kafka)
	if err != nil {
		logger.ErrorContext(ctx, "notifier",
			slog.String("new", err.Error()),
		)

		return osx.ExitFailure
	}

	// Listening for user changes and notify consumers.
	err = notifier.Run(ctx)
	if err != nil {
		logger.ErrorContext(ctx, "notifier",
			slog.String("run", err.Error()),
		)

		return osx.ExitFailure
	}

	logger.InfoContext(ctx, "exiting")

	return osx.ExitSuccess
}
