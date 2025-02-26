package dto

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/yolkhovyy/user/contract/proto"
	"github.com/yolkhovyy/user/internal/contract/domain"
)

// TODO: Replace these by dedicated DTO types if needed.

type User domain.User

type UserInput domain.UserInput

type UserUpdate domain.UserUpdate

type UserList struct {
	Users      []User `json:"users,omitempty"`
	TotalCount int    `json:"totalCount"`
	NextPage   int    `json:"nextPage"`
}

func UserFromDomain(user domain.User) User {
	return User(user)
}

func UserToDomain(user User) domain.User {
	return domain.User(user)
}

// func UserToStorage(user User) storage.User {
// 	return storage.User(user)
// }

func UserInputToDomain(userInput UserInput) domain.UserInput {
	return domain.UserInput(userInput)
}

func UserUpdateToDomain(userUpdate UserUpdate) domain.UserUpdate {
	return domain.UserUpdate(userUpdate)
}

func UserListFromDomain(userList domain.UserList) UserList {
	return UserList{
		Users:      UsersFromDomain(userList.Users),
		TotalCount: userList.TotalCount,
		NextPage:   userList.NextPage,
	}
}

// func UserListToDomain(userList UserList) domain.UserList {
// 	return domain.UserList(userList)
// }

func UsersFromDomain(domainUsers []domain.User) []User {
	users := make([]User, len(domainUsers))

	for i, u := range domainUsers {
		du := UserFromDomain(u)
		users[i] = du
	}

	return users
}

func UserFromProto(user *proto.User) (*User, error) {
	userID, err := uuid.Parse(user.GetId())
	if err != nil {
		return nil, fmt.Errorf("user from proto: %w", err)
	}

	return &User{
		ID:        userID,
		FirstName: user.GetFirstName(),
		LastName:  user.GetLastName(),
		Nickname:  user.GetNickname(),
		Email:     user.GetEmail(),
		Country:   user.GetCountry(),
		CreatedAt: user.GetCreatedAt().AsTime(),
		UpdatedAt: user.GetUpdatedAt().AsTime(),
	}, nil
}
