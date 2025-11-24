package api

//go:generate go run ../../cmd/syncversion/main.go

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const Version = "1.1.0"

// GetVersion returns the current version of the application
func (s *Server) GetVersion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"version": Version})
}
