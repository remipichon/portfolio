Kubernetes Job Assistant 
========================

Sometimes you have to run a K8s job to perform a predefined tasks at the time of
your choice : 
* perform a data migration after a contractor emailed you
* trigger a background checkup after an external event you can't predict happened
* perform an export for the marketing team 

The K8s job is already defined, it has been implemented and tested, but you need
to manually control when it is executed. 

This tiny tool provides a way for non-technical users to be autonomous in running
those K9s Jobs. Because the ops team doesn't want to be the one running 
`kubectl run` and you don't want to give access to ArgoCD to the Product Owner 
or the Marketing boss. 

With Kubernetes Job Assistant, allowed members of your organization can
* list Kubernetes Jobs based on the annotation `job-assistant:true` (customizable)
* run a Job
* kill a Job
* check basic Job stats


# Deploy Kubernetes Job Assistant in your cluster

Prerequisites : 
* Kubernetes 1.21+ (support `suspend` Job attribute)

> not implemented yet, see [Todo list](Todo.md)


# Configure your existing Jobs

Let's say you have an existing carefully crafted Job. To delegate its lifecycle
to KJA :
* Add the annotation `metadata.annotations.job-assistant: enable`
  But deploying it through your CI would run it right away, to prevent this:
* Add the `spec.suspend: true
```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: my-existing-job
  namespace : kja-demo
  annotations:
    job-assistant: enable   # to delegate the lifecycle to KJA
spec:
  suspend: true             # to delay the run
  template:
    spec:
      containers:
      - name: some-awesomely-tested-job
        image: busybox
        command: ["sh", "-c", "echo This is my awesome task which lasts 5 seconds!; sleep 5; echo This is the end of my awesome task"]
      restartPolicy: Never
```
Now your CI/CD can take over.


You can also patch an existing completed Job for KJA to take over its lifecycle.
```bash
kubectl annotate job my-existing-job job-assistant=enable --overwrite -n kja-demo
```
KJA will take over, you will be able to run the Job from the UI.


> Please note that KJA will tweak the `suspend` attribute value during the Job lifecycle. 
> If your CI/CD can't ignore the field, it's okay if it keeps setting it to `true`
> as it won't kill running Jobs but it could prevent Kubernetes from starting the pods
> in time.

# Contribute 

Prerequisites : 
* a running Kubernetes cluster configured through default `~/.kube/config`
* RBAC to create/delete Jobs in any namespace

Go backend which exposes a JSON Rest API
```bash
make backend
```
Listen on localhost:8080


React frontend which consumes the API
```bash
make frontend 
```
Listen on localhost:3000


## Backend testing

`internal/kube` covers the Job manager which directly interact with a live
Kubernetes cluster. Tests creates Jobs under the `kja-test-namespace` namespace.
One test also creates a Job in the `default` namespace.

All resources are created with the labels K8s `"testing-labels": "under-test-k8s-job-assistant"`
for easy cleaning. 

Prerequisites :
* a running Kubernetes cluster configured through default `~/.kube/config`
* RBAC 
  * create/delete namespace `kja-test-namespace`
  * crete/delete/patch Jobs in namespace `kja-test-namespace` and `default`

### Run the tests
```bash
cd goapp
go test ./...
```

### Keep resources
The tests setups the needed resources and tear them down after testing. You can 
disable the tearing down to investigate tests cases with 
```bash
go test ./... -keep-resources
```

### Tear down resources
If resources were kept, you will get the following error
If you get this error
```
Error:      	Received unexpected error:
                namespaces "kja-test-namespace" already exists
Test:       	TestKubeService
Messages:   	failed to create namespace
 ```

Simply manually clean up the resources before running the tests suite
```bash
go test ./... -tear-down
```

This flag is a bit sketchy. To not cluter the test suite, it simply panics
once the tearing down is done, the expected logs looks like this
```
Running in tear-down mode because -tear-down
Wait for namespace to be deleted (timeout 30s) run 'go test -tear-down' to keep trying if it times out:  kja-test-namespace
Tear down has run, you can now re-run the test without -tear-down
--- FAIL: TestKubeService (5.16s)
    proc.go:67: test panicked: unexpected call to os.Exit(0) during test
        goroutine 4 [running]:
```