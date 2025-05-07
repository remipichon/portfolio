This document is intended for Ops/Dev to improve/fix KJA. 

It is recommended to read [ARCHITECTURE.md](ARCHITECTURE.md) before. 

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
