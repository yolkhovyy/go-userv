package usergrpc

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/yolkhovyy/go-userv/contract/dto"
	"github.com/yolkhovyy/go-userv/contract/proto"
)

func dtoUserFromProto(user *proto.User) (*dto.User, error) {
	userID, err := uuid.Parse(user.GetId())
	if err != nil {
		return nil, fmt.Errorf("user from proto: %w", err)
	}

	return &dto.User{
		ID:        userID,
		FirstName: user.GetFirstName(),
		LastName:  user.GetLastName(),
		Nickname:  user.GetNickname(),
		Email:     user.GetEmail(),
		Country:   user.GetCountry(),
		CreatedAt: user.GetCreatedAt().AsTime(),
		UpdatedAt: user.GetUpdatedAt().AsTime(),
	}, nil
}
