package grpc

import (
	"context"
	"fmt"
	"log/slog"

	empty "github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"github.com/yolkhovyy/go-otelw/otelw/slogw"
	"github.com/yolkhovyy/go-userv/contract/dto"
	"github.com/yolkhovyy/go-userv/contract/proto"
	"github.com/yolkhovyy/go-userv/internal/contract/domain"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (c *Controller) Create(ctx context.Context, req *proto.UserInput) (*proto.User, error) {
	userInput := dto.UserInput{
		FirstName: req.GetFirstName(),
		LastName:  req.GetLastName(),
		Nickname:  req.GetNickname(),
		Email:     req.GetEmail(),
		Country:   req.GetCountry(),
		Password:  req.GetPassword(),
	}

	if err := userInput.Validate(); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	createdUser, err := c.domain.Create(ctx, dto.UserInputToDomain(userInput))
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return userToProto(createdUser), nil
}

func (c *Controller) Get(ctx context.Context, req *proto.UserID) (*proto.User, error) {
	var err error

	var userID uuid.UUID

	if req.GetId() != "" {
		userID, err = uuid.Parse(req.GetId())
		if err != nil {
			return nil, fmt.Errorf("user id: %w", err)
		}
	}

	user, err := c.domain.Get(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user:%w", err)
	}

	return userToProto(user), nil
}

func (c *Controller) List(ctx context.Context, req *proto.ListRequest) (*proto.Users, error) {
	logger := slogw.DefaultLogger()

	list, err := c.domain.List(ctx, int(req.GetPage()), int(req.GetLimit()), req.GetCountry())
	if err != nil {
		logger.ErrorContext(ctx, "domain",
			slog.String("list", err.Error()),
		)

		return nil, fmt.Errorf("list: %w", err)
	}

	protoUsers := make([]*proto.User, len(list.Users))

	for i, user := range list.Users {
		protoUsers[i] = userToProto(&user)
	}

	return &proto.Users{
		Users:      protoUsers,
		TotalCount: int64(list.TotalCount),
		NextPage:   int64(list.NextPage),
	}, nil
}

func (c *Controller) Update(ctx context.Context, req *proto.UserUpdate) (*proto.User, error) {
	var err error

	var userID uuid.UUID

	userID, err = uuid.Parse(req.GetId())
	if err != nil {
		return nil, fmt.Errorf("user id: %w", err)
	}

	userUpdate := dto.UserUpdate{
		ID:        userID,
		FirstName: req.GetFirstName(),
		LastName:  req.GetLastName(),
		Nickname:  req.GetNickname(),
		Email:     req.GetEmail(),
		Country:   req.GetCountry(),
		Password:  req.GetPassword(),
	}

	if err := userUpdate.Validate(); err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	updatedUser, err := c.domain.Update(ctx, dto.UserUpdateToDomain(userUpdate))
	if err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	return userToProto(updatedUser), nil
}

func (c *Controller) Delete(ctx context.Context, req *proto.UserID) (*empty.Empty, error) {
	var err error

	var userID uuid.UUID

	userID, err = uuid.Parse(req.GetId())
	if err != nil {
		return nil, fmt.Errorf("user id: %w", err)
	}

	err = c.domain.Delete(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("delete user: %w", err)
	}

	return &empty.Empty{}, nil
}

func userToProto(user *domain.User) *proto.User {
	return &proto.User{
		Id:        user.ID.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Nickname:  user.Nickname,
		Email:     user.Email,
		Country:   user.Country,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}
