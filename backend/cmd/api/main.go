package main

import (
	"backend/src/api"
	"backend/src/services"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// Create a new Gin router without default middleware
	r := gin.New()

	// Add recovery middleware
	r.Use(gin.Recovery())

	// Add custom logger middleware that skips /api/v1/status
	r.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/api/v1/status"},
	}))

	// Create a new Kubernetes clientset.
	config, err := clientcmd.BuildConfigFromFlags("", "")
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// Get the namespace from the environment variable.
	namespace := os.Getenv("POD_NAMESPACE")
	if namespace == "" {
		namespace = "default"
	}

	fmt.Printf("INFO: Starting application in namespace: %s\n", namespace)
	fmt.Printf("INFO: Starting application version: %s\n", api.Version)

	// Create a new config service.
	configService := services.NewConfigService(clientset, namespace)

	// Create a new dhcp service.
	dhcpService := services.NewDHCPService(clientset, namespace, configService)

	// Create a new status service.
	statusService := services.NewStatusService()

	// Create a new supervisor service.
	supervisorService := services.NewSupervisorService()

	// Create a new server.
	server := api.NewServer(configService, dhcpService, statusService, supervisorService)

	r.StaticFS("/static", http.Dir("../../../frontend/src"))

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/static/index.html")
	})

	// Register the API endpoints.
	v1 := r.Group("/api/v1")
	{
		v1.GET("/config", server.GetConfig)
		v1.PUT("/config", server.UpdateConfig)
		v1.POST("/dns/entries", server.AddDNSEntry)
		v1.GET("/dns/entries", server.GetDNSEntries)
		v1.DELETE("/dns/entries", server.DeleteDNSEntry)
		v1.PUT("/dns/entries", server.UpdateDNSEntry)
		v1.GET("/dhcp/leases", server.GetLeases)
		v1.PUT("/dhcp/leases", server.UpdateLease)
		v1.DELETE("/dhcp/leases", server.DeleteLease)
		v1.GET("/dhcp/reservations", server.GetReservations)
		v1.POST("/dhcp/reservations", server.AddReservation)
		v1.PUT("/dhcp/reservations", server.UpdateReservation)
		v1.DELETE("/dhcp/reservations", server.DeleteReservation)
		v1.GET("/status", server.GetStatus)
		v1.POST("/supervisor/:service/start", server.StartSupervisorService)
		v1.POST("/supervisor/:service/stop", server.StopSupervisorService)
		v1.POST("/supervisor/:service/restart", server.RestartSupervisorService)
		v1.GET("/version", server.GetVersion)
	}

	r.Run(":8080")
}
