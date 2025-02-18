package dto

import "errors"

var (
	ErrFirstName   = errors.New("first name")
	ErrLastName    = errors.New("last name")
	ErrNickname    = errors.New("nickname")
	ErrCountryCode = errors.New("country code")
	ErrEmail       = errors.New("email")
	ErrPassword    = errors.New("password")
)
