package grpc

import (
	"github.com/yolkhovyy/go-userv/contract/proto"
	"github.com/yolkhovyy/go-userv/internal/contract/domain"
)

type Controller struct {
	domain domain.Contract
	proto.UnimplementedUserServiceServer
}

func New(_ Config, domain domain.Contract) *Controller {
	user := Controller{
		domain: domain,
	}

	return &user
}
