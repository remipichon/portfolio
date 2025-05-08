package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"goapp/internal/handler"
	"goapp/internal/kube"
	"goapp/internal/service"
	"k8s.io/client-go/util/homedir"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {

	// Optional flag for kubeconfig
	var kubeconfigPath string
	if home := homedir.HomeDir(); home != "" {
		flag.StringVar(&kubeconfigPath, "kubeconfig", filepath.Join(home, ".kube", "config"),
			"(optional) absolute path to the kubeconfig file")
	} else {
		flag.StringVar(&kubeconfigPath, "kubeconfig", "",
			"(optional) absolute path to the kubeconfig file")
	}
	flag.Parse()

	router := gin.Default()

	// Setup Job Manager, Service and http Handler
	jobManager := kube.NewJobManager(kube.InitKubeClient(kubeconfigPath), "job-assistant")
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
