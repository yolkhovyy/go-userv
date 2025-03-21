package domain

import (
	"context"
	"fmt"
	"io"

	"github.com/google/uuid"
	"github.com/yolkhovyy/go-userv/internal/contract/storage"
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
	Update(ctx context.Context, user UserUpdate) (*User, error)
}

type Deleter interface {
	Delete(ctx context.Context, userID uuid.UUID) error
}

// TODO: Replace these by dedicated domain types if needed.

type User storage.User

type UserInput storage.UserInput

type UserUpdate storage.UserUpdate

type UserList struct {
	Users      []User `json:"users,omitempty"`
	TotalCount int    `json:"totalCount"`
	NextPage   int    `json:"nextPage"`
}

func UserFromStorage(user storage.User) User {
	return User(user)
}

func UserInputToStorage(userInput UserInput) storage.UserInput {
	return storage.UserInput(userInput)
}

func UserUpdateToStorage(userUpdate UserUpdate) storage.UserUpdate {
	return storage.UserUpdate(userUpdate)
}

func UsersFromStorage(storageUsers []storage.User) []User {
	users := make([]User, len(storageUsers))

	for i, u := range storageUsers {
		users[i] = UserFromStorage(u)
	}

	return users
}

func HashPassword(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("password hash: %w", err)
	}

	return string(passwordHash), nil
}
