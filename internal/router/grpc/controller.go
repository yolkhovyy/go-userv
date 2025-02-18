package grpc

import (
	"github.com/rs/zerolog/log"
	userpb "github.com/yolkhovyy/user/contract/proto"
	"github.com/yolkhovyy/user/internal/contract/domain"
)

type Controller struct {
	domain domain.Contract
	userpb.UnimplementedUserServiceServer
}

func New(_ Config, domain domain.Contract) *Controller {
	user := Controller{
		domain: domain,
	}

	return &user
}

func (c *Controller) Close() error {
	log.Debug().Msg("router closing")

	if err := c.domain.Close(); err != nil {
		log.Error().Err(err).Msg("router domain close")
	}

	log.Trace().Msg("router closed")

	return nil
}
