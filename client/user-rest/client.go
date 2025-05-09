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
	"github.com/yolkhovyy/go-otelw/otelw/tracew"
	"github.com/yolkhovyy/go-userv/client/internal/oteld"
	"github.com/yolkhovyy/go-userv/contract/dto"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// REST client for the user-rest service.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// Creates a new REST client.
func NewClient(baseURL string) *Client {
	transport := otelhttp.NewTransport(http.DefaultTransport)

	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{Transport: transport},
	}
}

// Creates a new user.
func (c *Client) Create(ctx context.Context, user dto.UserInput) (*dto.User, error) {
	var err error

	ctx, span := tracew.Start(ctx, "user rest client", "create")
	defer func() { span.End(err) }()

	ctx, err = oteld.ContextWithBaggageUserName(ctx, user.FirstName, user.LastName)
	if err != nil {
		return nil, fmt.Errorf("user name baggage: %w", err)
	}

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
	var err error

	ctx, span := tracew.Start(ctx, "user rest client", "get")
	defer func() { span.End(err) }()

	ctx, err = oteld.ContextWithBaggageUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user id baggage: %w", err)
	}

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
	var err error

	ctx, span := tracew.Start(ctx, "user rest client", "list")
	defer func() { span.End(err) }()

	ctx, err = oteld.ContextWithBaggagePage(ctx, page, limit, countryCode)
	if err != nil {
		return nil, fmt.Errorf("page baggage: %w", err)
	}

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
	var err error

	ctx, span := tracew.Start(ctx, "user rest client", "update")
	defer func() { span.End(err) }()

	ctx, err = oteld.ContextWithBaggageUserID(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("user id baggage: %w", err)
	}

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
	var err error

	ctx, span := tracew.Start(ctx, "user rest client", "delete")
	defer func() { span.End(err) }()

	ctx, err = oteld.ContextWithBaggageUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user id baggage: %w", err)
	}

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
