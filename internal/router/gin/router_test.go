package gin

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yolkhovyy/user/internal/contract/domain"
	"github.com/yolkhovyy/user/internal/contract/storage"
)

func TestUser_create(t *testing.T) {
	t.Parallel()
	t.Run("Create user", func(t *testing.T) {
		t.Parallel()
		gin.SetMode(gin.TestMode)

		userInput := domain.UserInput(
			storage.UserInput{
				FirstName: "John",
				LastName:  "Doe",
				Nickname:  "john.doe",
				Email:     "john.doe@example.com",
				Country:   "GB",
				Password:  "securepassword",
			},
		)

		userInputBytes, err := json.Marshal(userInput)
		require.NoError(t, err)

		createdUser := domain.User{
			ID:        uuid.New(),
			FirstName: userInput.FirstName,
			LastName:  userInput.LastName,
			Nickname:  userInput.Nickname,
			Email:     userInput.Email,
			Country:   userInput.Country,
		}
		createdUserBytes, err := json.Marshal(createdUser)
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		gctx, _ := gin.CreateTestContext(rec)
		gctx.Request = httptest.NewRequest(http.MethodPost, "/user", bytes.NewReader(userInputBytes))

		mockDomain := domain.MockContract{}
		mockDomain.EXPECT().
			Create(gctx.Request.Context(), userInput).
			Return(&createdUser, nil)

		router := Controller{
			domain: &mockDomain,
		}

		router.create(gctx)
		assert.Equal(t, http.StatusCreated, rec.Code)

		body, err := io.ReadAll(rec.Body)
		require.NoError(t, err)
		assert.JSONEq(t, string(createdUserBytes), string(body))
	})
}

func TestUser_update(t *testing.T) {
	t.Parallel()
	t.Run("Update user", func(t *testing.T) {
		t.Parallel()
		gin.SetMode(gin.TestMode)

		userUpdate := domain.UserUpdate(
			storage.UserUpdate{
				ID:        uuid.New(),
				FirstName: "John",
				LastName:  "Doe",
				Nickname:  "john.doe",
				Email:     "john.doe@example.com",
				Country:   "GB",
				Password:  "securepassword",
			},
		)

		userUpdateBytes, err := json.Marshal(userUpdate)
		require.NoError(t, err)

		updatedUser := domain.User{
			ID:        userUpdate.ID,
			FirstName: userUpdate.FirstName,
			LastName:  userUpdate.LastName,
			Nickname:  userUpdate.Nickname,
			Email:     userUpdate.Email,
			Country:   userUpdate.Country,
		}
		updatedUserBytes, err := json.Marshal(updatedUser)
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		gctx, _ := gin.CreateTestContext(rec)
		gctx.Request = httptest.NewRequest(http.MethodPut, "/user/"+updatedUser.ID.String(), bytes.NewReader(userUpdateBytes))
		gctx.Params = gin.Params{gin.Param{Key: "id", Value: updatedUser.ID.String()}}

		mockDomain := domain.MockContract{}
		mockDomain.EXPECT().
			Update(gctx.Request.Context(), userUpdate).
			Return(&updatedUser, nil)

		router := Controller{
			domain: &mockDomain,
		}

		router.update(gctx)
		assert.Equal(t, http.StatusOK, rec.Code)

		body, err := io.ReadAll(rec.Body)
		require.NoError(t, err)
		assert.JSONEq(t, string(updatedUserBytes), string(body))
	})
}

func TestUser_get(t *testing.T) {
	t.Parallel()
	t.Run("Get user", func(t *testing.T) {
		t.Parallel()
		gin.SetMode(gin.TestMode)

		gotUser := domain.User(
			storage.User{
				ID:        uuid.New(),
				FirstName: "John",
				LastName:  "Doe",
				Nickname:  "john.doe",
				Email:     "john.doe@example.com",
				Country:   "GB",
			},
		)

		gotUserBytes, err := json.Marshal(gotUser)
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		gctx, _ := gin.CreateTestContext(rec)
		gctx.Request = httptest.NewRequest(http.MethodGet, "/user/"+gotUser.ID.String(), nil)
		gctx.Params = gin.Params{gin.Param{Key: "id", Value: gotUser.ID.String()}}

		mockDomain := domain.MockContract{}
		mockDomain.EXPECT().
			Get(gctx.Request.Context(), gotUser.ID).
			Return(&gotUser, nil)

		router := Controller{
			domain: &mockDomain,
		}

		router.get(gctx)
		assert.Equal(t, http.StatusOK, rec.Code)

		body, err := io.ReadAll(rec.Body)
		require.NoError(t, err)
		assert.JSONEq(t, string(gotUserBytes), string(body))
	})
}

func TestUser_list(t *testing.T) {
	t.Parallel()
	t.Run("List users", func(t *testing.T) {
		t.Parallel()
		gin.SetMode(gin.TestMode)

		const (
			numUsers = 42
			page     = 1
			limit    = 43
			country  = "GB"
		)

		listUsers := domain.UserList(
			storage.UserList{
				Users:      make([]storage.User, 0, numUsers),
				TotalCount: numUsers,
				NextPage:   -1,
			},
		)

		for range numUsers {
			user := storage.User{
				ID:        uuid.New(),
				FirstName: "John",
				LastName:  "Doe",
				Nickname:  "john.doe",
				Email:     "john.doe@example.com",
				Country:   "GB",
			}
			listUsers.Users = append(listUsers.Users, user)
		}

		listUsersBytes, err := json.Marshal(listUsers)
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		gctx, _ := gin.CreateTestContext(rec)
		request := httptest.NewRequest(http.MethodGet, "/users", nil)
		query := request.URL.Query()
		query.Set("page", strconv.Itoa(page))
		query.Set("limit", strconv.Itoa(limit))
		query.Set("country", country)
		request.URL.RawQuery = query.Encode()
		gctx.Request = request

		mockDomain := domain.MockContract{}
		mockDomain.EXPECT().
			List(gctx.Request.Context(), page, limit, country).
			Return(&listUsers, nil)

		router := Controller{
			domain: &mockDomain,
		}

		router.list(gctx)
		assert.Equal(t, http.StatusOK, rec.Code)

		body, err := io.ReadAll(rec.Body)
		require.NoError(t, err)
		assert.JSONEq(t, string(listUsersBytes), string(body))
	})
}

func TestUser_delete(t *testing.T) {
	t.Parallel()
	t.Run("Delete user", func(t *testing.T) {
		t.Parallel()
		gin.SetMode(gin.TestMode)

		userID := uuid.New()
		rec := httptest.NewRecorder()

		gctx, _ := gin.CreateTestContext(rec)
		gctx.Request = httptest.NewRequest(http.MethodDelete, "/user/"+userID.String(), nil)
		gctx.Params = gin.Params{gin.Param{Key: "id", Value: userID.String()}}

		mockDomain := domain.MockContract{}
		mockDomain.EXPECT().Delete(gctx.Request.Context(), userID).Return(nil)

		controller := Controller{
			domain: &mockDomain,
		}

		controller.delete(gctx)
		// FIXME: It must return http.StatusNoContent.
		assert.Equal(t, http.StatusOK, rec.Code)

		body, err := io.ReadAll(rec.Body)
		require.NoError(t, err)
		assert.Empty(t, body)
	})
}
