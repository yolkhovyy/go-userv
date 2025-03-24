package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yolkhovyy/go-userv/internal/contract/domain"
)

type Controller struct {
	domain  domain.Contract
	handler *gin.Engine
}

func New(config Config, domain domain.Contract, handlers ...gin.HandlerFunc) *Controller {
	controller := Controller{
		domain: domain,
	}

	// Gin routing engine.
	gin.SetMode(config.Mode)
	engine := gin.New()
	engine.RedirectTrailingSlash = false
	engine.Use(gin.Recovery(), Logger())
	engine.Use(handlers...)

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

func (c *Controller) Handler() http.Handler {
	return c.handler
}
