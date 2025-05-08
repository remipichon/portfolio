This document is intended for Ops/Dev to improve/fix KJA. 

It is recommended to read [ARCHITECTURE.md](ARCHITECTURE.md) before. 

# Contribute

Prerequisites :
* a running Kubernetes cluster configured through default `~/.kube/config`
* RBAC to create the resources listed in [kustomize/backend-testing](kustomize/backend-testing) 

Go backend which exposes a JSON Rest API
```bash
make backend
```
Listen on localhost:8080

> To understand more about the RBAC needed to test, check [set-up-the-least-privileges-sa](#set-up-the-least-privileges-sa) 

> Check [tear-down-the-least-privileges-sa](#tear-down-the-least-privileges-sa) once
> you are done with your testing 

React frontend which consumes the API
```bash
make frontend 
```
Listen on localhost:3000


## Backend testing

`internal/kube` covers the Job manager which directly interact with a live
Kubernetes cluster. Tests creates Jobs under the `kja-test-resources` namespace.

All resources are created with the labels K8s `"testing-labels": "under-test-k8s-job-assistant"`
for easy cleaning.

### Set up the least privileges SA

To simulate working in a _least privileges_ setup, we need : 
* a [ServiceAccount](kustomize/backend-testing/service_account.yaml) and a [token](kustomize/backend-testing/secret.yaml)
to generate a Kube Config (kja-sa-kubeconfig-test.yaml)
* a [ClusterRoleBinding](kustomize/backend-testing/cluster-role-binding.yaml) bind 
to the [ClusterRole](kustomize/base/cluster-role.yaml) used for the production setup 
* for the sake of test suite setup and test suite tear down, the test SA also can create/delete namespace but no more

To generate those K8s resources and the associated KubeConfig :
```bash 
kubectl apply -k kustomize/backend-testing
bash kustomize/backend-testing/seed_kubeconfig.sh > goapp/internal/kube/kja-sa-kubeconfig-test.yaml
```
> This is where you need a valid KubeConfig with enough permissions 
> the Make target `make setup-test` does it for you, `make backend` depends on it.

> This setup can not be done through the Go SetupSuite because of the RBAC needed
> to create the SA/token and ClusterRoleBinding. 

> It is possible `kubectl get secret kube-job-assistant-token -n kja-test-deploy`
> will fail if the cluster didn't have time to populate the secret with the token. 
> If it happens, run the command twice. Kustomize won't update existing resources
> and the bash script will be able to retrieve the already token populated by then. 

Then you can run the following test. 

### Run the tests
```bash
cd goapp
go test -v ./... -kubeconfig .kja-sa-kubeconfig-test.yaml
```
> `make backend-test` does the SA setup + run the test

Existing tests: 
```
cd goapp && go test -v ./... -kubeconfig ./kja-sa-kubeconfig-test.yaml
?   	goapp	[no test files]
?   	goapp/internal/handler	[no test files]
=== RUN   TestKubeService
Starting Job Manager tests in namespace kja-test-resources with labels map[testing-labels:under-test-k8s-job-assistant] and KJA annotation under-test-for-job-assist
=== RUN   TestKubeService/TestKillJob
    job_manager_test.go:327: Run has started
=== RUN   TestKubeService/TestKillJobNonExisting
=== RUN   TestKubeService/TestListJob
=== RUN   TestKubeService/TestListJobEmpty
=== RUN   TestKubeService/TestListJobStatusSuspended
=== RUN   TestKubeService/TestRunAfterKillJob
    job_manager_test.go:346: Run has started
    job_manager_test.go:354: Run has started
=== RUN   TestKubeService/TestRunJobAfterCompletion
    job_manager_test.go:288: Run first time
    job_manager_test.go:292: First run has started
    job_manager_test.go:295: First run has completed
    job_manager_test.go:297: Run second time
    job_manager_test.go:301: Second run has started, test is over
=== RUN   TestKubeService/TestRunJobAfterCreate
=== RUN   TestKubeService/TestRunJobNonExisting
=== RUN   TestKubeService/TestRunJobWhileRunning
    job_manager_test.go:309: Run first time
    job_manager_test.go:313: First run has started
    job_manager_test.go:315: Run second time (without waiting for first completion
=== RUN   TestKubeService/TestRunJobWithoutSuspend
=== RUN   TestKubeService/TestSetup
Wait for namespace to be deleted (timeout 30s) run 'go test -tear-down' to keep trying if it times out:  kja-test-resources
--- PASS: TestKubeService (49.09s)
    --- PASS: TestKubeService/TestKillJob (2.90s)
    --- PASS: TestKubeService/TestKillJobNonExisting (0.00s)
    --- PASS: TestKubeService/TestListJob (0.20s)
    --- PASS: TestKubeService/TestListJobEmpty (0.60s)
    --- PASS: TestKubeService/TestListJobStatusSuspended (0.80s)
    --- PASS: TestKubeService/TestRunAfterKillJob (18.58s)
    --- PASS: TestKubeService/TestRunJobAfterCompletion (5.39s)
    --- PASS: TestKubeService/TestRunJobAfterCreate (1.60s)
    --- PASS: TestKubeService/TestRunJobNonExisting (0.60s)
    --- PASS: TestKubeService/TestRunJobWhileRunning (1.80s)
    --- PASS: TestKubeService/TestRunJobWithoutSuspend (5.40s)
    --- PASS: TestKubeService/TestSetup (0.40s)
PASS
ok  	goapp/internal/kube	49.413s
?   	goapp/internal/model	[no test files]
?   	goapp/internal/service	[no test files]
```

### Keep resources
The test setups the needed resources and tear them down after testing. You can
disable the tearing down to investigate tests cases with :
```bash
go test ./... -keep-resources -kubeconfig .kja-sa-kubeconfig-test.yaml
```

### Tear down resources
If resources were kept, running the test again will get you the following error
```
Error:      	Received unexpected error:
                namespaces "kja-test-resources" already exists
Test:       	TestKubeService
Messages:   	failed to create namespace
 ```

Simply manually clean up the resources before running the tests suite
```bash
go test ./... -tear-down -kubeconfig .kja-sa-kubeconfig-test.yaml
```

This flag is a bit sketchy. To not clutter the test suite, it simply panics
once the tearing down is done, the expected logs looks like this
```
Running in tear-down mode because -tear-down
Wait for namespace to be deleted (timeout 30s) run 'go test -tear-down' to keep trying if it times out:  kja-test-resources
Tear down has run, you can now re-run the test without -tear-down
--- FAIL: TestKubeService (5.16s)
    proc.go:67: test panicked: unexpected call to os.Exit(0) during test
        goroutine 4 [running]:
```

### Tear down the least privileges SA

To clean up the tests SA, simply run 
```bash
make tear-down-test
```
