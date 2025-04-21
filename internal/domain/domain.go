package domain

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/yolkhovyy/go-otelw/otelw/slogw"
	"github.com/yolkhovyy/go-userv/internal/contract/domain"
)

func (u Controller) Create(ctx context.Context, userInput domain.UserInput) (*domain.User, error) {
	hashedPassword, err := domain.HashPassword(userInput.Password)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	userInput.Password = hashedPassword

	user, err := u.storage.Create(ctx, domain.UserInputToStorage(userInput))
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	createdUser := domain.User(*user)

	slogw.DefaultLogger().InfoContext(ctx, "created ",
		slog.Any("user", createdUser),
	)

	return &createdUser, nil
}

func (u Controller) Update(ctx context.Context, userUpdate domain.UserUpdate) (*domain.User, error) {
	var hashedPassword string

	if userUpdate.Password != "" {
		var err error

		hashedPassword, err = domain.HashPassword(userUpdate.Password)
		if err != nil {
			return nil, fmt.Errorf("create user: %w", err)
		}
	}

	userUpdate.Password = hashedPassword

	user, err := u.storage.Update(ctx, domain.UserUpdateToStorage(userUpdate))
	if err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	updatedUser := domain.User(*user)

	slogw.DefaultLogger().InfoContext(ctx, "updated",
		slog.Any("user", updatedUser),
	)

	return &updatedUser, nil
}

func (u Controller) Get(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	user, err := u.storage.Get(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	domainUser := domain.User(*user)

	slogw.DefaultLogger().InfoContext(ctx, "retrieved",
		slog.Any("user", domainUser),
	)

	return &domainUser, nil
}

func (u Controller) List(ctx context.Context, page, limit int, countryCode string) (*domain.UserList, error) {
	users, count, err := u.storage.List(ctx, page, limit, countryCode)
	if err != nil {
		return nil, fmt.Errorf("get users: %w", err)
	}

	if len(users) > 0 {
		slogw.DefaultLogger().InfoContext(ctx, "retrieved",
			slog.Int("number of users", len(users)),
			slog.Any("first user", users[0]),
		)
	} else {
		slogw.DefaultLogger().InfoContext(ctx, "retrieved no users")
	}

	nextPage := page + 1
	if len(users) < limit {
		nextPage = -1
	}

	listUsers := domain.UserList{
		Users:      domain.UsersFromStorage(users),
		TotalCount: count,
		NextPage:   nextPage,
	}

	return &listUsers, nil
}

func (u Controller) Delete(ctx context.Context, userID uuid.UUID) error {
	if err := u.storage.Delete(ctx, userID); err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	slogw.DefaultLogger().InfoContext(ctx, "deleted",
		slog.Any("user", userID.String()),
	)

	return nil
}
