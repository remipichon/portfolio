// integration_test.go
package kube

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/suite"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"os"
	"testing"
	"time"
)

type KubeServiceIntegrationTestSuite struct {
	suite.Suite
	// KubeClient to check resources right in the test (should be used very lightly, prefer the app service)
	kubeClient *kubernetes.Clientset
	jobMgr     JobManager
	// create all namespaced resources in this one
	Namespace string
	// to easily cleanup resources, make sure to create all resources with these
	TestLabels map[string]string
	// to 'export' the JobManager annotation
	jobAssistAnnotation string
	// Kubernetes resources lifecycle
	tearDown      bool
	keepResources bool
}

func (s *KubeServiceIntegrationTestSuite) SetupSuite() {
	ctx := context.Background()
	var err error

	s.tearDown = TearDown
	s.keepResources = KeepResources

	// configure test suite
	s.TestLabels = map[string]string{
		"testing-labels": "under-test-k8s-job-assistant",
	}
	s.Namespace = "kja-test-namespace"
	s.jobAssistAnnotation = "under-test-for-job-assist"

	// init Kube client for tests assertions
	s.kubeClient = InitKubeClient()

	// init component under test
	s.jobMgr = NewJobManager(s.kubeClient, s.jobAssistAnnotation)

	if s.tearDown {
		fmt.Println("Running in tear-down mode because -tear-down")
		s.TearDownSuite()
		fmt.Println("Tear down has run, you can now re-run the test without -tear-down")
		os.Exit(0)
	}

	// namespace
	_, err = s.kubeClient.CoreV1().Namespaces().Create(ctx, &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: s.Namespace,
		}},
		metav1.CreateOptions{})
	s.Require().NoError(err, "failed to create namespace")

}

func (s *KubeServiceIntegrationTestSuite) TearDownTest() {
	if s.keepResources {
		s.T().Logf("Skipping jobs deletion because keep-resources is set.")
		return
	}

	//deleteJobAndWaitForDeletion all created jobs and wait for complete deletion to isolate test cases
	policy := metav1.DeletePropagationForeground
	err := s.kubeClient.BatchV1().Jobs(s.Namespace).DeleteCollection(
		context.Background(),
		metav1.DeleteOptions{
			PropagationPolicy: &policy,
		},
		metav1.ListOptions{
			LabelSelector: metav1.FormatLabelSelector(&metav1.LabelSelector{
				MatchLabels: s.TestLabels,
			}),
		},
	)
	if err != nil {
		s.T().Logf("Error deleting job: %v", err)
	}

	// the one in default namespace
	err = s.kubeClient.BatchV1().Jobs("default").DeleteCollection(
		context.Background(),
		metav1.DeleteOptions{
			PropagationPolicy: &policy,
		},
		metav1.ListOptions{
			LabelSelector: metav1.FormatLabelSelector(&metav1.LabelSelector{
				MatchLabels: s.TestLabels,
			}),
		},
	)
	if err != nil {
		s.T().Logf("Error deleting job: %v", err)
	}

}

func (s *KubeServiceIntegrationTestSuite) TearDownSuite() {
	if s.keepResources {
		s.T().Logf("Skipping namespace deletion because keep-resources is set.")
		return
	}
	ctx := context.Background()

	// B. Jobs are deleted after each test case

	// A. Namespace
	err := s.kubeClient.CoreV1().Namespaces().Delete(ctx, s.Namespace, metav1.DeleteOptions{})
	s.Require().NoError(err, "failed to deleteJobAndWaitForDeletion namespace")

	fmt.Println("Wait for namespace to be deleted (timeout 30s) run 'go test -tear-down' to keep trying if it times out: ", s.Namespace)
	end := time.Now().Add(time.Second * time.Duration(30)) //TODO export timeout
	for {
		_, err := s.kubeClient.CoreV1().Namespaces().Get(ctx, s.Namespace, metav1.GetOptions{})
		if errors.IsNotFound(err) {
			break
		}
		if time.Now().After(end) {
			s.Assert().FailNow("Timeout waiting for namespace to be deleted")
		}
		time.Sleep(200 * time.Millisecond) // polling interval
	}

}

func (s *KubeServiceIntegrationTestSuite) TestSetup() {

}

