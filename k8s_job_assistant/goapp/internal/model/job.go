package model

type DecoratedJob struct {
	Namespace   string            `json:"namespace"`
	Name        string            `json:"name"`
	Completions *int32            `json:"completions,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
}

type ListJobs struct {
	Jobs  []DecoratedJob `json:"jobs"`
	Count int            `json:"count"`
}
