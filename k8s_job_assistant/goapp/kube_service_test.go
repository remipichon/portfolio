// integration_test.go
package main

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestKubeServiceIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(KubeServiceIntegrationTestSuite))
}

func (s *KubeServiceIntegrationTestSuite) TestListJob() {
	jobs, err := s.ks.list()
	s.Require().NoError(err)

	s.Assert().Len(jobs, 1)
	s.Assert().Equal(jobs[0].Name, s.BaseJobName)
}
