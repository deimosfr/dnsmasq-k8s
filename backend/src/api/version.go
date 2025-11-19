package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetVersion returns the current version of the application
func (s *Server) GetVersion(c *gin.Context) {
	// This could be injected at build time, but for now hardcoding matches the Chart.yaml
	version := "1.0.0"
	c.JSON(http.StatusOK, gin.H{"version": version})
}
