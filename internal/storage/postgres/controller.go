package postgres

import (
	"context"
	"fmt"
	"net"
	"strconv"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yolkhovyy/go-otelw/otelw/slogw"
	"github.com/yolkhovyy/go-userv/internal/contract/storage"
)

type Controller struct {
	pool *pgxpool.Pool
}

//nolint:ireturn
func New(ctx context.Context, config Config) (storage.Contract, error) {
	// TODO: fix sslmode
	connString := "postgres://" +
		config.Username + ":" +
		config.Password + "@" +
		net.JoinHostPort(config.Host, strconv.Itoa(config.Port)) + "/" +
		config.Database + "?sslmode=disable"

	cfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	cfg.ConnConfig.Tracer = otelpgx.NewTracer()

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("database connect: %w", err)
	}

	slogw.DefaultLogger().DebugContext(ctx, "database connected")

	return &Controller{
		pool: pool,
	}, nil
}

func (c *Controller) Close() error {
	logger := slogw.DefaultLogger()
	logger.Debug("database closing")

	c.pool.Close()

	logger.Debug("database closed")

	return nil
}
