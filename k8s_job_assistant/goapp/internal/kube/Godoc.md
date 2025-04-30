# kube
--
    import "."

Package kube implements service methods to manage Kubernetes Jobs

## Usage

#### func  InitKubeClient

```go
func InitKubeClient() *kubernetes.Clientset
```
InitClient instantiate a Kubernetes client based on local kubeconfig.

#### type JobAlreadyRunningError

```go
type JobAlreadyRunningError struct {
}
```


#### func (*JobAlreadyRunningError) Error

```go
func (e *JobAlreadyRunningError) Error() string
```

#### type JobManager

```go
type JobManager interface {
	List() ([]batchv1.Job, error)
	Run(namespace, jobName string) error
	Kill(namespace, jobName string) error
	Status(namespace, jobName string) (error, *batchv1.JobStatus)
}
```


#### func  NewJobManager

```go
func NewJobManager(kubeClient *kubernetes.Clientset, jobAssistAnnotation string) JobManager
```
