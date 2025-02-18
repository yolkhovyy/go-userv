package dto

import (
	"github.com/yolkhovyy/user/internal/contract/domain"
)

// TODO: Replace these by dedicated DTO types if needed.

type User domain.User

type UserInput domain.UserInput

type Users domain.Users

func UserFromDomain(user domain.User) User {
	return User(user)
}

func UserInputToDomain(userInput UserInput) domain.UserInput {
	return domain.UserInput(userInput)
}

func UsersFromDomain(users domain.Users) Users {
	return Users(users)
}
