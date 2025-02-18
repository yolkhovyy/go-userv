package usergrpc

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/yolkhovyy/go-userv/contract/dto"
	"github.com/yolkhovyy/go-userv/contract/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn   *grpc.ClientConn
	client proto.UserServiceClient
}

func NewClient(address string) (*Client, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %w", err)
	}

	return &Client{
		conn:   conn,
		client: proto.NewUserServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	err := c.conn.Close()
	if err != nil {
		return fmt.Errorf("failed to close gRPC connection: %w", err)
	}

	return nil
}

func (c *Client) Create(ctx context.Context, input dto.UserInput) (*dto.User, error) {
	req := &proto.UserInput{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Nickname:  input.Nickname,
		Password:  input.Password,
		Email:     input.Email,
		Country:   input.Country,
	}

	resp, err := c.client.Create(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	user, err := dtoUserFromProto(resp)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}

func (c *Client) Get(ctx context.Context, id uuid.UUID) (*dto.User, error) {
	req := &proto.UserID{
		Id: id.String(),
	}

	resp, err := c.client.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	user, err := dtoUserFromProto(resp)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	return user, nil
}

func (c *Client) List(ctx context.Context, page, limit int, country string) (*dto.UserList, error) {
	req := &proto.ListRequest{
		Page:    int64(page),
		Limit:   int64(limit),
		Country: country,
	}

	resp, err := c.client.List(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}

	users := make([]dto.User, len(resp.GetUsers()))

	for i, u := range resp.GetUsers() {
		user, err := dtoUserFromProto(u)
		if err != nil {
			return nil, fmt.Errorf("list users: %w", err)
		}

		users[i] = *user
	}

	userList := dto.UserList{
		TotalCount: int(resp.GetTotalCount()),
		NextPage:   int(resp.GetNextPage()),
		Users:      users,
	}

	return &userList, nil
}

func (c *Client) Update(ctx context.Context, update dto.UserUpdate) (*dto.User, error) {
	req := &proto.UserUpdate{
		Id:        update.ID.String(),
		FirstName: update.FirstName,
		LastName:  update.LastName,
		Nickname:  update.Nickname,
		Password:  update.Password,
		Email:     update.Email,
		Country:   update.Country,
	}

	resp, err := c.client.Update(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	user, err := dtoUserFromProto(resp)
	if err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	return user, nil
}

func (c *Client) Delete(ctx context.Context, id uuid.UUID) error {
	req := &proto.UserID{
		Id: id.String(),
	}

	_, err := c.client.Delete(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
