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
	testFunc func(t *testing.T, tcase testCase)
	numUsers int
	pageSize int
	country  string
}

type TestSuite struct {
	suite.Suite
	client *userrest.Client
	cases  map[string]testCase
}

// testCases = map[string]testCase{
// 	"Create": {testFunc: ,numUsers: 100},
// }

func TestSuiteRun(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) SetupSuite() {
	fmt.Println("SETUP SUITE")
	s.client = userrest.NewClient("http://localhost:8080")
	s.cases = map[string]testCase{
		"Create": {testFunc: s.createUsers, numUsers: 100},
	}
}

func (s *TestSuite) TearDownSuite() {
	fmt.Println("TEAR DOWN SUITE")
}

func (s *TestSuite) SetupTest() {
	fmt.Println("SETUP TEST")
}

func (s *TestSuite) TearDownTest() {
	fmt.Println("TEAR DOWN TEST")
}

func (s *TestSuite) TestOne() {
	s.T().Parallel()
	for tcn, tc := range s.cases {
		s.T().Run(tcn, func(t *testing.T) {
			t.Parallel()
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

	// Verify that we can list all created listUsers.
	listUsers, err := s.client.List(context.Background(), 1, tcase.numUsers+1, tcase.country)
	require.NoError(t, err)
	assert.Equal(t, tcase.numUsers, listUsers.TotalCount)
	assert.Equal(t, -1, listUsers.NextPage)
	assert.Equal(t, tcase.numUsers, len(listUsers.Users))
}
