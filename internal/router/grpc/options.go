package grpc

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/yolkhovyy/go-otelw/pkg/slogw"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
)

func Options() []grpc.ServerOption {
	logOpts := []logging.Option{
		logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
	}

	traceHandler := otelgrpc.NewServerHandler(
		otelgrpc.WithTracerProvider(otel.GetTracerProvider()),
	)

	return []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			recovery.UnaryServerInterceptor(panicRecovery),
			logging.UnaryServerInterceptor(logging.LoggerFunc(slogWrapper), logOpts...),
		),
		grpc.ChainStreamInterceptor(
			recovery.StreamServerInterceptor(panicRecovery),
			logging.StreamServerInterceptor(logging.LoggerFunc(slogWrapper), logOpts...),
		),
		grpc.StatsHandler(traceHandler),
	}
}

//nolint:gochecknoglobals
var panicRecovery = recovery.WithRecoveryHandler(func(p any) error {
	return fmt.Errorf("%v recovered from panic: %w", p, ErrPanic)
})

func slogWrapper(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
	attrs := make([]slog.Attr, 0, len(fields))

	for i := 0; i < len(fields)-1; i += 2 {
		key, ok := fields[i].(string)
		if !ok {
			continue
		}

		attrs = append(attrs, slog.Any(key, fields[i+1]))
	}

	// Convert grpc-middleware log level to slog log level
	var slogLevel slog.Level

	switch lvl {
	case logging.LevelDebug:
		slogLevel = slog.LevelDebug
	case logging.LevelInfo:
		slogLevel = slog.LevelInfo
	case logging.LevelWarn:
		slogLevel = slog.LevelWarn
	case logging.LevelError:
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelInfo
	}

	logger := slogw.DefaultLogger()
	logger.LogAttrs(ctx, slogLevel, msg, attrs...)
}
