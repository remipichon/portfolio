This document is intended for Ops to deploy and configure KJA in a K8s cluster. 

# Deploy Kubernetes Job Assistant in your cluster

Prerequisites :
* Kubernetes 1.21+ (support `suspend` Job attribute)

Add your own Kustomization from [kustomize/overlays](kustomize/overlays). You can 
easily configure
* the namespace where KJA is deployed
* COMING : the annotation KJA use to take ownership of Job
> You can get inspiration from the [demo](kustomize/overlays/demo)

Add your Ingress setup to expose it through your existing setup. The Service port
is 8080 and defined as a ClusterIP.
You can tweak it in [kustomize/base/service.yaml](kustomize/base/service.yaml)

Deploy KJA via Kustomize
```bash
kubectl apply -k kustomize/overlays/your_kustomization
```

# Configure your existing Jobs

Let's say you have an existing carefully crafted Job. To delegate its lifecycle
to KJA :
* Add the annotation `metadata.annotations.job-assistant: enable`
> the annotation can be changed through the Kustomization, see above

Deploying it through your CI would run it right away, to prevent this:
* Add the `spec.suspend: true

For example, with the values defined in the [demo](kustomize/overlays/demo) : 
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
