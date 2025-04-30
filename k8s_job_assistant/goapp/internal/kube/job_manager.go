// Package kube implements service methods to manage Kubernetes Jobs
package kube

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/utils/pointer"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// JobManager provides helper methods to interact with Kubernetes Jobs.
type jobManager struct {
	kubeClient          *kubernetes.Clientset
	jobAssistAnnotation string
}

type JobManager interface {
	List() ([]batchv1.Job, error)
	Run(namespace, jobName string) error
	Kill(namespace, jobName string) error
	Status(namespace, jobName string) (error, *batchv1.JobStatus)
}

func NewJobManager(kubeClient *kubernetes.Clientset, jobAssistAnnotation string) JobManager {
	return &jobManager{kubeClient: kubeClient, jobAssistAnnotation: jobAssistAnnotation}
}

// List lists Jobs with annotation 'job-assistant' set to true on any namespace.
func (j *jobManager) List() ([]batchv1.Job, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	jobs, err := j.kubeClient.BatchV1().Jobs("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var filtered []batchv1.Job
	for _, job := range jobs.Items {
		if val, ok := job.Annotations[j.jobAssistAnnotation]; ok && val == "enable" {
			filtered = append(filtered, job)
		}
	}
	return filtered, nil
}

// Run runs a Job, fails if already running, handle Suspend:true and clean re-create when needed.
func (j *jobManager) Run(namespace, jobName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	job, err := j.kubeClient.BatchV1().Jobs(namespace).Get(ctx, jobName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	if isJobRunning(job.Status) {
		return &JobAlreadyRunningError{}
	}

	// suspended=true, set it to false for Kube to run the Job right away
	if job.Spec.Suspend != nil && *job.Spec.Suspend {
		job.Spec.Suspend = newFalse()
		_, err = j.kubeClient.BatchV1().Jobs(namespace).Update(ctx, job, metav1.UpdateOptions{})
		return err
	}

	// suspended=false or absent, deleteJobAndWaitForDeletion Job then recreate it
	err = deleteJobAndWaitForDeletion(j.kubeClient, namespace, job.Name)
	if err != nil {
		return err
	}

	_, err = j.kubeClient.BatchV1().Jobs(namespace).Create(ctx, cleanJobForRecreate(job), metav1.CreateOptions{})
	return err
}

// Status returns the full Kubernetes status of job, without any decoration.
func (j *jobManager) Status(namespace, jobName string) (error, *batchv1.JobStatus) {
	job, err := j.kubeClient.BatchV1().Jobs(namespace).Get(context.Background(), jobName, metav1.GetOptions{})
	if err != nil {
		return err, nil
	}

	return nil, &job.Status
}

// Kill suspends the Job and delete all of its running pod.
func (j *jobManager) Kill(namespace, jobName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second) //TODO configure deletion timeout
	defer cancel()
	//Job is kept for later usage

	// suspend the Job to prevent Kubernetes from recreating the pods
	patch := []byte(`{"spec":{"suspend":true}}`)
	_, err := j.kubeClient.BatchV1().Jobs(namespace).Patch(ctx, jobName, types.MergePatchType, patch, metav1.PatchOptions{})
	if err != nil {
		return err
	}

	// deleteJobAndWaitForDeletion the pods by labels
	err = j.kubeClient.CoreV1().Pods(namespace).DeleteCollection(
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
		pods, getPodErr := j.kubeClient.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{
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

func deleteJobAndWaitForDeletion(kubeClient *kubernetes.Clientset, namespace, jobName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second) //TODO configure deletion timeout
	defer cancel()

	policy := metav1.DeletePropagationForeground
	err := kubeClient.BatchV1().Jobs(namespace).Delete(ctx, jobName, metav1.DeleteOptions{
		PropagationPolicy: &policy,
	})
	if err != nil {
		return err
	}

	// wait for actual full deletion
	for {
		_, err := kubeClient.BatchV1().Jobs(namespace).Get(ctx, jobName, metav1.GetOptions{})
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
		if (cond.Type == batchv1.JobComplete || cond.Type == batchv1.JobFailed || cond.Type == batchv1.JobSuspended) && cond.Status == corev1.ConditionTrue {
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
	jobNameLabel := "job-name" // how Kubernetes link Jobs to Pods, it's standard

	if _, ok := job.Spec.Selector.MatchLabels[jobNameLabel]; !ok {
		job.Spec.Selector.MatchLabels[jobNameLabel] = job.Name
	}
	job.Spec.Template.Labels[jobNameLabel] = job.Name
}
