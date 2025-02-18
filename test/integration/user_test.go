//go:build integration_tests

package test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// Test template

type testCaseEdgar struct{}

var testCases = map[string]testCaseEdgar{
	"test_case_one": {},
	"test_case_two": {},
}

type TestSuite struct {
	suite.Suite
	// other test suite vars here
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) SetupSuite() {
	fmt.Println("SETUP SUITE")
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
	for tcn := range testCases {
		s.T().Run(tcn, func(t *testing.T) {
			t.Parallel()
			assert.True(t, true)
			require.True(t, true)
		})
	}
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
