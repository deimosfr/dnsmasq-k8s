package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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

func (s *Server) RestartSupervisorService(c *gin.Context) {
	serviceName := c.Param("service")
	if serviceName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "service name is required"})
		return
	}

	err := s.supervisorService.RestartService(serviceName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "service restarted"})
}
