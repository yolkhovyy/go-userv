package postgres

import (
	"context"
	"fmt"
	"net"
	"strconv"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"
	"github.com/yolkhovyy/user/internal/contract/storage"
)

type Controller struct {
	connString string
	pool       *pgxpool.Pool
}

//nolint:ireturn
func New(ctx context.Context, config Config) (storage.Contract, error) {
	// TODO: fix sslmode
	connString := "postgres://" +
		config.Username + ":" +
		config.Password + "@" +
		net.JoinHostPort(config.Host, strconv.Itoa(config.Port)) + "/" +
		config.Database + "?sslmode=disable"

	pool, err := pgxpool.Connect(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("database connect: %w", err)
	}

	log.Debug().Msg("database connected")

	return &Controller{
		connString: connString,
		pool:       pool,
	}, nil
}

func (c *Controller) Close() error {
	log.Debug().Msg("database closing")

	c.pool.Close()

	log.Trace().Msg("database closed")

	return nil
}
