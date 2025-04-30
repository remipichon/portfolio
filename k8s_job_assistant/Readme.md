Kubernetes Job Assistant 
========================

Sometimes you have to run a K8s job to perform a predefined tasks at the time of
your choice : 
* perform a data migration after a contractor emailed you
* trigger a background checkup after an external event you can't predict happened
* perform an export for the marketing team 

The job is already defined, it has been tested and all, but you need to manually
control when it is executed. 

This tiny tool provides a way for non-technical users to be autonomous. 
Because the ops team doesn't want to be the one running `kubectl run` and you 
don't want to give access to ArgoCD to the Product Owner or the Marketing boss. 

With Kubernetes Job Assistant, allowed members of your organization can
* list Jobs based on the annotation `job-assistant:true` (customizable)
* run a Job
* kill a Job
* check basic Job stats
