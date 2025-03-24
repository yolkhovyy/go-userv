package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"strconv"

	"github.com/yolkhovyy/go-otelw/pkg/slogw"
	"github.com/yolkhovyy/go-userv/contract/proto"
	"github.com/yolkhovyy/go-userv/internal/contract/server"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	server *grpc.Server
	config Config
}

//nolint:ireturn
func New(config Config, serviceServer proto.UserServiceServer, opts ...grpc.ServerOption) server.Contract {
	serverHandler := otelgrpc.NewServerHandler(
		otelgrpc.WithTracerProvider(otel.GetTracerProvider()),
	)

	opts = append(opts, grpc.StatsHandler(serverHandler))

	grpcServer := grpc.NewServer(opts...)
	proto.RegisterUserServiceServer(grpcServer, serviceServer)

	if config.Reflection {
		reflection.Register(grpcServer)
	}

	return &Server{
		server: grpcServer,
		config: config,
	}
}

func (s *Server) Run(ctx context.Context) error {
	logger := slogw.DefaultLogger()

	listener, err := net.Listen("tcp", net.JoinHostPort("", strconv.Itoa(s.config.Port)))
	if err != nil {
		return fmt.Errorf("net listen: %w", err)
	}

	errChan := make(chan error)

	go func() {
		logger.InfoContext(ctx, "grpc server starting",
			slog.String("address", listener.Addr().String()),
		)

		if err := s.server.Serve(listener); err != nil {
			errChan <- err
		}

		close(errChan)
	}()

	select {
	case <-ctx.Done():
	case err := <-errChan:
		if err != nil {
			return fmt.Errorf("grpc server start: %w", err)
		}
	}

	logger.DebugContext(ctx, "grpc server shutting down")

	s.server.GracefulStop()

	logger.DebugContext(ctx, "grpc server shutdown complete")

	return nil
}
