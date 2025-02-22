package http

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"
	cserver "github.com/yolkhovyy/user/internal/contract/server"
)

type Server struct {
	*http.Server
	config Config
}

//nolint:ireturn
func New(config Config, handler http.Handler) cserver.Contract {
	return &Server{
		Server: &http.Server{
			Addr:              net.JoinHostPort("", strconv.Itoa(config.Port)),
			Handler:           handler,
			ReadHeaderTimeout: config.ReadHeaderTimeout,
		},
		config: config,
	}
}

func (s *Server) Run(ctx context.Context) error {
	errChan := make(chan error)

	go func() {
		log.Info().Msgf("http server starting on %s", s.Server.Addr)

		if err := s.Server.ListenAndServe(); err != nil {
			errChan <- err
		}

		close(errChan)
	}()

	select {
	case <-ctx.Done():
	case err := <-errChan:
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("http server start: %w", err)
		}
	}

	log.Debug().Msg("http server shutting down")

	ctx, timeout := context.WithTimeout(ctx, s.config.ShutdownTimeout)
	defer timeout()

	if err := s.Server.Shutdown(ctx); err != nil {
		return fmt.Errorf("http server shutdown: %w", err)
	}

	log.Trace().Msg("http server shutdown complete")

	return nil
}
