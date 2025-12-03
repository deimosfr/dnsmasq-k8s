package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// StartSupervisorService starts a supervisor service
// @Summary      Start a supervisor service
// @Description  Starts the specified service managed by supervisor
// @Tags         supervisor
// @Accept       json
// @Produce      json
// @Param        service  path      string  true  "Service Name"
// @Success      200      {object}  map[string]string
// @Failure      400      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /supervisor/{service}/start [post]
func (s *Server) StartSupervisorService(c *gin.Context) {
	serviceName := c.Param("service")
	if serviceName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "service name is required"})
		return
	}

	err := s.supervisorService.StartService(serviceName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "service started"})
}

// StopSupervisorService stops a supervisor service
// @Summary      Stop a supervisor service
// @Description  Stops the specified service managed by supervisor
// @Tags         supervisor
// @Accept       json
// @Produce      json
// @Param        service  path      string  true  "Service Name"
// @Success      200      {object}  map[string]string
// @Failure      400      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /supervisor/{service}/stop [post]
func (s *Server) StopSupervisorService(c *gin.Context) {
	serviceName := c.Param("service")
	if serviceName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "service name is required"})
		return
	}

	err := s.supervisorService.StopService(serviceName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "service stopped"})
}

// RestartSupervisorService restarts a supervisor service
// @Summary      Restart a supervisor service
// @Description  Restarts the specified service managed by supervisor
// @Tags         supervisor
// @Accept       json
// @Produce      json
// @Param        service  path      string  true  "Service Name"
// @Success      200      {object}  map[string]string
// @Failure      400      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /supervisor/{service}/restart [post]
func (s *Server) RestartSupervisorService(c *gin.Context) {
	serviceName := c.Param("service")
	if serviceName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "service name is required"})
		return
	}

	err := s.supervisorService.RestartService(serviceName)
	if err != nil {
		// Log the error to stdout so it appears in pod logs
		// The error message now includes the supervisorctl output
		println(fmt.Sprintf("ERROR: RestartSupervisorService failed: %v", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "service restarted"})
}
