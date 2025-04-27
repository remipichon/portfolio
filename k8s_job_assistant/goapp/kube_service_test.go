// integration_test.go
package main

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestKubeServiceIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(KubeServiceIntegrationTestSuite))
}

// SetupSuite also create a non valid Job
func (s *KubeServiceIntegrationTestSuite) TestListJob() {
	jobs, err := s.ks.list()
	s.Require().NoError(err)

	s.Assert().Len(jobs, 2)
	s.Assert().Equal(jobs[0].Name, s.BaseJobName)
	s.Assert().Equal(jobs[1].Name, s.BaseJobName)

	var ns []string
	for _, job := range jobs {
		ns = append(ns, job.Namespace)
	}

	s.Assert().Contains(ns, s.Namespace)
	s.Assert().Contains(ns, "default")
}
