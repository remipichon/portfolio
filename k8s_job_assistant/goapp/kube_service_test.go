// integration_test.go
package main

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/suite"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
	"time"
)

func TestKubeService(t *testing.T) {
	suite.Run(t, new(KubeServiceIntegrationTestSuite))
}

func (s *KubeServiceIntegrationTestSuite) TestSetup() {

}

// SetupSuite also creates a non-valid Job
func (s *KubeServiceIntegrationTestSuite) TestListJob() {
	// Correctly configured Job
	job1, jobName := validJob("correct-job-list", s.TestLabels, 0)
	s.createJob(job1, false)
	// Correctly configured Job in default namespace
	_, err := s.ks.kubeClient.BatchV1().Jobs("default").Create(context.Background(), job1, metav1.CreateOptions{})
	s.Require().NoError(err, "failed to create job")

	jobs, err := s.ks.list()
	s.Require().NoError(err)

	s.Assert().Len(jobs, 2)
	s.Assert().Equal(jobs[0].Name, jobName)
	s.Assert().Equal(jobs[1].Name, jobName)

	var ns []string
	for _, job := range jobs {
		ns = append(ns, job.Namespace)
	}

	s.Assert().Contains(ns, s.Namespace)
	s.Assert().Contains(ns, "default")
}

func (s *KubeServiceIntegrationTestSuite) TestRunJobNonExisting() {
	err := s.ks.run(s.Namespace, "non-existing")
	s.Require().Error(err)
	s.Assert().Contains(err.Error(), "jobs.batch")
	s.Assert().Contains(err.Error(), "not found")
}

/*func (s *KubeServiceIntegrationTestSuite) TestRunJobNonSchedulable() {
	validButUnschedulableJob := validJob(s.BaseJobName, s.TestLabels, 0)
	validButUnschedulableJob.Spec.Template.Spec.NodeSelector = map[string]string{
		"kubernetes.io/hostname": "nonexistent-node",
	}
	_, err = s.ks.kubeClient.BatchV1().Jobs("default").Create(context.Background(), validButUnschedulableJob, metav1.CreateOptions{})
	s.Require().NoError(err, "failed to create job")

	err := s.ks.run("default", s.BaseJobName)
	s.Require().NoError(err)

	//this test only care that the Job scheduled at least one pod
	s.Require().Eventually(func() bool {
		err, jobStatus := s.ks.status(s.Namespace, s.BaseJobName)
		s.Require().NoError(err)

		if jobStatus.StartTime != nil {
			return true
		}

		return false
	}, 5*time.Second, 200*time.Millisecond, "Job did not started (scheduled one pod) within timeout")
}*/

// case where suspend=true
func (s *KubeServiceIntegrationTestSuite) TestRunJobAfterCreate() {
	job1, jobName := validJob("correct-job-run", s.TestLabels, 0)
	s.createJob(job1, true)

	err := s.ks.run(s.Namespace, jobName)
	s.Require().NoError(err)

	assertJobStarted(s, jobName)
}

// this test only care that the Job scheduled at least one pod
func assertJobStarted(s *KubeServiceIntegrationTestSuite, jobName string) {
	s.Require().Eventually(func() bool {
		err, jobStatus := s.ks.status(s.Namespace, jobName)
		s.Require().NoError(err)

		if jobStatus.StartTime != nil {
			return true
		}

		return false
	}, 5*time.Second, 200*time.Millisecond, "Job did not started (scheduled one pod) within timeout")
}

