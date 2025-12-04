package main

import (
	"backend/src/api"
	"backend/src/services"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	_ "backend/src/docs" // Import generated docs
)

// @title           Dnsmasq K8s API
// @version         1.0
// @description     API for managing Dnsmasq in Kubernetes
// @termsOfService  http://swagger.io/terms/

// @contact.name    API Support
// @contact.url     http://www.swagger.io/support
// @contact.email   support@swagger.io

// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

func main() {
	// Get configuration from environment variables
	namespace := os.Getenv("POD_NAMESPACE")
	if namespace == "" {
		namespace = "default"
	}

	webPort := os.Getenv("WEB_PORT")
	if webPort == "" {
		webPort = "8080"
	}

	apiPort := os.Getenv("API_PORT")
	if apiPort == "" {
		apiPort = "8081"
	}

	fmt.Printf("INFO: Starting application in namespace: %s\n", namespace)
	fmt.Printf("INFO: Starting application version: %s\n", api.Version)
	fmt.Printf("INFO: Web server listening on port: %s\n", webPort)
	fmt.Printf("INFO: API server listening on port: %s\n", apiPort)

	// Create a new Kubernetes clientset.
	config, err := clientcmd.BuildConfigFromFlags("", "")
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// Create services
	configService := services.NewConfigService(clientset, namespace)
	dhcpService := services.NewDHCPService(clientset, namespace, configService)
	statusService := services.NewStatusService()
	supervisorService := services.NewSupervisorService()
	server := api.NewServer(configService, dhcpService, statusService, supervisorService)

	// --- Server Setup ---
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/api/v1/status"},
	}))
	// CORS might not be strictly necessary for same-origin, but good to keep if we want to allow external access or dev mode
	router.Use(api.CORSMiddleware())

	// API Routes
	v1 := router.Group("/api/v1")
	{
		v1.GET("/config", server.GetConfig)
		v1.GET("/config/tags", server.GetTags)
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

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Web Routes
	router.StaticFS("/static", http.Dir("../../../frontend/src"))

	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/static/index.html")
	})

	// Serve env.js to configure API URL
	router.GET("/env.js", func(c *gin.Context) {
		// Since we are now serving from the same origin, API_URL can be empty (relative path)
		// or we can explicitly set it to the current origin if needed.
		// Empty string usually implies "same origin" in our frontend logic if we update it,
		// or we can just let it be empty and frontend uses relative paths?
		// The frontend code uses `${window.env.API_URL}/api/v1/...`
		// If API_URL is "", it becomes `/api/v1/...` which is correct for same origin.

		c.Header("Content-Type", "application/javascript")
		c.String(http.StatusOK, `window.env = { API_URL: "" };`)
	})

	// Run server
	if err := router.Run(":" + webPort); err != nil {
		panic(fmt.Sprintf("Server failed: %v", err))
	}
}
