package handler

import (
	"github.com/gin-gonic/gin"
	"goapp/internal/model"
	"goapp/internal/service"
	"net/http"
)

func DecorateRouterWithJobHandlers(router *gin.Engine, jobSvc service.JobService) {
	router.GET("/list", func(c *gin.Context) {
		jobs, err := jobSvc.ListDecoratedJobs()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		listJobs := model.ListJobs{
			Jobs:  jobs,
			Count: len(jobs),
		}

		c.JSON(http.StatusOK, listJobs)
	})

	router.GET("/run/:namespace/:name", func(c *gin.Context) {
		namespace := c.Param("namespace")
		name := c.Param("name")
		if namespace == "" || name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid path"})
			return
		}
		if err := jobSvc.Run(namespace, name); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusOK)
	})

	router.GET("/kill/:namespace/:name", func(c *gin.Context) {
		namespace := c.Param("namespace")
		name := c.Param("name")
		if namespace == "" || name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid path"})
			return
		}
		if err := jobSvc.Kill(namespace, name); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusOK)
	})
}
