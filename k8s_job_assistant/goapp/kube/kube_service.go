// Package kube implements service methods to manage Kubernetes Jobs
package kube

import (
	"context"
	"flag"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/utils/pointer"
	"path/filepath"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Service provides helper methods to interact with Kubernetes Jobs.
type Service struct {
	kubeClient          *kubernetes.Clientset
	jobAssistAnnotation string
}

// InitClient instantiate a Kubernetes client based on local kubeconfig.
func (ks *Service) InitClient() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	ks.kubeClient = kubeClient

}

// SetJobAssistAnnotation configures the annotation used to take ownerships of Jobs
//
// Only Jobs with this annotation will be considered
func (ks *Service) SetJobAssistAnnotation(jobAssistAnnotation string) {
	fmt.Println(fmt.Sprintf("Setting JobAssist annotation to %s (only job with annotation %s=enabled will be considered)", jobAssistAnnotation, jobAssistAnnotation))
	ks.jobAssistAnnotation = jobAssistAnnotation
}

// List lists Jobs with annotation 'job-assistant' set to true on any namespace.
func (ks *Service) List() ([]batchv1.Job, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	jobs, err := ks.kubeClient.BatchV1().Jobs("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var filtered []batchv1.Job
	for _, job := range jobs.Items {
		if val, ok := job.Annotations[ks.jobAssistAnnotation]; ok && val == "enable" {
			filtered = append(filtered, job)
		}
	}
	return filtered, nil
}

// Run runs a Job, fails if already running, handle Suspend:true and clean re-create when needed.
func (ks *Service) Run(namespace, jobName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	job, err := ks.kubeClient.BatchV1().Jobs(namespace).Get(ctx, jobName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	if isJobRunning(job.Status) {
		return &JobAlreadyRunningError{}
	}

	// suspended=true, set it to false for Kube to run the Job right away
	if job.Spec.Suspend != nil && *job.Spec.Suspend {
		job.Spec.Suspend = newFalse()
		_, err = ks.kubeClient.BatchV1().Jobs(namespace).Update(ctx, job, metav1.UpdateOptions{})
		return err
	}

	// suspended=false or absent, deleteJobAndWaitForDeletion Job then recreate it
	err = ks.deleteJobAndWaitForDeletion(namespace, job.Name)
	if err != nil {
		return err
	}

	_, err = ks.kubeClient.BatchV1().Jobs(namespace).Create(ctx, cleanJobForRecreate(job), metav1.CreateOptions{})
	return err
}

// Status returns the full Kubernetes status of job, without any decoration.
func (ks *Service) Status(namespace, jobName string) (error, *batchv1.JobStatus) {
	job, err := ks.kubeClient.BatchV1().Jobs(namespace).Get(context.Background(), jobName, metav1.GetOptions{})
	if err != nil {
		return err, nil
	}

	return nil, &job.Status
}

// Kill suspends the Job and delete all of its running pod.
//
// TODO kill pods (to keep them and their logs) instead of deleting pods
func (ks *Service) Kill(namespace, jobName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second) //TODO configure deletion timeout
	defer cancel()
	//Job is kept for later usage

	// suspend the Job to prevent Kubernetes from recreating the pods
	patch := []byte(`{"spec":{"suspend":true}}`)
	_, err := ks.kubeClient.BatchV1().Jobs(namespace).Patch(ctx, jobName, types.MergePatchType, patch, metav1.PatchOptions{})
	if err != nil {
		return err
	}

	// deleteJobAndWaitForDeletion the pods by labels
	err = ks.kubeClient.CoreV1().Pods(namespace).DeleteCollection(
		ctx,
		metav1.DeleteOptions{},
		metav1.ListOptions{
			LabelSelector: fmt.Sprintf("job-name=%s", jobName),
		},
	)
	if err != nil {
		return err
	}

	// wait for actual pods deletion
	for {
		pods, getPodErr := ks.kubeClient.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{
			LabelSelector: fmt.Sprintf("job-name=%s", jobName),
		})
		if getPodErr != nil {
			return getPodErr
		}
		if len(pods.Items) == 0 {
			break
		}
		if ctx.Err() != nil {
			return fmt.Errorf("timed out waiting for Job's pods deletion: %v", ctx.Err())
		}
		time.Sleep(200 * time.Millisecond) // polling interval
	}

	return nil
}

func (ks *Service) deleteJobAndWaitForDeletion(namespace, jobName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second) //TODO configure deletion timeout
	defer cancel()

	policy := metav1.DeletePropagationForeground
	err := ks.kubeClient.BatchV1().Jobs(namespace).Delete(ctx, jobName, metav1.DeleteOptions{
		PropagationPolicy: &policy,
	})
	if err != nil {
		return err
	}

	// wait for actual full deletion
	for {
		_, err := ks.kubeClient.BatchV1().Jobs(namespace).Get(ctx, jobName, metav1.GetOptions{})
		if errors.IsNotFound(err) {
			break
		}
		if ctx.Err() != nil {
			return fmt.Errorf("timed out waiting for job deletion: %v", ctx.Err())
		}
		time.Sleep(200 * time.Millisecond) // polling interval
	}

	return nil
}

func isJobRunning(status batchv1.JobStatus) bool {
	if status.StartTime == nil {
		return false // Not pod scheduled, not even started yet
	}

	for _, cond := range status.Conditions {
		if (cond.Type == batchv1.JobComplete || cond.Type == batchv1.JobFailed) && cond.Status == corev1.ConditionTrue {
			return false // Already completed or failed
		}
	}

	return true // Started and not complete/failed yet
}

// cleanJobForRecreate returns a clean copy of the given Job, ready for recreation.
func cleanJobForRecreate(original *batchv1.Job) *batchv1.Job {
	job := original.DeepCopy()

	// Clean metadata: remove fields that must not be reused
	job.ResourceVersion = ""
	job.UID = ""
	job.Generation = 0
	job.CreationTimestamp = metav1.Time{}
	job.ManagedFields = nil
	job.SelfLink = ""
	job.Finalizers = nil      // optional: clear if you don't want old finalizers
	job.OwnerReferences = nil // optional: clear if you don't want old owners

	// Status must be empty
	job.Status = batchv1.JobStatus{}

	// Fix selector + pod template labels
	ensureJobSelectorMatchesTemplate(job)

	return job
}

// ensureJobSelectorMatchesTemplate deals with spec.selector :
// - keep the ones set by the user
// - mimic Kubernetes selector 'job-name'
func ensureJobSelectorMatchesTemplate(job *batchv1.Job) {
	if job.Spec.Selector == nil {
		job.Spec.Selector = &metav1.LabelSelector{
			MatchLabels: map[string]string{},
		}
	}
	if job.Spec.Template.Labels == nil {
		job.Spec.Template.Labels = map[string]string{}
	}

	// since we are now setting the selector ourselves
	job.Spec.ManualSelector = pointer.Bool(true)

	// we have to mimic what Kubernetes is doing automatically
	jobNameLabel := "job-name" // how Kubernetes links Jobs to Pods, it's standard

	if _, ok := job.Spec.Selector.MatchLabels[jobNameLabel]; !ok {
		job.Spec.Selector.MatchLabels[jobNameLabel] = job.Name
	}
	job.Spec.Template.Labels[jobNameLabel] = job.Name
}
