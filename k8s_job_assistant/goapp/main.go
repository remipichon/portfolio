package main

import (
	"github.com/gin-gonic/gin"
	"goapp/kube"
	"log"
	"net/http"
)

func main() {

	kubeService := kube.Service{}
	kubeService.InitClient()
	kubeService.SetJobAssistAnnotation("job-assistant")

	r := gin.Default()
	r.GET("/list", func(c *gin.Context) {
		handleList(c, &kubeService)
	})
	r.GET("/run/:namespace/:name", func(c *gin.Context) {
		handleRun(c, &kubeService)
	})
	//TODO handleStatus
	r.GET("/kill/:namespace/:name", func(c *gin.Context) {
		handleKill(c, &kubeService)
	})

	log.Println("Server starting on :8080...")
	log.Fatal(r.Run(":8080"))
}

type ListJobs struct {
	Jobs  []ResponseJob `json:"jobs"`
	Count int           `json:"count"`
}

// TODO decorate more as needed by the UI
type ResponseJob struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

func handleList(c *gin.Context, kubeService *kube.Service) {
	jobs, err := kubeService.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	returnJobs := make([]ResponseJob, 0, len(jobs))
	for _, job := range jobs {
		returnJobs = append(returnJobs, ResponseJob{
			Namespace: job.Namespace,
			Name:      job.Name,
		})
	}

	listJobs := ListJobs{
		Jobs:  returnJobs,
		Count: len(returnJobs),
	}

	c.JSON(http.StatusOK, listJobs)
}

func handleRun(c *gin.Context, kubeService *kube.Service) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	if namespace == "" || name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid path"})
		return
	}
	if err := kubeService.Run(namespace, name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func handleKill(c *gin.Context, kubeService *kube.Service) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	if namespace == "" || name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid path"})
		return
	}
	if err := kubeService.Kill(namespace, name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
