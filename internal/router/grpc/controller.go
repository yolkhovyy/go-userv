package grpc

import (
	usergrpc "github.com/yolkhovyy/user/contract/proto"
	"github.com/yolkhovyy/user/internal/contract/domain"
)

type Controller struct {
	domain domain.Contract
	usergrpc.UnimplementedUserServiceServer
}

func New(_ Config, domain domain.Contract) *Controller {
	user := Controller{
		domain: domain,
	}

	return &user
}
