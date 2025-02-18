package domain

import (
	"context"
	"fmt"

	"github.com/yolkhovyy/user/internal/contract/domain"
	"github.com/yolkhovyy/user/internal/contract/storage"
	storagec "github.com/yolkhovyy/user/internal/storage"
)

type Controller struct {
	storage storage.Contract
}

//nolint:ireturn
func New(
	ctx context.Context,
	storageConfig storagec.Config,
) (domain.Contract, error) {
	var err error

	controller := Controller{}

	controller.storage, err = storagec.New(ctx, storageConfig)
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
