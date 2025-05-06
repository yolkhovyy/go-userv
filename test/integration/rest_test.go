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
	"github.com/yolkhovyy/go-otelw/otelw/otlp"
	"github.com/yolkhovyy/go-otelw/otelw/tracew"
	restclient "github.com/yolkhovyy/go-userv/client/user-rest"
	"github.com/yolkhovyy/go-userv/contract/dto"
)

type testCaseRest struct {
	name     string
	testFunc func(t *testing.T, tcase testCaseRest)
	numUsers int
	pageSize int
}

type testSuiteRest struct {
	suite.Suite
	client        *restclient.Client
	cases         []testCaseRest
	createCountry string
	updateCountry string
	tracer        *tracew.Tracer
}

func TestRestSuiteRun(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(testSuiteRest))
}

func (s *testSuiteRest) SetupSuite() {
	var err error

	if s.tracer, err = tracew.Configure(context.Background(), tracew.Config{
		Enable: true,
		OTLP: otlp.Config{
			Protocol: otlp.GRPC,
			Endpoint: "localhost:4317",
			Insecure: true,
		},
	}, nil); err != nil {
		panic(err)
	}

	s.client = restclient.NewClient("http://localhost:8080")
	s.cases = []testCaseRest{
		{name: "CreateUsers", testFunc: s.createUsers, numUsers: 100},
		{name: "GetUpdateDeleteUser", testFunc: s.getUpdateDeleteUser, numUsers: 100},
		{name: "ListUsers", testFunc: s.listUsers, numUsers: 100, pageSize: 12},
		{name: "DeleteAll", testFunc: s.deleteAllUsers, numUsers: 100},
	}

	s.createCountry = "XA"
	s.updateCountry = "XB"
}

func (s *testSuiteRest) TearDownSuite() {
	_ = s.tracer.Shutdown(context.Background())
}

func (s *testSuiteRest) SetupTest() {
	// TODO: setup test.
}

func (s *testSuiteRest) TearDownTest() {
	// TODO: tear down test.
}

func (s *testSuiteRest) TestRest() {
	for _, tc := range s.cases {
		s.T().Run(tc.name, func(t *testing.T) {
			tc.testFunc(t, tc)
		})
	}
}

func (s *testSuiteRest) createUsers(t *testing.T, tcase testCaseRest) {
	for i := 0; i < tcase.numUsers; i++ {
		userID := uuid.New()
		userInput := dto.UserInput{
			FirstName: fmt.Sprintf("Rest"),
			LastName:  fmt.Sprintf("User, %d", i),
			Nickname:  fmt.Sprintf("restuser%s", userID.String()),
			Email:     fmt.Sprintf("rest.user.%s@example.com", userID.String()),
			Country:   s.createCountry,
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
	list, err := s.client.List(context.Background(), 1, tcase.numUsers+1, s.createCountry)
	require.NoError(t, err)
	assert.Equal(t, tcase.numUsers, list.TotalCount)
	assert.Equal(t, -1, list.NextPage)
	assert.Equal(t, tcase.numUsers, len(list.Users))
}

func (s *testSuiteRest) getUpdateDeleteUser(t *testing.T, tcase testCaseRest) {
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
		FirstName: "Rest",
		LastName:  "User",
		Nickname:  fmt.Sprintf("restuser.%s", user.ID.String()),
		Email:     fmt.Sprintf("rest.user.%s@example.com", user.ID.String()),
		Country:   s.updateCountry,
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
	assert.Greater(t, updatedUser.UpdatedAt, oneUser.CreatedAt)

	err = s.client.Delete(context.Background(), user.ID)
	require.NoError(t, err)
}

func (s *testSuiteRest) listUsers(t *testing.T, tcase testCaseRest) {
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

func (s *testSuiteRest) deleteAllUsers(t *testing.T, tcase testCaseRest) {
	list, err := s.client.List(context.Background(), 1, tcase.numUsers+1, s.createCountry)
	require.NoError(t, err)

	for i := 0; i < len(list.Users); i++ {
		err := s.client.Delete(context.Background(), list.Users[i].ID)
		require.NoError(t, err)
	}
}
