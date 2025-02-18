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
	user := Controller{
		domain: domain,
	}

	// Gin routing engine.
	gin.SetMode(config.Mode)
	engine := gin.New()
	engine.RedirectTrailingSlash = false
	engine.Use(gin.Recovery(), Logger())

	// Health check.
	engine.GET("/health", user.health)

	// API endpoints.
	group := engine.Group("/api/v1")
	{
		group.POST("/user", user.create)
		group.GET("/user/:id", user.get)
		group.GET("/users", user.list)
		group.PUT("/user/:id", user.update)
		group.DELETE("/user/:id", user.delete)
	}

	user.handler = engine

	return &user
}

func (u *Controller) Close() error {
	log.Debug().Msg("router closing")

	if err := u.domain.Close(); err != nil {
		log.Error().Err(err).Msg("router domain close")
	}

	log.Trace().Msg("router closed")

	return nil
}

func (u *Controller) Handler() http.Handler {
	return u.handler
}
