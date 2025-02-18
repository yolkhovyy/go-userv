package userrest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/google/uuid"
	"github.com/yolkhovyy/go-userv/contract/dto"
)

// REST client for the user-rest service.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// Creates a new REST client.
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

// Creates a new user.
func (c *Client) Create(ctx context.Context, user dto.UserInput) (*dto.User, error) {
	body, err := json.Marshal(user)
	if err != nil {
		return nil, fmt.Errorf("marshal user input: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/v1/user", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("create: %w %d", ErrUnexpectedStatusCode, resp.StatusCode)
	}

	var createdUser dto.User
	if err := json.NewDecoder(resp.Body).Decode(&createdUser); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &createdUser, nil
}

// Retrieves a user by ID.
func (c *Client) Get(ctx context.Context, userID uuid.UUID) (*dto.User, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/api/v1/user/%s", c.baseURL, userID), nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get: %w %d", ErrUnexpectedStatusCode, resp.StatusCode)
	}

	var user dto.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &user, nil
}

// Retrieves a list of users.
func (c *Client) List(ctx context.Context, page, limit int, countryCode string) (*dto.UserList, error) {
	query := url.Values{}
	query.Set("page", strconv.Itoa(page))
	query.Set("limit", strconv.Itoa(limit))

	if countryCode != "" {
		query.Set("country", countryCode)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		fmt.Sprintf("%s/api/v1/users?%s", c.baseURL, query.Encode()), nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("list: %w %d", ErrUnexpectedStatusCode, resp.StatusCode)
	}

	var userList dto.UserList
	if err := json.NewDecoder(resp.Body).Decode(&userList); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &userList, nil
}

// Updates an existing user.
func (c *Client) Update(ctx context.Context, user dto.UserUpdate) (*dto.User, error) {
	body, err := json.Marshal(user)
	if err != nil {
		return nil, fmt.Errorf("marshal user input: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut,
		fmt.Sprintf("%s/api/v1/user/%s", c.baseURL, user.ID.String()), bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("update: %w %d", ErrUnexpectedStatusCode, resp.StatusCode)
	}

	var updatedUser dto.User
	if err := json.NewDecoder(resp.Body).Decode(&updatedUser); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &updatedUser, nil
}

// Deletes a user by ID.
func (c *Client) Delete(ctx context.Context, userID uuid.UUID) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete,
		fmt.Sprintf("%s/api/v1/user/%s", c.baseURL, userID), nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)

		return fmt.Errorf("delete: %w %d, body: %s", ErrUnexpectedStatusCode, resp.StatusCode, string(body))
	}

	return nil
}
