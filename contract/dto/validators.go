package dto

import (
	"fmt"
	"regexp"
)

func (u *UserInput) ValidateOnCreate() error {
	if u.FirstName == "" {
		return fmt.Errorf("invalid: %w %s", ErrFirstName, u.FirstName)
	}

	if u.LastName == "" {
		return fmt.Errorf("invalid: %w %s", ErrLastName, u.LastName)
	}

	if u.Nickname == "" {
		return fmt.Errorf("invalid: %w", ErrNickname)
	}

	if err := ValidateEmail(u.Email); err != nil {
		return fmt.Errorf("invalid: %w", err)
	}

	if err := ValidateCountryCode(u.Country); err != nil {
		return fmt.Errorf("invalid: %w", err)
	}

	if err := ValidatePassword(u.Password); err != nil {
		return fmt.Errorf("invalid: %w", err)
	}

	return nil
}

func (u *UserInput) ValidateOnUpdate() error {
	if u.Email != "" {
		if err := ValidateEmail(u.Email); err != nil {
			return fmt.Errorf("invalid: %w", err)
		}
	}

	if u.Country != "" {
		if err := ValidateCountryCode(u.Country); err != nil {
			return fmt.Errorf("invalid: %w", err)
		}
	}

	if u.Password != "" {
		if err := ValidatePassword(u.Password); err != nil {
			return fmt.Errorf("invalid: %w", err)
		}
	}

	return nil
}

func ValidateCountryCode(countryCode string) error {
	if countryCode == "" {
		return fmt.Errorf("code: %w", ErrCountryCode)
	}

	if match, _ := regexp.MatchString("^[A-Z]{2}$", countryCode); !match {
		return fmt.Errorf("country code: %w %s", ErrCountryCode, countryCode)
	}

	return nil
}

func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email: %w %s", ErrEmail, email)
	}

	if match, _ := regexp.MatchString(`^[^@\s]+@[^@\s]+\.[^@\s]+$`, email); !match {
		return fmt.Errorf("email: %w %s", ErrEmail, email)
	}

	return nil
}

func ValidatePassword(_ string) error {
	// TODO: Make strong password validation.
	return nil
}
