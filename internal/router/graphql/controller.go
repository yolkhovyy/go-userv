package graphql

import (
	"fmt"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/rs/zerolog/log"
	"github.com/yolkhovyy/user/internal/contract/domain"
)

type Controller struct {
	domain  domain.Contract
	handler *handler.Handler
}

func New(_ Config, domain domain.Contract) (*Controller, error) {
	controller := Controller{
		domain: domain,
	}

	schema, err := graphql.NewSchema(controller.schemaConfig())
	if err != nil {
		return nil, fmt.Errorf("new graphql router: %w", err)
	}

	controller.handler = handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	http.Handle("/api/v1/graphql", controller.handler)

	return &controller, nil
}

func (c *Controller) Close() error {
	log.Debug().Msg("router closing")

	if err := c.domain.Close(); err != nil {
		log.Error().Err(err).Msg("router domain close")
	}

	log.Trace().Msg("router closed")

	return nil
}

func (c *Controller) Handler() http.Handler {
	return c.handler
}

func (c *Controller) Schema() *graphql.Schema {
	return c.handler.Schema
}
