# Prerequisites

* Kubernetes 1.21+ (requires `suspend` Job attribute)

* demo namespace
```
kubectl create namespace  kja-demo
```

# Prepare Jobs

Let's say you have an existing carefully crafted Job, deploying it through your
CI would run it right away. To delegate its lifecycle to KJA : 
* Add the `spec.suspend: true` to delay the run
* Add the annotation `metadata.annotations.job-assistant: enable`
```yaml
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
        command: ["sh", "-c", "echo This is my awesome task which lasts 5 seconds!; sleep 5; This is the end of my awesome task"]
      restartPolicy: Never
```
Now your CI can take over


!!!
? how to deal with `suspend` being removed by KJA ?
!!!

# Deploy the Kubernetes Job Assistant UI



