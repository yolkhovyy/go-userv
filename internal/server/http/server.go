package http

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strconv"

	"github.com/yolkhovyy/go-otelw/otelw/slogw"
	cserver "github.com/yolkhovyy/go-userv/internal/contract/server"
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
	logger := slogw.DefaultLogger()

	errChan := make(chan error)

	go func() {
		logger.InfoContext(ctx, "http server starting",
			slog.String("addr", s.Server.Addr))

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

	logger.DebugContext(ctx, "http server shutting down")

	ctx, timeout := context.WithTimeout(ctx, s.config.ShutdownTimeout)
	defer timeout()

	if err := s.Server.Shutdown(ctx); err != nil {
		return fmt.Errorf("http server shutdown: %w", err)
	}

	logger.DebugContext(ctx, "http server shutting down complete")

	return nil
}
