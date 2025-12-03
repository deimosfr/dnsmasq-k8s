package api

//go:generate go run ../../cmd/syncversion/main.go

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const Version = "1.2.0"

// GetVersion returns the current version of the application
// @Summary      Get application version
// @Description  Returns the current version of the application
// @Tags         version
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /version [get]
func (s *Server) GetVersion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"version": Version})
}
