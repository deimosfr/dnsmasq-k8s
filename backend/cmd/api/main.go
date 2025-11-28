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

	// --- API Server ---
	apiRouter := gin.New()
	apiRouter.Use(gin.Recovery())
	apiRouter.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/api/v1/status"},
	}))
	apiRouter.Use(api.CORSMiddleware())

	v1 := apiRouter.Group("/api/v1")
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

	// --- Web Server ---
	webRouter := gin.New()
	webRouter.Use(gin.Recovery())
	webRouter.Use(gin.Logger())

	webRouter.StaticFS("/static", http.Dir("../../../frontend/src"))

	webRouter.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/static/index.html")
	})

	// Serve env.js to configure API URL
	webRouter.GET("/env.js", func(c *gin.Context) {
		// In a real scenario, this might be dynamic based on external URL
		// For now, we assume the API is on the same host but different port
		// If accessed via k8s service/ingress, this logic might need adjustment
		// But for local/docker, this works.
		// Actually, for browser access, we need the public URL.
		// Let's assume relative URL doesn't work across ports.
		// We can try to infer from request host, but port is different.
		// Let's just set it to empty string if we want to use relative path (same origin),
		// but here we are different origin (port).
		// So we need absolute URL.
		// Or we can use a proxy? No, user asked to dissociate.

		// Simple approach: Use the same hostname as the request, but change the port.
		// host := c.Request.Host // e.g. localhost:8080
		// We need to strip port and add API port
		// This is a bit hacky for production behind ingress, but works for direct access.
		// For production, we might want an env var for PUBLIC_API_URL.

		publicApiUrl := os.Getenv("PUBLIC_API_URL")
		if publicApiUrl == "" {
			// Fallback to constructing from request
			// This assumes http (not https) if not specified, which might be wrong.
			// But for internal/local it's fine.
			// scheme := "http"
			// if c.Request.TLS != nil {
			// 	scheme = "https"
			// }
			// If behind a proxy, X-Forwarded-Proto might be needed.

			// Let's just return a script that constructs it client side if possible?
			// No, client side doesn't know the API port unless we tell it.
			// So we must tell it the API port.

			// We'll just inject the API_PORT into the JS and let JS construct the URL
			// using window.location.hostname
			c.Header("Content-Type", "application/javascript")
			c.String(http.StatusOK, fmt.Sprintf(`window.env = { API_PORT: "%s", API_URL: "" };
if (!window.env.API_URL) {
	window.env.API_URL = window.location.protocol + "//" + window.location.hostname + ":%s";
}`, apiPort, apiPort))
			return
		}

		c.Header("Content-Type", "application/javascript")
		c.String(http.StatusOK, fmt.Sprintf(`window.env = { API_URL: "%s" };`, publicApiUrl))
	})

	// Run servers in goroutines
	go func() {
		if err := apiRouter.Run(":" + apiPort); err != nil {
			panic(fmt.Sprintf("API server failed: %v", err))
		}
	}()

	if err := webRouter.Run(":" + webPort); err != nil {
		panic(fmt.Sprintf("Web server failed: %v", err))
	}
}
