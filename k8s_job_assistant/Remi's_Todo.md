Todo
=======

Project under work, mishmash tasks

* Go backend
  * unit test for internal/service which mocked Kube JobManager
  * export the annotation to be configured via ENV from a configmap and exposed through the Kustomization 
  * add a job-manager unit test to deal with a CI/CD which would always set suspend to true + better explain the behaviour and race condition 
  * OpenAPI auto generated doc (care for types and examples)
  * GoDoc 
  * add security layer with JWT auth 
  * transfer code docs from ARCHITECTURE.md into the go code to be dealt with via Godoc
  * **NEXT** [GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.

* Github actions CI 
  * gotest : discuss how to run the kube manager test in a CI (needs a valid Kube somewhere)
  * gotest : at least run the job_service test (with mocked job_manager)
  * build and publish according to Github release / git tag  

* feature 
  * oAuth via Github app
    * keep it generic to support any oAuth provider 
    * support two roles : Admin and Read
    * simple integration into the ReactApp (two roles to support and oAuth flow)

* going further
  * deep dive in scaling and how would it behave with hundreds of Jobs
  * write benchmark test
