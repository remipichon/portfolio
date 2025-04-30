package service

import (
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
	for _, j := range jobs {
		result = append(result, model.DecoratedJob{
			Namespace:   j.Namespace,
			Name:        j.Name,
			Completions: j.Spec.Completions,
			Labels:      j.Labels,
		})
	}
	return result, nil
}

func (s *jobService) Run(namespace, jobName string) error {
	return s.jobManager.Run(namespace, jobName)
}

func (s *jobService) Kill(namespace, jobName string) error {
	return s.jobManager.Kill(namespace, jobName)
}
