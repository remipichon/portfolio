package main

import (
	"github.com/gin-gonic/gin"
	"goapp/internal/handler"
	"goapp/internal/kube"
	"goapp/internal/service"
	"log"
)

func main() {

	router := gin.Default()

	// Setup Job Manager, Service and http Handler
	jobManager := kube.NewJobManager(kube.InitKubeClient(), "job-assistant")
	jobService := service.NewJobService(jobManager)
	handler.DecorateRouterWithJobHandlers(router, jobService)

	// Serve Static React app
	//r.LoadHTMLFiles("ui/index.html")
	//r.GET("/", func(c *gin.Context) {
	//	c.HTML(http.StatusOK, "index.html", gin.H{})
	//})

	// Start server
	err := router.Run(":8080")
	if err != nil {
		log.Fatal(router.Run(":8080"))
	}
	log.Println("Server starting on :8080...")
}
