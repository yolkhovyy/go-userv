package grpc

import (
	"context"
	"fmt"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func Interceptors() []grpc.ServerOption {
	logOpts := []logging.Option{
		logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
	}

	return []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(recovery.UnaryServerInterceptor(panicRecovery)),
		grpc.ChainStreamInterceptor(recovery.StreamServerInterceptor(panicRecovery)),
		grpc.ChainUnaryInterceptor(logging.UnaryServerInterceptor(zerologInterceptor(log.Logger), logOpts...)),
		grpc.ChainStreamInterceptor(logging.StreamServerInterceptor(zerologInterceptor(log.Logger), logOpts...)),
	}
}

//nolint:gochecknoglobals
var panicRecovery = recovery.WithRecoveryHandler(func(p any) error {
	return fmt.Errorf("%v recovered from panic: %w", p, ErrPanic)
})

//nolint:ireturn
func zerologInterceptor(logger zerolog.Logger) logging.Logger {
	return logging.LoggerFunc(func(_ context.Context, level logging.Level, msg string, fields ...any) {
		switch level {
		case logging.LevelDebug:
			logger.Debug().Fields(fields).Msg(msg)
		case logging.LevelInfo:
			logger.Info().Fields(fields).Msg(msg)
		case logging.LevelWarn:
			logger.Warn().Fields(fields).Msg(msg)
		case logging.LevelError:
			logger.Error().Fields(fields).Msg(msg)
		default:
			logger.Trace().Fields(fields).Msg(msg)
		}
	})
}
