package domain

import (
	"context"
	"fmt"

	"github.com/yolkhovyy/user/internal/contract/domain"
	storagec "github.com/yolkhovyy/user/internal/contract/storage"
	storage "github.com/yolkhovyy/user/internal/storage/postgres"
)

type Controller struct {
	storage storagec.Contract
}

//nolint:ireturn
func New(
	ctx context.Context,
	config storage.Config,
) (domain.Contract, error) {
	var err error

	controller := Controller{}

	controller.storage, err = storage.New(ctx, config)
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
