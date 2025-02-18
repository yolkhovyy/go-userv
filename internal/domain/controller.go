package domain

import (
	"context"
	"fmt"

	"github.com/yolkhovyy/go-userv/internal/contract/domain"
	"github.com/yolkhovyy/go-userv/internal/contract/storage"
	"github.com/yolkhovyy/go-userv/internal/storage/postgres"
)

type Controller struct {
	storage storage.Contract
}

//nolint:ireturn
func New(
	ctx context.Context,
	config postgres.Config,
) (domain.Contract, error) {
	var err error

	controller := Controller{}

	controller.storage, err = postgres.New(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("storage: %w", err)
	}

	return controller, nil
}

func (u Controller) Close() error {
	var errAll error

	if u.storage != nil {
		if err := u.storage.Close(); err != nil {
			errAll = fmt.Errorf("storage: %w", err)
		}
	}

	return errAll
}
