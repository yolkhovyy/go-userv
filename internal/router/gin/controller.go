package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/yolkhovyy/user/internal/contract/domain"
)

type Controller struct {
	domain  domain.Contract
	handler *gin.Engine
}

func New(config Config, domain domain.Contract) *Controller {
	controller := Controller{
		domain: domain,
	}

	// Gin routing engine.
	gin.SetMode(config.Mode)
	engine := gin.New()
	engine.RedirectTrailingSlash = false
	engine.Use(gin.Recovery(), Logger())

	// Health check.
	engine.GET("/health", controller.health)

	// API endpoints.
	group := engine.Group("/api/v1")
	{
		group.POST("/user", controller.create)
		group.GET("/user/:id", controller.get)
		group.GET("/users", controller.list)
		group.PUT("/user/:id", controller.update)
		group.DELETE("/user/:id", controller.delete)
	}

	controller.handler = engine

	return &controller
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
