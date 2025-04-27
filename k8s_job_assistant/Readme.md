Kubernetes Job Assistant 
========================

Sometimes you have to run a K8s job to perform a predefined tasks at the time of
your choice : 
* perform a data migration after a 3rd party provided emailed you
* trigger a background checkup after an external event you can't automate happened
* perform an export for the marketing team 

The job is already defined, it has been tested and all, but you need to manually
control when it is executed. 

This tiny tool provides several ways for non-technical users to be autonomous. 
Because the ops team doesn't want to be the one running `kubectl run` and you 
don't want to give access to ArgoCD to the Product Owner or the Marketing boss. 



# Specs 


* create job with a given annotation 
  * annotation is configurable, but there is a default
* a UI 
  * auth via any Auth0 provider
    * basic roles from the provider : admin or viewer 
  * list available job and their latest K8s events
  * run a job, get its logs
  * kill a job
  * run a job with on-the-spot parameters 

# infra

* stateless 
  * only based on annotation  
* IaC
  * Deployment 
    * distro less container 
    * via Flux : can override scaling, affinity, resources 
  * Service 
    * to be connected to your ingress
    * port forward for the sake of the demo 
  * RBAC 
    * get/create/delete Jobs configured namespace
    * namepsace configurable via Flux
    * ? can I just "run" existing Job ? Does it have to be CronJob ? 
* app
  * Go orchestrator with Kube SDK 
  * React UI embedded and served by Go server 
  * (just one container)
  * exposed comprehensive REST API 
    *  ? which security ? 
* docs
  * OpenAPI
  * Go Code
  * IaC (Kustomize)

# Go app


# UI
* cookie to store auth token 



# demo setup 

* local Kube 
* oAuth avec Github (mais rester generique)


# perf
* naive implementation
 * each "Refresh" perform a scan for matching jobs
 * 
