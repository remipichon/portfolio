# kube
--
    import "."

Package kube implements service methods to manage Kubernetes Jobs

## Usage

#### type JobAlreadyRunningError

```go
type JobAlreadyRunningError struct {
}
```


#### func (*JobAlreadyRunningError) Error

```go
func (e *JobAlreadyRunningError) Error() string
```

#### type Service

```go
type Service struct {
}
```

Service provides helper methods to interact with Kubernetes Jobs.

#### func (*Service) InitClient

```go
func (ks *Service) InitClient()
```
InitClient instantiate a Kubernetes client based on local kubeconfig.

#### func (*Service) Kill

```go
func (ks *Service) Kill(namespace, jobName string) error
```
Kill suspends the Job and delete all of its running pod.

TODO kill pods (to keep them and their logs) instead of deleting pods

#### func (*Service) List

```go
func (ks *Service) List() ([]batchv1.Job, error)
```
List lists Jobs with annotation 'job-assistant' set to true on any namespace.

#### func (*Service) Run

```go
func (ks *Service) Run(namespace, jobName string) error
```
Run runs a Job, fails if already running, handle Suspend:true and clean
re-create when needed.

#### func (*Service) Status

```go
func (ks *Service) Status(namespace, jobName string) (error, *batchv1.JobStatus)
```
Status returns the full Kubernetes status of job, without any decoration.
