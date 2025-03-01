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
	grpcclient "github.com/yolkhovyy/go-userv/client/user-grpc"
	"github.com/yolkhovyy/go-userv/contract/dto"
)

type testCaseGRPC struct {
	name     string
	testFunc func(t *testing.T, tcase testCaseGRPC)
	numUsers int
	pageSize int
}

type testSuiteGRPC struct {
	suite.Suite
	client        *grpcclient.Client
	testCases     []testCaseGRPC
	createCountry string
	updateCountry string
}

func TestGRPCSuiteRun(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(testSuiteGRPC))
}

func (s *testSuiteGRPC) SetupSuite() {
	var err error
	s.client, err = grpcclient.NewClient("localhost:50051")
	require.Nil(s.T(), err)
	s.testCases = []testCaseGRPC{
		{name: "CreateUsers", testFunc: s.createUsers, numUsers: 100},
		{name: "GetUpdateDeleteUser", testFunc: s.getUpdateDeleteUser, numUsers: 100},
		{name: "ListUsers", testFunc: s.listUsers, numUsers: 100, pageSize: 12},
		{name: "DeleteAll", testFunc: s.deleteAllUsers, numUsers: 100},
	}
	s.createCountry = "ZB"
	s.updateCountry = "ZC"
}

func (s *testSuiteGRPC) TearDownSuite() {
	// TODO: tear down suite.
}

func (s *testSuiteGRPC) SetupTest() {
	// TODO: setup test.
}

func (s *testSuiteGRPC) TearDownTest() {
	// TODO: tear down test.
}

func (s *testSuiteGRPC) TestGRPC() {
	for _, tc := range s.testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			tc.testFunc(t, tc)
		})
	}
}

func (s *testSuiteGRPC) createUsers(t *testing.T, tcase testCaseGRPC) {
	for i := 0; i < tcase.numUsers; i++ {
		userInput := dto.UserInput{
			FirstName: fmt.Sprintf("gRPC"),
			LastName:  fmt.Sprintf("User, %d", i),
			Nickname:  fmt.Sprintf("grpcuser%d", i),
			Email:     fmt.Sprintf("grpc.user.%d@example.com", i),
			Country:   s.createCountry,
			Password:  fmt.Sprintf("securepassword.%d", i),
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

	// Verify that we can list all created users.
	list, err := s.client.List(context.Background(), 1, tcase.numUsers+1, s.createCountry)
	require.NoError(t, err)
	assert.Equal(t, tcase.numUsers, list.TotalCount)
	assert.Equal(t, -1, list.NextPage)
	assert.Equal(t, tcase.numUsers, len(list.Users))
}

func (s *testSuiteGRPC) getUpdateDeleteUser(t *testing.T, tcase testCaseGRPC) {
	list, err := s.client.List(context.Background(), 1, 1, s.createCountry)
	require.NoError(t, err)
	assert.Equal(t, tcase.numUsers, list.TotalCount)
	assert.Equal(t, 2, list.NextPage)
	assert.Equal(t, 1, len(list.Users))

	oneUser := list.Users[0]

	user, err := s.client.Get(context.Background(), oneUser.ID)
	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, oneUser, *user)

	userUpdate := dto.UserUpdate{
		ID:        user.ID,
		FirstName: "gRPC",
		LastName:  "User",
		Nickname:  fmt.Sprintf("grpcuser.%s", user.ID.String()),
		Email:     fmt.Sprintf("grpc.user.%s@example.com", user.ID.String()),
		Country:   s.updateCountry,
		Password:  "newsecurepassword",
	}

	updatedUser, err := s.client.Update(context.Background(), userUpdate)
	require.NoError(t, err)
	assert.NotNil(t, updatedUser)

	assert.Equal(t, userUpdate.ID, updatedUser.ID)
	assert.Equal(t, userUpdate.FirstName, updatedUser.FirstName)
	assert.Equal(t, userUpdate.LastName, updatedUser.LastName)
	assert.Equal(t, userUpdate.Nickname, updatedUser.Nickname)
	assert.Equal(t, userUpdate.Email, updatedUser.Email)
	assert.Equal(t, userUpdate.Country, updatedUser.Country)
	assert.NotZero(t, updatedUser.CreatedAt)
	assert.Greater(t, updatedUser.UpdatedAt, oneUser.CreatedAt)

	err = s.client.Delete(context.Background(), user.ID)
	require.NoError(t, err)
}

func (s *testSuiteGRPC) listUsers(t *testing.T, tcase testCaseGRPC) {
	lastPage := 1 + (tcase.numUsers+tcase.pageSize/2)/tcase.pageSize
	for page := 1; page <= lastPage; page++ {
		list, err := s.client.List(context.Background(), page, tcase.pageSize, s.createCountry)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, tcase.numUsers, list.TotalCount)
		if page < lastPage {
			assert.Equal(t, page+1, list.NextPage)
			assert.Equal(t, tcase.pageSize, len(list.Users))
		} else {
			assert.Equal(t, -1, list.NextPage)
			assert.GreaterOrEqual(t, tcase.numUsers%tcase.pageSize, len(list.Users))
		}
	}
}

func (s *testSuiteGRPC) deleteAllUsers(t *testing.T, tcase testCaseGRPC) {
	list, err := s.client.List(context.Background(), 1, tcase.numUsers+1, s.createCountry)
	require.NoError(t, err)

	for i := 0; i < len(list.Users); i++ {
		err := s.client.Delete(context.Background(), list.Users[i].ID)
		require.NoError(t, err)
	}
}
