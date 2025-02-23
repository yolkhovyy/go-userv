package gin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yolkhovyy/user/contract/dto"
)

func (c *Controller) health(gctx *gin.Context) {
	gctx.JSON(http.StatusOK, gin.H{"message": "healthy"})
}

func (c *Controller) create(gctx *gin.Context) {
	var userInput dto.UserInput
	if err := gctx.ShouldBindJSON(&userInput); err != nil {
		gctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	if err := userInput.Validate(); err != nil {
		gctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	createdUser, err := c.domain.Create(gctx.Request.Context(), dto.UserInputToDomain(userInput))
	if err != nil {
		gctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	gctx.JSON(http.StatusCreated, createdUser)
}

func (c *Controller) update(gctx *gin.Context) {
	userID, err := uuid.Parse(gctx.Param("id"))
	if err != nil {
		gctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	var user dto.UserUpdate
	if err := gctx.ShouldBindJSON(&user); err != nil {
		gctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	user.ID = userID

	if err := user.Validate(); err != nil {
		gctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	updatedUser, err := c.domain.Update(gctx.Request.Context(), dto.UserUpdateToDomain(user))
	if err != nil {
		gctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	gctx.JSON(http.StatusOK, updatedUser)
}

func (c *Controller) get(gctx *gin.Context) {
	userID, err := uuid.Parse(gctx.Param("id"))
	if err != nil {
		gctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	user, err := c.domain.Get(gctx.Request.Context(), userID)
	if err != nil {
		gctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

		return
	}

	gctx.JSON(http.StatusOK, dto.UserFromDomain(*user))
}

func (c *Controller) list(gctx *gin.Context) {
	const (
		defaultPage  = 1
		defaultLimit = 10
	)

	page, err := strconv.Atoi(gctx.DefaultQuery("page", strconv.Itoa(defaultPage)))
	if err != nil || page <= 0 {
		gctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	limit, err := strconv.Atoi(gctx.DefaultQuery("limit", strconv.Itoa(defaultLimit)))
	if err != nil || page <= 0 {
		gctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	country := gctx.Query("country")
	if err := dto.ValidateCountryCode(country); err != nil && country != "" {
		gctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	list, err := c.domain.List(gctx.Request.Context(), page, limit, country)
	if err != nil {
		gctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	gctx.JSON(http.StatusOK, dto.UserListFromDomain(*list))
}

func (c *Controller) delete(gctx *gin.Context) {
	userID, err := uuid.Parse(gctx.Param("id"))
	if err != nil {
		gctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	err = c.domain.Delete(gctx.Request.Context(), userID)
	if err != nil {
		gctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	gctx.Status(http.StatusNoContent)
}
