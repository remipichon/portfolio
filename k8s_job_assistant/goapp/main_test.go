// main_test.go
package main

import (
	"context"
	"flag"
	"github.com/stretchr/testify/suite"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"testing"
)

var (
	keepResources bool
)

func TestMain(m *testing.M) {
	// Custom flag to keep resources
	flag.BoolVar(&keepResources, "keep-resources", false, "Keep test resources after test run")
	flag.Parse()

	os.Exit(m.Run())
}

type KubeServiceIntegrationTestSuite struct {
	suite.Suite
	// KubeClient to check resources right in the test (should be used very lightly)
	ks *KubeService
	// create all namespaced resources in this one
	Namespace string
	// to easily cleanup resources, make sure to create all resources with these
	TestLabels  map[string]string
	BaseJobName string
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
		"testing-labels": "k8s-job-assitant",
	}
	s.Namespace = "kja-test-namespace"
	s.BaseJobName = "the-first-job"

	// init Kube client
	s.ks = &KubeService{}
	s.ks.InitClient()

	// Create Kube resources

	// A. Namespace
	_, err = s.ks.kubeClient.CoreV1().Namespaces().Create(ctx, &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: s.Namespace,
		}},
		metav1.CreateOptions{})
	s.Require().NoError(err, "failed to create namespace")

	// B. Correctly configured Jobs
	_, err = s.ks.kubeClient.BatchV1().Jobs(s.Namespace).Create(ctx, &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      s.BaseJobName,
			Namespace: s.Namespace,
			Annotations: map[string]string{
				"job-assistant": "enable",
			},
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
								"echo This is my awesome task which lasts 5 seconds!; sleep 5; This is the end of my awesome task",
							},
						},
					},
					RestartPolicy: corev1.RestartPolicyNever,
				},
			},
		},
	}, metav1.CreateOptions{})
	s.Require().NoError(err, "failed to create job")

}

func (s *KubeServiceIntegrationTestSuite) TearDownSuite() {
	if keepResources {
		s.T().Logf("Skipping namespace deletion because keep-resources is set.")
		return
	}
	ctx := context.Background()

	// B . Jobs
	err := s.ks.kubeClient.BatchV1().Jobs(s.Namespace).Delete(ctx, s.BaseJobName, metav1.DeleteOptions{})
	s.Require().NoError(err, "failed to delete job")

	// A. Namespace
	err = s.ks.kubeClient.CoreV1().Namespaces().Delete(ctx, s.Namespace, metav1.DeleteOptions{})
	s.Require().NoError(err, "failed to delete namespace")

}