// case where suspend is not present (taking over already created Job)
func (s *KubeServiceIntegrationTestSuite) TestRunJobWithoutSuspend() {
	jobName := "suspend-not-present"
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: jobName,
			Annotations: map[string]string{
				"job-assistant": "enable",
			},
			Labels: s.TestLabels,
		},
		Spec: batchv1.JobSpec{
			// no Suspend at all
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "some-awesomely-tested-job",
							Image: "busybox",
							Command: []string{
								"sh",
								"-c",
								fmt.Sprintf("echo This is my awesome task which lasts 5 seconds!; sleep %.3f; echo This is the end of my awesome task", 0.0),
							},
						},
					},
					RestartPolicy: corev1.RestartPolicyNever,
				},
			},
		},
	}
	s.createJob(job, true)

	//because the Job is created without suspend, we need to wait for it to complete
	//before running the actual test
	s.WaitForJobCompletion(s.Namespace, jobName, 20)

	err := s.ks.run(s.Namespace, jobName)
	s.Require().NoError(err)

	assertJobStarted(s, jobName)
}

// case where suspend=false (configured and run twice)
func (s *KubeServiceIntegrationTestSuite) TestRunASecondtimeJobAfterCompletion() {
	job, jobName := validJob("suspend-there-run-twice", s.TestLabels, 0)
	s.createJob(job, true)

	s.T().Logf("Run first time")
	err := s.ks.run(s.Namespace, jobName)
	s.Require().NoError(err)
	assertJobStarted(s, jobName)
	s.T().Logf("First run has started")

	s.WaitForJobCompletion(s.Namespace, jobName, 600)
	s.T().Logf("First run has completed")

	s.T().Logf("Run second time")
	err = s.ks.run(s.Namespace, jobName)
	s.Require().NoError(err)
	assertJobStarted(s, jobName)
	s.T().Logf("Second run has started, test is over")

}

func (s *KubeServiceIntegrationTestSuite) TestRunJobWhileRunning() {
	job, jobName := validJob("suspend-there-run-twice-without-waiting-for-completion", s.TestLabels, 15)
	s.createJob(job, true)

	s.T().Logf("Run first time")
	err := s.ks.run(s.Namespace, jobName)
	s.Require().NoError(err)
	assertJobStarted(s, jobName)
	s.T().Logf("First run has started")

	s.T().Logf("Run second time (without waiting for first completion")
	err = s.ks.run(s.Namespace, jobName)
	s.Require().Error(err, &JobAlreadyRunningError{})
}

func (s *KubeServiceIntegrationTestSuite) WaitForJobCompletion(namespace string, jobName string, waitSecond int) {
	timeout := time.After(time.Duration(waitSecond) * time.Second)
	tick := time.Tick(200 * time.Millisecond) // poll every 200ms

	for {
		select {
		case <-timeout:
			s.FailNow("timed out waiting for Job to complete")
		case <-tick:
			err, jobStatus := s.ks.status(namespace, jobName)
			s.Require().NoError(err)

			for _, condition := range jobStatus.Conditions {
				if condition.Type == batchv1.JobComplete && condition.Status == corev1.ConditionTrue {
					s.Assert().Equal(corev1.ConditionTrue, condition.Status)
					return // done, exit test
				}
			}
		}
	}
}

func (s *KubeServiceIntegrationTestSuite) TestKillJob() {
	job, jobName := validJob("suspend-there-run-to-kill", s.TestLabels, 15)
	s.createJob(job, true)

	err := s.ks.run(s.Namespace, jobName)
	s.Require().NoError(err)
	assertJobStarted(s, jobName)
	s.T().Logf("Run has started")

	err = s.ks.kill(s.Namespace, jobName)

	job, err = s.ks.kubeClient.BatchV1().Jobs(s.Namespace).Get(context.Background(), jobName, metav1.GetOptions{})
	s.Require().NoError(err)

	s.Assert().Equal(int32(0), job.Status.Active)

}

/*

	s.Require().Eventually(func() bool {
		err, jobStatus := s.ks.status(s.Namespace, s.BaseJobName)
		s.Require().NoError(err)

		for _, condition := range jobStatus.Conditions {
			if condition.Type == batchv1.JobComplete && condition.Status == corev1.ConditionTrue {
				return true
			}
		}
		return false
	}, 5*time.Second, 200*time.Millisecond, "Job did not complete within timeout")
*/
