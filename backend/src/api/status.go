package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) GetStatus(c *gin.Context) {
	status := s.statusService.GetStatus()
	c.JSON(http.StatusOK, status)
}
