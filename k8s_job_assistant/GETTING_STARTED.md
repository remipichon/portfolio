This document is intended for any users with access to a Kubernetes to quickly
deploy Kubernetes Job Assistant. 

Kubernetes Job Assistant getting started
========================================

At the end of this document you will have 
* KJA running in your cluster
* have access to KJA via a port forward
* created a few dummy Jobs for demo usage

Prerequisites : 
* Kubernetes 1.21+ (support `suspend` Job attribute)
* Kubectl installed
* you can kubectl into your cluster (your `~/kube/config` is configured accordingly)
* enough RBAC to create namespace, Deployments and Jobs
* your cluster can pull images from dockerhub.io 

> such cluster can obviously be a local one provided through Docker Desktop

For the purpose of the demo we will work in a namespace named `kja-demo`.

Deploy KJA via Kustomize 
```bash
kubectl apply -k kustomize/overlays/demo
```

Then create the port-forward to access it locally without worrying about Ingress,
SSL and DNS. 
```bash
kubectl port-forward -n kja-demo service/kube-job-assistant 8080:8080
```
visit [http://localhost:8080/](http://localhost:8080/)

> It can take a few moments for the image to be pulled and the pod to start.

You should see the list populated with dummy Jobs

![kja demo list with dummy jobs](doc/kja_demo.png)

Check the [USER_GUIDE.md](USER_GUIDE.md) to know what can be done. 

Tear down the demo with
```bash
kubectl delete -k kustomize/overlays/demo
```
