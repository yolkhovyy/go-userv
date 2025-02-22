package domain

import (
	"context"
	"fmt"
	"io"

	"github.com/google/uuid"
	"github.com/yolkhovyy/user/internal/contract/storage"
	"golang.org/x/crypto/bcrypt"
)

type Contract interface {
	CRUD
	io.Closer
}

type CRUD interface {
	Creator
	Reader
	Updater
	Deleter
}

type Creator interface {
	Create(ctx context.Context, user UserInput) (*User, error)
}

type Reader interface {
	Get(ctx context.Context, userID uuid.UUID) (*User, error)
	List(ctx context.Context, page, limit int, countryCode string) (*UserList, error)
}

type Updater interface {
	Update(ctx context.Context, user UserInput) (*User, error)
}

type Deleter interface {
	Delete(ctx context.Context, userID uuid.UUID) error
}

// TODO: Replace these by dedicated domain types if needed.

type User storage.User

type UserInput storage.UserInput

type UserList storage.UserList

func UserFromStorage(user storage.User) User {
	return User(user)
}

func UserInputToStorage(userInput UserInput) storage.UserInput {
	return storage.UserInput(userInput)
}

func UsersFromStorage(users storage.UserList) UserList {
	return UserList(users)
}

func (u *UserInput) HashPassword() error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("password hash: %w", err)
	}

	u.Password = string(passwordHash)

	return nil
}
