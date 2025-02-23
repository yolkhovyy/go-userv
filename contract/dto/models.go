package dto

import (
	"github.com/yolkhovyy/user/internal/contract/domain"
)

// TODO: Replace these by dedicated DTO types if needed.

type User domain.User

type UserInput domain.UserInput

type UserUpdate domain.UserUpdate

type UserList domain.UserList

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
	return UserList(userList)
}

func UserListToDomain(userList UserList) domain.UserList {
	return domain.UserList(userList)
}
