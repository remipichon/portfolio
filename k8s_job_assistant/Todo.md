Todo
=======

Project under work, mishmash tasks

* Go backend
  * unit test for internal/service which mocked Kube JobManager
  * export the annotation to be configured via ENV from a configmap 
  * add a job-manager unit test to deal with CI/CD always setting suspend to true + better explain the behaviour and race condition 
  * OpenAPI auto generated doc (care for types and examples)

* React app
  * better handle auto refresh 
  * nicer CSS 

* Github actions CI (gotest, build docker, push to docker hub)

* prod ready setup
  * multi layer Docker image with Goapp + React (distroless)
  * Github Action CD to build, validate and push to public Docker registry
  * K8s Deployment + Service to deploy public image
  * minimum working RBAC
  * Make to deploy to Kube
  * Make to setup portforward (demo purpose, to skip having a proper Ingress/DNS/SSL)

* feature 
  * oAuth via Github app
    * keep it generic to support any oAuth provider 
    * support two roles : Admin and Read
    * simple integration into the ReactApp (two roles to support and oAuth flow)

* going further
  * deep dive in scaling and how would it behave with hundreds oj Jobs
  * write benchmark test
