package gin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yolkhovyy/user/contract/dto"
)

func (u *Controller) health(gctx *gin.Context) {
	gctx.JSON(http.StatusOK, gin.H{"message": "healthy"})
}

func (u *Controller) create(gctx *gin.Context) {
	var userInput dto.UserInput
	if err := gctx.ShouldBindJSON(&userInput); err != nil {
		gctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	if err := userInput.ValidateOnCreate(); err != nil {
		gctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	createdUser, err := u.domain.Create(gctx.Request.Context(), dto.UserInputToDomain(userInput))
	if err != nil {
		gctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	gctx.JSON(http.StatusCreated, createdUser)
}

func (u *Controller) update(gctx *gin.Context) {
	var user dto.UserInput
	if err := gctx.ShouldBindJSON(&user); err != nil {
		gctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	userID, err := uuid.Parse(gctx.Param("id"))
	if err != nil {
		gctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	if user.ID == uuid.Nil {
		user.ID = userID
	} else if user.ID != userID {
		gctx.JSON(http.StatusBadRequest, gin.H{"error": ErrUUIDConflict.Error()})

		return
	}

	if err := user.ValidateOnUpdate(); err != nil {
		gctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	updatedUser, err := u.domain.Update(gctx.Request.Context(), dto.UserInputToDomain(user))
	if err != nil {
		gctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	gctx.JSON(http.StatusOK, updatedUser)
}

func (u *Controller) get(gctx *gin.Context) {
	userID, err := uuid.Parse(gctx.Param("id"))
	if err != nil {
		gctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	user, err := u.domain.Get(gctx.Request.Context(), userID)
	if err != nil {
		gctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

		return
	}

	gctx.JSON(http.StatusOK, dto.UserFromDomain(*user))
}

func (u *Controller) list(gctx *gin.Context) {
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

	list, err := u.domain.List(gctx.Request.Context(), page, limit, country)
	if err != nil {
		gctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	gctx.JSON(http.StatusOK, dto.UsersFromDomain(*list))
}

func (u *Controller) delete(gctx *gin.Context) {
	userID, err := uuid.Parse(gctx.Param("id"))
	if err != nil {
		gctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	err = u.domain.Delete(gctx.Request.Context(), userID)
	if err != nil {
		gctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	gctx.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}
