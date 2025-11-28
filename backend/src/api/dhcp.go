package api

import (
	"backend/src/services"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) GetLeases(c *gin.Context) {
	leases, err := s.dhcpService.GetLeases(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"leases": leases})
}

func (s *Server) GetReservations(c *gin.Context) {
	reservations, err := s.dhcpService.GetReservations(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"reservations": reservations})
}

func (s *Server) AddReservation(c *gin.Context) {
	var json struct {
		MACAddress string `json:"mac_address"`
		IPAddress  string `json:"ip_address"`
		Hostname   string `json:"hostname"`
		Comment    string `json:"comment"`
	}

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := net.ParseMAC(json.MACAddress); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid MAC address"})
		return
	}

	if net.ParseIP(json.IPAddress) == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid IP address"})
		return
	}

	err := s.dhcpService.AddReservation(c.Request.Context(), json.MACAddress, json.IPAddress, json.Hostname, json.Comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (s *Server) UpdateReservation(c *gin.Context) {
	var json struct {
		Old services.DHCPReservation `json:"old"`
		New services.DHCPReservation `json:"new"`
	}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := s.dhcpService.UpdateReservation(c.Request.Context(), json.Old, json.New); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (s *Server) DeleteReservation(c *gin.Context) {
	var json services.DHCPReservation
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := s.dhcpService.DeleteReservation(c.Request.Context(), json); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (s *Server) UpdateLease(c *gin.Context) {
	var json struct {
		Old services.DHCPLease `json:"old"`
		New services.DHCPLease `json:"new"`
	}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := s.dhcpService.UpdateLease(c.Request.Context(), json.Old, json.New); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (s *Server) DeleteLease(c *gin.Context) {
	var json services.DHCPLease
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := s.dhcpService.DeleteLease(c.Request.Context(), json); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
