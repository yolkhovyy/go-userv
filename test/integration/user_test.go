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

type testCaseEdgar struct {
	numUsers int
}

var testCases = map[string]testCaseEdgar{
	"test_case_one": {numUsers: 100},
	"test_case_two": {},
}

type TestSuite struct {
	suite.Suite
	client *userrest.Client
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) SetupSuite() {
	fmt.Println("SETUP SUITE")
	s.client = userrest.NewClient("http://localhost:8080") // Adjust the URL as needed
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
	for tcn, tc := range testCases {
		s.T().Run(tcn, func(t *testing.T) {
			t.Parallel()
			if tcn == "test_case_one" {
				s.createUsers(t, tc.numUsers)
			} else {
				assert.True(t, true)
				require.True(t, true)
			}
		})
	}
}

func (s *TestSuite) createUsers(t *testing.T, numUsers int) {
	for i := 0; i < numUsers; i++ {
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
	}

	// Verify that we can list all created users
	// users, err := s.client.List(context.Background(), 1, numUsers, "")
	// require.NoError(t, err)
	// assert.Equal(t, numUsers, len(users.Users))
}

func (s *TestSuite) TestTwo() {
	s.T().Parallel()
	for tcn := range testCases {
		s.T().Run(tcn, func(t *testing.T) {
			t.Parallel()
			assert.True(t, true)
			require.True(t, true)
		})
	}
}
