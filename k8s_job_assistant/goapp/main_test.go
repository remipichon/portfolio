// main_test.go
package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/stretchr/testify/suite"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"testing"
	"time"
)

var (
	keepResources bool
	tearDown      bool
)

func TestMain(m *testing.M) {
	// Custom flag to keep resources
	flag.BoolVar(&keepResources, "keep-resources", false, "Keep test resources after test run")
	flag.BoolVar(&tearDown, "tear-down", false, "Attempt to delete resources before testing, in case of leftovers")
	flag.Parse()

	if tearDown {
		fmt.Println("-tear-down flag set to true")
	}
	if keepResources {
		fmt.Println("-keep-resources flag set to true")
	}

	os.Exit(m.Run())
}

type KubeServiceIntegrationTestSuite struct {
	suite.Suite
	// KubeClient to check resources right in the test (should be used very lightly)
	ks *KubeService
	// create all namespaced resources in this one
	Namespace string
	// to easily cleanup resources, make sure to create all resources with these
	TestLabels map[string]string
}

// Create resources common to all tests
//
//	init Kube client
//	prepare Namespace with resources
func (s *KubeServiceIntegrationTestSuite) SetupSuite() {
	ctx := context.Background()

	var err error

	// configure test suite
	s.TestLabels = map[string]string{
		"testing-labels": "k8s-job-assistant",
	}
	s.Namespace = "kja-test-namespace"

	// init Kube client
	s.ks = &KubeService{}
	s.ks.InitClient()

	if tearDown {
		fmt.Println("Running in tear-down mode because -tear-down")
		s.TearDownSuite()
		fmt.Println("Tear down has run, you can now re-run the test without -tear-down")
		os.Exit(0)
	}

	// Create Kube resources

	// A. Namespace
	_, err = s.ks.kubeClient.CoreV1().Namespaces().Create(ctx, &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: s.Namespace,
		}},
		metav1.CreateOptions{})
	s.Require().NoError(err, "failed to create namespace")

}

func (s *KubeServiceIntegrationTestSuite) TearDownTest() {
	if keepResources {
		s.T().Logf("Skipping jobs deletion because keep-resources is set.")
		return
	}

	//delete all created jobs and wait for complete deletion to isolate test cases
	// B . Jobs (based on testing labels)
	policy := metav1.DeletePropagationForeground
	err := s.ks.kubeClient.BatchV1().Jobs(s.Namespace).DeleteCollection(
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
		// don't s.T().Logf as we don't want to halt there
		fmt.Println("Error deleting job:", err)
	}

	// B the one in default namespace
	err = s.ks.kubeClient.BatchV1().Jobs("default").DeleteCollection(
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
		// don't s.T().Logf as we don't want to halt there
		fmt.Println("Error deleting job:", err)
	}

}

func (s *KubeServiceIntegrationTestSuite) createJob(job *batchv1.Job, waitForStable bool) {
	_, err := s.ks.kubeClient.BatchV1().Jobs(s.Namespace).Create(context.Background(), job, metav1.CreateOptions{})
	s.Require().NoError(err, "failed to create job")
	if waitForStable {
		s.waitForJobStable(s.Namespace, job.Name)
	}
}

// Thanks ChatGPT for this
// WaitForJobStable polls the Job until its resourceVersion is stable.
// It returns the fresh, stable Job object.
func (s *KubeServiceIntegrationTestSuite) waitForJobStable(namespace, name string) {
	ctx := context.Background()
	end := time.Now().Add(time.Second * 1)

	var lastRV string
	var job *batchv1.Job
	var err error

	for {
		job, err = s.ks.kubeClient.BatchV1().Jobs(namespace).Get(ctx, name, metav1.GetOptions{})
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

func validJob(jobName string, testLabels map[string]string, sleepSecond float32) (*batchv1.Job, string) {
	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: jobName,
			Annotations: map[string]string{
				"job-assistant": "enable",
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

func (s *KubeServiceIntegrationTestSuite) TearDownSuite() {
	if keepResources {
		s.T().Logf("Skipping namespace deletion because keep-resources is set.")
		return
	}
	ctx := context.Background()

	// B. Jobs are deleted after each test case

	// A. Namespace
	err := s.ks.kubeClient.CoreV1().Namespaces().Delete(ctx, s.Namespace, metav1.DeleteOptions{})
	s.Require().NoError(err, "failed to delete namespace")

	fmt.Println("Wait for namespace to be deleted (timeout 30s) run 'go test -tear-down' to keep trying if it times out: ", s.Namespace)
	end := time.Now().Add(time.Second * time.Duration(30)) //TODO export timeout
	for {
		_, err := s.ks.kubeClient.CoreV1().Namespaces().Get(ctx, s.Namespace, metav1.GetOptions{})
		if errors.IsNotFound(err) {
			break
		}
		if time.Now().After(end) {
			s.Assert().FailNow("Timeout waiting for namespace to be deleted")
		}
		time.Sleep(200 * time.Millisecond) // polling interval
	}

}
