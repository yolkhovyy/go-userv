package domain

import (
	"context"
	"fmt"
	"io"

	"github.com/google/uuid"
	usergrpc "github.com/yolkhovyy/user/contract/proto"
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
	Update(ctx context.Context, user UserUpdate) (*User, error)
}

type Deleter interface {
	Delete(ctx context.Context, userID uuid.UUID) error
}

// TODO: Replace these by dedicated domain types if needed.

type User storage.User

type UserInput storage.UserInput

type UserUpdate storage.UserUpdate

type UserList storage.UserList

func UserFromStorage(user storage.User) User {
	return User(user)
}

func UserToStorage(user User) storage.User {
	return storage.User(user)
}

func UserInputToStorage(userInput UserInput) storage.UserInput {
	return storage.UserInput(userInput)
}

func UserUpdateToStorage(userUpdate UserUpdate) storage.UserUpdate {
	return storage.UserUpdate(userUpdate)
}

func UserListFromStorage(userList storage.UserList) UserList {
	return UserList(userList)
}

func UserFromGrpc(user *usergrpc.User) User {
	return User{
		ID:        uuid.MustParse(user.GetId()),
		FirstName: user.GetFirstName(),
		LastName:  user.GetLastName(),
		Nickname:  user.GetNickname(),
		Email:     user.GetEmail(),
		Country:   user.GetCountry(),
	}
}

func UserListFromGrpc(userList *usergrpc.Users) UserList {
	users := make([]storage.User, len(userList.GetUsers()))

	for i, u := range userList.GetUsers() {
		du := UserToStorage(UserFromGrpc(u))
		users[i] = du
	}

	return UserList{
		Users:      users,
		TotalCount: int(userList.GetTotalCount()),
		NextPage:   int(userList.GetNextPage()),
	}
}

func HashPassword(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("password hash: %w", err)
	}

	return string(passwordHash), nil
}
