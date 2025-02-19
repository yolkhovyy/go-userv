package grpc

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func Interceptors() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.UnaryInterceptor(logUnary),
		grpc.StreamInterceptor(logStream),
	}
}

func logUnary(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	start := time.Now()

	log.Info().
		Str("fullMethod", info.FullMethod).
		Msg("grpc request begin")

	resp, err := handler(ctx, req)
	if err != nil {
		log.Error().Err(err).
			Str("fullMethod", info.FullMethod).
			Dur("duration(ms)", time.Since(start).Round(time.Millisecond)).
			Msg("grpc request end")

		return resp, err
	}

	log.Info().
		Str("fullMethod", info.FullMethod).
		Dur("duration(ms)", time.Since(start).Round(time.Millisecond)).
		Msg("grpc request end")

	return resp, nil
}

func logStream(
	srv interface{},
	stream grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	start := time.Now()

	log.Info().
		Str("fullMethod", info.FullMethod).
		Msg("grpc request begin")

	err := handler(srv, stream)
	if err != nil {
		log.Error().Err(err).
			Str("fullMethod", info.FullMethod).
			Dur("duration(ms)", time.Since(start).Round(time.Millisecond)).
			Msg("grpc request")

		return err
	}

	log.Info().
		Str("fullMethod", info.FullMethod).
		Dur("duration(ms)", time.Since(start).Round(time.Millisecond)).
		Msg("grpc request end")

	return nil
}
