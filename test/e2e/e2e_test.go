package e2e_test

import (
	"context"
	"github.com/icon-project/centralized-relay/test/e2e/tests"
	"github.com/icon-project/centralized-relay/test/testsuite"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, new(E2ETest))
}

type E2ETest struct {
	testsuite.E2ETestSuite
}

func (s *E2ETest) TestE2E_all() {
	//go panicOnTimeout(10 * time.Hour) // custom timeout

	t := s.T()
	ctx := context.TODO()
	s.Require().NoError(s.SetCfg())
	_ = s.SetupChainsAndRelayer(ctx)
	xcall := tests.XCallTestSuite{
		E2ETestSuite: &s.E2ETestSuite,
		T:            t,
	}
	t.Run("test xcall", func(t *testing.T) {
		xcall.TextXCall()
	})

	t.Run("test xcall packet flush", func(t *testing.T) {
		//xcall.TestXCallFlush()
	})

}
