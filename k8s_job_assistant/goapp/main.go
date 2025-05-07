package main

import (
	"github.com/gin-gonic/gin"
	"goapp/internal/handler"
	"goapp/internal/kube"
	"goapp/internal/service"
	"log"
	"net/http"
	"os"
)

func main() {

	router := gin.Default()

	// Setup Job Manager, Service and http Handler
	jobManager := kube.NewJobManager(kube.InitKubeClient(), "job-assistant")
	jobService := service.NewJobService(jobManager)
	handler.DecorateRouterWithJobHandlers(router, jobService)

	//Serve Static React app
	if _, err := os.Stat("ui/index.html"); err == nil {
		router.LoadHTMLFiles("ui/index.html")
		router.Static("/assets", "./ui/assets")
		router.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.html", gin.H{})
		})
	} else {
		log.Println("ui/index.html not found, not serving static file")
	}

	// Start server
	err := router.Run(":8080")
	if err != nil {
		log.Fatal(router.Run(":8080"))
	}
	log.Println("Server starting on :8080...")
}
