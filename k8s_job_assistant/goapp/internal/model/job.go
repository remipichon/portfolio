package model

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DecoratedJob struct {
	Namespace                         string       `json:"namespace"`
	Name                              string       `json:"name"`
	LastSuccessfullyRunStarTime       *metav1.Time `json:"lastSuccessfullyRunStarTime,omitempty"`
	LastStatus                        LastStatus   `json:"lastStatus"`
	LastSuccessfullyRunCompletionTime *metav1.Time `json:"lastSuccessfullyRunCompletionTime,omitempty"`
}

type ListJobs struct {
	Jobs  []DecoratedJob `json:"jobs"`
	Count int            `json:"count"`
}

type LastStatus struct {
	Type    string `json:"type"`
	Message string `json:"message,omitempty"`
}
