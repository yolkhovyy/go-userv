package usergraphql

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/yolkhovyy/go-userv/contract/dto"
)

// GraphQL client for the user-graphql service.
type Client struct {
	baseURL    string
	httpClient *http.Client
	replacer   *strings.Replacer
}

// Creates a new GraphQL client.
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{},
		replacer:   strings.NewReplacer("\t", "", "\n", " "),
	}
}

// Executes a GraphQL request.
func (c *Client) execute(
	ctx context.Context,
	query string,
	variables map[string]interface{},
	result interface{},
) error {
	body, err := json.Marshal(map[string]interface{}{
		"query":     query,
		"variables": variables,
	})
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/v1/graphql", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http request: %w %d", ErrUnexpectedStatusCode, resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	return nil
}

// Creates a new user.
func (c *Client) Create(ctx context.Context, user dto.UserInput) (*dto.User, error) {
	query := c.replacer.Replace(
		`mutation CreateUser($input: UserCreate!) {
			create(input: $input) {
				id
				firstName
				lastName
				nickname
				email
				country
				createdAt
				updatedAt
			}
		}`)

	variables := map[string]interface{}{
		"input": user,
	}

	var response struct {
		Data struct {
			Create dto.User `json:"create"`
		} `json:"data"`
	}

	if err := c.execute(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return &response.Data.Create, nil
}

// Retrieves a user by ID.
func (c *Client) Get(ctx context.Context, userID uuid.UUID) (*dto.User, error) {
	query := `
        query GetUser($id: ID!) {
            user(id: $id) {
                id
                firstName
                lastName
                nickname
                email
                country
                createdAt
                updatedAt
            }
        }
    `
	variables := map[string]interface{}{
		"id": userID.String(),
	}

	var response struct {
		Data struct {
			User dto.User `json:"user"`
		} `json:"data"`
	}

	if err := c.execute(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return &response.Data.User, nil
}

// Retrieves a list of users.
func (c *Client) List(ctx context.Context, page, limit int, countryCode string) (*dto.UserList, error) {
	query := `
        query ListUsers($page: Int!, $limit: Int!, $country: String) {
            users(page: $page, limit: $limit, country: $country) {
                users {
                    id
                    firstName
                    lastName
                    nickname
                    email
                    country
                    createdAt
                    updatedAt
                }
                totalCount
                nextPage
            }
        }
    `
	variables := map[string]interface{}{
		"page":    page,
		"limit":   limit,
		"country": countryCode,
	}

	var response struct {
		Data struct {
			Users dto.UserList `json:"users"`
		} `json:"data"`
	}

	if err := c.execute(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return &response.Data.Users, nil
}

// Updates an existing user.
func (c *Client) Update(ctx context.Context, user dto.UserUpdate) (*dto.User, error) {
	query := `
        mutation UpdateUser($input: UserUpdate!) {
            update(input: $input) {
                id
                firstName
                lastName
                nickname
                email
                country
                createdAt
                updatedAt
            }
        }
    `
	variables := map[string]interface{}{
		"input": user,
	}

	var response struct {
		Data struct {
			Update dto.User `json:"update"`
		} `json:"data"`
	}

	if err := c.execute(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return &response.Data.Update, nil
}

// Deletes a user by ID.
func (c *Client) Delete(ctx context.Context, userID uuid.UUID) error {
	query := `
        mutation DeleteUser($id: ID!) {
            delete(id: $id)
        }
    `
	variables := map[string]interface{}{
		"id": userID.String(),
	}

	var response struct {
		Data struct {
			Delete bool `json:"delete"`
		} `json:"data"`
	}

	if err := c.execute(ctx, query, variables, &response); err != nil {
		return fmt.Errorf("graphql delete: %w", err)
	}

	if !response.Data.Delete {
		return fmt.Errorf("delete user: %w", ErrDeleteFailure)
	}

	return nil
}
