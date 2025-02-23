//go:build integration_tests

package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	userrest "github.com/yolkhovyy/user/client/user-rest"
	"github.com/yolkhovyy/user/internal/contract/domain"
)

type testCase struct {
	name     string
	testFunc func(t *testing.T, tcase testCase)
	numUsers int
	pageSize int
	country  string
}

type TestSuite struct {
	suite.Suite
	client *userrest.Client
	cases  []testCase
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) SetupSuite() {
	s.client = userrest.NewClient("http://localhost:8080")
	s.cases = []testCase{
		{name: "CreateUsers", testFunc: s.createUsers, numUsers: 100},
		{name: "GetUpdateUser", testFunc: s.getUpdateUser, numUsers: 100},
		{name: "ListUsers", testFunc: s.listUsers, numUsers: 100, pageSize: 12},
		{name: "DeleteAll", testFunc: s.deleteAllUsers, numUsers: 100},
	}
}

func (s *TestSuite) TearDownSuite() {
	// TODO: tear down suite.
}

func (s *TestSuite) SetupTest() {
	// TODO: setup test.
}

func (s *TestSuite) TearDownTest() {
	// TODO: tear down test.
}

func (s *TestSuite) TestREST() {
	for _, tc := range s.cases {
		s.T().Run(tc.name, func(t *testing.T) {
			tc.testFunc(t, tc)
		})
	}
}

func (s *TestSuite) createUsers(t *testing.T, tcase testCase) {
	for i := 0; i < tcase.numUsers; i++ {
		userID := uuid.New()
		userInput := domain.UserInput{
			FirstName: fmt.Sprintf("FirstName%d", i),
			LastName:  fmt.Sprintf("LastName%d", i),
			Nickname:  fmt.Sprintf("user-%s", userID.String()),
			Email:     fmt.Sprintf("user.%s@example.com", userID.String()),
			Country:   "US",
			Password:  fmt.Sprintf("securepassword%d", i),
		}

		createdUser, err := s.client.Create(context.Background(), userInput)
		require.NoError(t, err)
		assert.NotNil(t, createdUser)

		assert.NotEmpty(t, createdUser.ID)
		_, err = uuid.Parse(createdUser.ID.String())
		require.NoError(t, err)

		assert.Equal(t, userInput.FirstName, createdUser.FirstName)
		assert.Equal(t, userInput.LastName, createdUser.LastName)
		assert.Equal(t, userInput.Nickname, createdUser.Nickname)
		assert.Equal(t, userInput.Email, createdUser.Email)
		assert.Equal(t, userInput.Country, createdUser.Country)
		assert.NotZero(t, createdUser.CreatedAt)
		assert.NotZero(t, createdUser.UpdatedAt)
		assert.Equal(t, createdUser.CreatedAt, createdUser.UpdatedAt)
	}

	// Verify that we can list all created list.
	list, err := s.client.List(context.Background(), 1, tcase.numUsers+1, tcase.country)
	require.NoError(t, err)
	assert.Equal(t, tcase.numUsers, list.TotalCount)
	assert.Equal(t, -1, list.NextPage)
	assert.Equal(t, tcase.numUsers, len(list.Users))
}

func (s *TestSuite) getUpdateUser(t *testing.T, tcase testCase) {
	list, err := s.client.List(context.Background(), 1, 1, tcase.country)
	require.NoError(t, err)
	assert.Equal(t, tcase.numUsers, list.TotalCount)
	assert.Equal(t, 2, list.NextPage)
	assert.Equal(t, 1, len(list.Users))

	oneUser := list.Users[0]

	user, err := s.client.Get(context.Background(), oneUser.ID)
	require.NoError(t, err)
	assert.NotNil(t, user)

	assert.Equal(t, oneUser.ID, user.ID)
	assert.Equal(t, oneUser.FirstName, user.FirstName)
	assert.Equal(t, oneUser.LastName, user.LastName)
	assert.Equal(t, oneUser.Nickname, user.Nickname)
	assert.Equal(t, oneUser.Email, user.Email)
	assert.Equal(t, oneUser.Country, user.Country)
	assert.Equal(t, oneUser.CreatedAt, user.CreatedAt)
	assert.Equal(t, user.CreatedAt, user.UpdatedAt)

	userUpdate := domain.UserUpdate{
		ID:        user.ID,
		FirstName: "Bob",
		LastName:  "Smith",
		Nickname:  "bsmith",
		Email:     "bsmith@example.com",
		Country:   "GB",
		Password:  "newsecurepassword",
	}

	updatedUser, err := s.client.Update(context.Background(), userUpdate)
	require.NoError(t, err)
	assert.NotNil(t, user)

	assert.Equal(t, userUpdate.ID, updatedUser.ID)
	assert.Equal(t, userUpdate.FirstName, updatedUser.FirstName)
	assert.Equal(t, userUpdate.LastName, updatedUser.LastName)
	assert.Equal(t, userUpdate.Nickname, updatedUser.Nickname)
	assert.Equal(t, userUpdate.Email, updatedUser.Email)
	assert.Equal(t, userUpdate.Country, updatedUser.Country)
	assert.NotZero(t, updatedUser.CreatedAt)
	assert.Greater(t, updatedUser.UpdatedAt, oneUser.UpdatedAt)
}

func (s *TestSuite) listUsers(t *testing.T, tcase testCase) {
	lastPage := 1 + (tcase.numUsers+tcase.pageSize/2)/tcase.pageSize
	for page := 1; page <= lastPage; page++ {
		list, err := s.client.List(context.Background(), page, tcase.pageSize, "")
		require.NoError(t, err)
		assert.Equal(t, tcase.numUsers, list.TotalCount)
		if page < lastPage {
			assert.Equal(t, page+1, list.NextPage)
			assert.Equal(t, tcase.pageSize, len(list.Users))
		} else {
			assert.Equal(t, -1, list.NextPage)
			assert.Equal(t, tcase.numUsers%tcase.pageSize, len(list.Users))
		}
	}
}

func (s *TestSuite) deleteAllUsers(t *testing.T, tcase testCase) {
	list, err := s.client.List(context.Background(), 1, tcase.numUsers+1, tcase.country)
	require.NoError(t, err)

	for i := 0; i < len(list.Users); i++ {
		err := s.client.Delete(context.Background(), list.Users[i].ID)
		require.NoError(t, err)
	}
}