func (s *KubeServiceIntegrationTestSuite) TestListJob() {
	// Correctly configured Job
	job1, jobName := s.validJob("correct-job-list", s.TestLabels, 0)
	s.createJob(job1, false)
	// Correctly configured Job in default namespace
	_, err := s.kubeClient.BatchV1().Jobs("default").Create(context.Background(), job1, metav1.CreateOptions{})
	s.Require().NoError(err, "failed to create job")

	jobs, err := s.jobMgr.List()
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

func (s *KubeServiceIntegrationTestSuite) TestListJobEmpty() {
	jobs, err := s.jobMgr.List()
	s.Require().NoError(err)

	s.Assert().Len(jobs, 0)
}

//status: "suspended" | "running" | "failed" | "scheduled";

func (s *KubeServiceIntegrationTestSuite) TestListJobStatusSuspended() {
	// Correctly configured Job
	job1, jobName := s.validJob("correct-job-suspended", s.TestLabels, 0)
	s.createJob(job1, false)

	jobs, err := s.jobMgr.List()
	s.Require().NoError(err)

	s.Assert().Len(jobs, 1)
	s.Assert().Equal(jobs[0].Name, jobName)

	var ns []string
	for _, job := range jobs {
		ns = append(ns, job.Namespace)
	}

	s.Assert().Contains(ns, s.Namespace)
}

func (s *KubeServiceIntegrationTestSuite) TestRunJobNonExisting() {
	err := s.jobMgr.Run(s.Namespace, "non-existing")
	s.Require().Error(err)
	s.Assert().Contains(err.Error(), "jobs.batch")
	s.Assert().Contains(err.Error(), "not found")
}

/*func (s *KubeServiceIntegrationTestSuite) TestRunJobNonSchedulable() {
	validButUnschedulableJob := s.validJob(s.BaseJobName, s.TestLabels, 0)
	validButUnschedulableJob.Spec.Template.Spec.NodeSelector = map[string]string{
		"kubernetes.io/hostname": "nonexistent-node",
	}
	_, err = s.kubeClient.BatchV1().Jobs("default").Create(context.Background(), validButUnschedulableJob, metav1.CreateOptions{})
	s.Require().NoError(err, "failed to create job")

	err := s.jobMgr.Run("default", s.BaseJobName)
	s.Require().NoError(err)

	//this test only care that the Job scheduled at least one pod
	s.Require().Eventually(func() bool {
		err, jobStatus := s.jobMgr.Status(s.Namespace, s.BaseJobName)
		s.Require().NoError(err)

		if jobStatus.StartTime != nil {
			return true
		}

		return false
	}, 5*time.Second, 200*time.Millisecond, "Job did not started (scheduled one pod) within timeout")
}*/

// case where suspend=true
func (s *KubeServiceIntegrationTestSuite) TestRunJobAfterCreate() {
	job1, jobName := s.validJob("correct-job-run", s.TestLabels, 0)
	s.createJob(job1, true)

	err := s.jobMgr.Run(s.Namespace, jobName)
	s.Require().NoError(err)

	s.assertJobStarted(jobName)
}

// case where suspend is not present (taking over already created Job)
func (s *KubeServiceIntegrationTestSuite) TestRunJobWithoutSuspend() {
	jobName := "suspend-not-present"
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: jobName,
			Annotations: map[string]string{
				s.jobAssistAnnotation: "enable",
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
	s.waitForJobCompletion(s.Namespace, jobName, 20)

	err := s.jobMgr.Run(s.Namespace, jobName)
	s.Require().NoError(err)

	s.assertJobStarted(jobName)
}

// case where suspend=false (configured and run twice)
func (s *KubeServiceIntegrationTestSuite) TestRunJobAfterCompletion() {
	job, jobName := s.validJob("suspend-there-run-twice", s.TestLabels, 0)
	s.createJob(job, true)

	s.T().Logf("Run first time")
	err := s.jobMgr.Run(s.Namespace, jobName)
	s.Require().NoError(err)
	s.assertJobStarted(jobName)
	s.T().Logf("First run has started")

	s.waitForJobCompletion(s.Namespace, jobName, 600)
	s.T().Logf("First run has completed")

	s.T().Logf("Run second time")
	err = s.jobMgr.Run(s.Namespace, jobName)
	s.Require().NoError(err)
	s.assertJobStarted(jobName)
	s.T().Logf("Second run has started, test is over")

}

func (s *KubeServiceIntegrationTestSuite) TestRunJobWhileRunning() {
	job, jobName := s.validJob("suspend-there-run-twice-without-waiting-for-completion", s.TestLabels, 15)
	s.createJob(job, true)

	s.T().Logf("Run first time")
	err := s.jobMgr.Run(s.Namespace, jobName)
	s.Require().NoError(err)
	s.assertJobStarted(jobName)
	s.T().Logf("First run has started")

	s.T().Logf("Run second time (without waiting for first completion")
	err = s.jobMgr.Run(s.Namespace, jobName)
	s.Require().Error(err, &JobAlreadyRunningError{})
}

func (s *KubeServiceIntegrationTestSuite) TestKillJob() {
	job, jobName := s.validJob("suspend-there-run-to-kill", s.TestLabels, 15)
	s.createJob(job, true)

	err := s.jobMgr.Run(s.Namespace, jobName)
	s.Require().NoError(err)
	s.assertJobStarted(jobName)
	s.T().Logf("Run has started")

	err = s.jobMgr.Kill(s.Namespace, jobName)

	job, err = s.kubeClient.BatchV1().Jobs(s.Namespace).Get(context.Background(), jobName, metav1.GetOptions{})
	s.Require().NoError(err)

	s.Assert().Equal(int32(0), job.Status.Active)

}

func (s *KubeServiceIntegrationTestSuite) TestRunAfterKillJob() {
	job, jobName := s.validJob("suspend-there-run-after-kill", s.TestLabels, 15)
	s.createJob(job, true)

	err := s.jobMgr.Run(s.Namespace, jobName)
	s.Require().NoError(err)
	s.assertJobStarted(jobName)
	s.T().Logf("Run has started")

	err = s.jobMgr.Kill(s.Namespace, jobName)

	err = s.jobMgr.Run(s.Namespace, jobName)
	s.Require().NoError(err)
	s.assertJobStarted(jobName)
	s.T().Logf("Run has started")

}

func (s *KubeServiceIntegrationTestSuite) TestKillJobNonExisting() {
	err := s.jobMgr.Kill(s.Namespace, "non-existing")
	s.Require().Error(err)
	s.Assert().Contains(err.Error(), "jobs.batch")
	s.Assert().Contains(err.Error(), "not found")
}

func TestKubeService(t *testing.T) {
	suite.Run(t, new(KubeServiceIntegrationTestSuite))
}

// helpers functions

func (s *KubeServiceIntegrationTestSuite) assertJobStarted(jobName string) {
	// this test only care that the Job scheduled at least one pod
	s.Require().Eventually(func() bool {
		err, jobStatus := s.jobMgr.Status(s.Namespace, jobName)
		s.Require().NoError(err)

		if jobStatus.StartTime != nil {
			return true
		}

		return false
	}, 5*time.Second, 200*time.Millisecond, "Job did not started (scheduled one pod) within timeout")
}

func (s *KubeServiceIntegrationTestSuite) waitForJobCompletion(namespace string, jobName string, waitSecond int) {
	timeout := time.After(time.Duration(waitSecond) * time.Second)
	tick := time.Tick(200 * time.Millisecond) // poll every 200ms

	for {
		select {
		case <-timeout:
			s.FailNow("timed out waiting for Job to complete")
		case <-tick:
			err, jobStatus := s.jobMgr.Status(namespace, jobName)
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

func (s *KubeServiceIntegrationTestSuite) waitForJobStable(namespace, name string) {
	// Thanks for ChatGPT for this
	// WaitForJobStable polls the Job until its resourceVersion is stable.
	// It returns the fresh, stable Job object.
	ctx := context.Background()
	end := time.Now().Add(time.Second * 1)

	var lastRV string
	var job *batchv1.Job
	var err error

	for {
		job, err = s.kubeClient.BatchV1().Jobs(namespace).Get(ctx, name, metav1.GetOptions{})
		s.Require().NoError(err, "failed to fetch job")

		if job.ResourceVersion == lastRV {
			// No change detected → stable
			return
		}

		// Still mutating → wait and check again
		lastRV = job.ResourceVersion

		if time.Now().After(end) {
			s.Require().FailNow("Timeout waiting for Job to stabilize")
		}
		time.Sleep(100 * time.Millisecond) // polling interval
	}
}

func (s *KubeServiceIntegrationTestSuite) validJob(jobName string, testLabels map[string]string, sleepSecond float32) (*batchv1.Job, string) {
	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: jobName,
			Annotations: map[string]string{
				s.jobAssistAnnotation: "enable",
			},
			Labels: testLabels,
		},
		Spec: batchv1.JobSpec{
			Suspend: newTrue(),
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "some-awesomely-tested-job",
							Image: "busybox",
							Command: []string{
								"sh",
								"-c",
								fmt.Sprintf("echo This is my awesome task which lasts 5 seconds!; sleep %.3f; echo This is the end of my awesome task", sleepSecond),
							},
						},
					},
					RestartPolicy: corev1.RestartPolicyNever,
				},
			},
		},
	}, jobName
}

func (s *KubeServiceIntegrationTestSuite) createJob(job *batchv1.Job, waitForStable bool) {
	_, err := s.kubeClient.BatchV1().Jobs(s.Namespace).Create(context.Background(), job, metav1.CreateOptions{})
	s.Require().NoError(err, "failed to create job")
	if waitForStable {
		s.waitForJobStable(s.Namespace, job.Name)
	}
}
