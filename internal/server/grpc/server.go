package grpc

import (
	"context"
	"fmt"
	"net"
	"strconv"

	"github.com/rs/zerolog/log"
	userpb "github.com/yolkhovyy/user/contract/proto"
	"github.com/yolkhovyy/user/internal/contract/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	server *grpc.Server
	config Config
}

//nolint:ireturn
func New(config Config, serviceServer userpb.UserServiceServer, opts ...grpc.ServerOption) server.Contract {
	grpcServer := grpc.NewServer(opts...)
	userpb.RegisterUserServiceServer(grpcServer, serviceServer)

	if config.Reflection {
		reflection.Register(grpcServer)
	}

	return &Server{
		server: grpcServer,
		config: config,
	}
}

func (s *Server) Run(ctx context.Context) error {
	listener, err := net.Listen("tcp", net.JoinHostPort("", strconv.Itoa(s.config.Port)))
	if err != nil {
		return fmt.Errorf("net listen: %w", err)
	}

	errChan := make(chan error)

	go func() {
		log.Info().Msgf("grpc server starting on %s", listener.Addr())

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

	log.Debug().Msg("grpc server shutting down")

	s.server.GracefulStop()

	log.Trace().Msg("grpc server shutdown complete")

	return nil
}
