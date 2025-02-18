package server

import "context"

// Contract defines the interface for a server that can be run and shut down gracefully.
type Contract interface {
	// Run starts the server and blocks until it receives a shutdown signal
	// (e.g. syscall.SIGTERM or syscall.SIGINT) or the provided context is canceled.
	// The given context `ctx` should be configured to close its Done channel
	// when a shutdown signal is received.  Run will then perform a graceful
	// shutdown of the server before returning.  If the context is canceled
	// before a shutdown signal, Run will return ctx.Err() after attempting
	// a graceful shutdown.  Any errors encountered during startup, running,
	// or shutdown should be returned.
	Run(ctx context.Context) error
}
