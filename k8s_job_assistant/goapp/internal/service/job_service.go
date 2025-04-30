package service

import (
	"fmt"
	"goapp/internal/kube"
	"goapp/internal/model"
)

type JobService interface {
	ListDecoratedJobs() ([]model.DecoratedJob, error)
	Run(namespace, jobName string) error
	Kill(namespace, jobName string) error
}

type jobService struct {
	jobManager kube.JobManager
}

func NewJobService(j kube.JobManager) JobService {
	return &jobService{jobManager: j}
}

func (s *jobService) ListDecoratedJobs() ([]model.DecoratedJob, error) {
	jobs, err := s.jobManager.List()
	if err != nil {
		return nil, err
	}

	// Transform into decorated format
	result := make([]model.DecoratedJob, 0, len(jobs))
	for _, job := range jobs {
		decoratedJob := model.DecoratedJob{
			Namespace: job.Namespace,
			Name:      job.Name,
		}

		if job.Status.Active > 0 {
			decoratedJob.LastStatus = model.LastStatus{
				Type:    "Running",
				Message: fmt.Sprintf("%d pod(s)", job.Status.Active),
			}
		} else {
			if len(job.Status.Conditions) > 0 {
				latest := &job.Status.Conditions[0]
				for i := range job.Status.Conditions {
					if job.Status.Conditions[i].LastTransitionTime.After(latest.LastTransitionTime.Time) {
						latest = &job.Status.Conditions[i]
					}
				}
				decoratedJob.LastStatus = model.LastStatus{
					Type:    string(latest.Type),
					Message: latest.Message,
				}
			}
		}

		decoratedJob.LastSuccessfullyRunStarTime = job.Status.StartTime
		decoratedJob.LastSuccessfullyRunCompletionTime = job.Status.CompletionTime

		result = append(result, decoratedJob)
	}
	return result, nil
}

func (s *jobService) Run(namespace, jobName string) error {
	return s.jobManager.Run(namespace, jobName)
}

func (s *jobService) Kill(namespace, jobName string) error {
	return s.jobManager.Kill(namespace, jobName)
}
