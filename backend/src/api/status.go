package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetStatus returns the current status of the application
// @Summary      Get application status
// @Description  Returns the current status of the application
// @Tags         status
// @Produce      json
// @Success      200  {object}  services.Status
// @Router       /status [get]
func (s *Server) GetStatus(c *gin.Context) {
	status := s.statusService.GetStatus()
	c.JSON(http.StatusOK, status)
}
