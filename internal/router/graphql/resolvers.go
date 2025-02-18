package graphql

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
	"github.com/yolkhovyy/go-userv/internal/contract/domain"
)

var (
	ErrTypeAssertion = errors.New("type assertion")
	ErrMissingKey    = errors.New("missing key")
)

func (c *Controller) user() graphql.FieldResolveFn {
	return withLogging(func(params graphql.ResolveParams) (any, error) {
		var input inputMap = params.Args

		userID, err := input.uuidValue("id")
		if err != nil {
			return nil, fmt.Errorf("user resolver: %w", ErrMissingKey)
		}

		user, err := c.domain.Get(params.Context, userID)
		if err != nil {
			return nil, fmt.Errorf("user resolver: %w", err)
		}

		return user, nil
	})
}

func (c *Controller) users() graphql.FieldResolveFn {
	return withLogging(func(params graphql.ResolveParams) (any, error) {
		const (
			defaultPage  = 1
			defaultLimit = 10
		)

		var input inputMap = params.Args

		page := defaultPage
		if value, ok := input.intValue("page"); ok {
			page = value
		}

		limit := defaultLimit
		if value, ok := input.intValue("limit"); ok {
			limit = value
		}

		var country string
		if value, ok := input.stringValue("country"); ok {
			country = value
		}

		list, err := c.domain.List(params.Context, page, limit, country)
		if err != nil {
			return nil, fmt.Errorf("users resolver: %w", err)
		}

		return list, nil
	})
}

func (c *Controller) create() graphql.FieldResolveFn {
	return withLogging(func(params graphql.ResolveParams) (any, error) {
		user, err := userInput(params)
		if err != nil {
			return nil, fmt.Errorf("create user resolver: %w", err)
		}

		createdUser, err := c.domain.Create(params.Context, *user)
		if err != nil {
			return nil, fmt.Errorf("create user resolver: %w", err)
		}

		return createdUser, nil
	})
}

func (c *Controller) update() graphql.FieldResolveFn {
	return withLogging(func(params graphql.ResolveParams) (any, error) {
		user, err := userUpdate(params)
		if err != nil {
			return nil, fmt.Errorf("update user input: %w", err)
		}

		updatedUser, err := c.domain.Update(params.Context, *user)
		if err != nil {
			return nil, fmt.Errorf("update user resolver: %w", err)
		}

		return updatedUser, nil
	})
}

func (c *Controller) delete() graphql.FieldResolveFn {
	return withLogging(func(params graphql.ResolveParams) (any, error) {
		var input inputMap = params.Args

		userID, err := input.uuidValue("id")
		if err != nil {
			return nil, fmt.Errorf("delete user resolver: %w", ErrMissingKey)
		}

		if err := c.domain.Delete(params.Context, userID); err != nil {
			return false, fmt.Errorf("delete user resolver: %w", err)
		}

		return true, nil
	})
}

func foo(params graphql.ResolveParams) (inputMap, error) {
	inputRaw, exists := params.Args["input"]
	if !exists {
		return nil, fmt.Errorf("user input raw: %w", ErrTypeAssertion)
	}

	inputAny, exists := inputRaw.(map[string]any)
	if !exists {
		return nil, fmt.Errorf("user input map: %w", ErrTypeAssertion)
	}

	return inputMap(inputAny), nil
}

func userInput(params graphql.ResolveParams) (*domain.UserInput, error) {
	input, err := foo(params)
	if err != nil {
		return nil, fmt.Errorf("user input: %w", err)
	}

	userInput := &domain.UserInput{}

	if value, ok := input.stringValue("firstName"); ok {
		userInput.FirstName = value
	}

	if value, ok := input.stringValue("lastName"); ok {
		userInput.LastName = value
	}

	if value, ok := input.stringValue("nickname"); ok {
		userInput.Nickname = value
	}

	if value, ok := input.stringValue("email"); ok {
		userInput.Email = value
	}

	if value, ok := input.stringValue("country"); ok {
		userInput.Country = value
	}

	if value, ok := input.stringValue("password"); ok {
		userInput.Password = value
	}

	return userInput, nil
}

func userUpdate(params graphql.ResolveParams) (*domain.UserUpdate, error) {
	input, err := foo(params)
	if err != nil {
		return nil, fmt.Errorf("user update: %w", err)
	}

	userUpdate := &domain.UserUpdate{}

	value, err := input.uuidValue("id")
	if err != nil {
		return nil, fmt.Errorf("user update %w", err)
	}

	userUpdate.ID = value

	if value, ok := input.stringValue("firstName"); ok {
		userUpdate.FirstName = value
	}

	if value, ok := input.stringValue("lastName"); ok {
		userUpdate.LastName = value
	}

	if value, ok := input.stringValue("nickname"); ok {
		userUpdate.Nickname = value
	}

	if value, ok := input.stringValue("email"); ok {
		userUpdate.Email = value
	}

	if value, ok := input.stringValue("country"); ok {
		userUpdate.Country = value
	}

	if value, ok := input.stringValue("password"); ok {
		userUpdate.Password = value
	}

	return userUpdate, nil
}

type inputMap map[string]any

func (im inputMap) stringValue(key string) (string, bool) {
	input, exists := map[string]any(im)[key]
	stringValue, valid := input.(string)

	return stringValue, exists && valid
}

func (im inputMap) intValue(key string) (int, bool) {
	input, exists := map[string]any(im)[key]
	intValue, valid := input.(int)

	return intValue, exists && valid
}

func (im inputMap) uuidValue(key string) (uuid.UUID, error) {
	value, exists := map[string]any(im)[key]
	if !exists {
		return uuid.Nil, fmt.Errorf("resolvercuuid: %w", ErrMissingKey)
	}

	valueStr, ok := value.(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("resolver cuuid: %w", ErrTypeAssertion)
	}

	uuidValue, err := uuid.Parse(valueStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("resolver cuuid: %w", err)
	}

	return uuidValue, nil
}
