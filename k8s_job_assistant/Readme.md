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

Check out [GETTING_STARTED](GETTING_STARTED.md) for a quick easy demo setup.

Check out [USER_GUIDE.md](USER_GUIDE.md) for end user usage of the UI.

Check out [RUNBOOK.md](RUNBOOK.md) to deploy KJA in any Kubernetes cluster.

Check out [ARCHITECTURE.md](ARCHITECTURE.md) to understand what KJA is made of. 

Check out [CONTRIBUTE.md](CONTRIBUTE.md) to work on KJA in a local Kubernetes cluster.
