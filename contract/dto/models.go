package dto

import (
	"github.com/yolkhovyy/go-userv/internal/contract/domain"
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

func UsersFromDomain(domainUsers []domain.User) []User {
	users := make([]User, len(domainUsers))

	for i, u := range domainUsers {
		du := UserFromDomain(u)
		users[i] = du
	}

	return users
}
