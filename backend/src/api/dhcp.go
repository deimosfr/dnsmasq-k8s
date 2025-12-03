package api

import (
	"backend/src/services"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AddReservationRequest struct {
	MACAddress string `json:"mac_address"`
	IPAddress  string `json:"ip_address"`
	Hostname   string `json:"hostname"`
	Comment    string `json:"comment"`
}

type UpdateReservationRequest struct {
	Old services.DHCPReservation `json:"old"`
	New services.DHCPReservation `json:"new"`
}

type UpdateLeaseRequest struct {
	Old services.DHCPLease `json:"old"`
	New services.DHCPLease `json:"new"`
}

// GetLeases returns all DHCP leases
// @Summary      Get DHCP leases
// @Description  Returns all DHCP leases
// @Tags         dhcp
// @Produce      json
// @Success      200  {object}  map[string][]services.DHCPLease
// @Failure      500  {object}  map[string]string
// @Router       /dhcp/leases [get]
func (s *Server) GetLeases(c *gin.Context) {
	leases, err := s.dhcpService.GetLeases(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"leases": leases})
}

// GetReservations returns all DHCP reservations
// @Summary      Get DHCP reservations
// @Description  Returns all DHCP reservations
// @Tags         dhcp
// @Produce      json
// @Success      200  {object}  map[string][]services.DHCPReservation
// @Failure      500  {object}  map[string]string
// @Router       /dhcp/reservations [get]
func (s *Server) GetReservations(c *gin.Context) {
	reservations, err := s.dhcpService.GetReservations(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"reservations": reservations})
}

// AddReservation adds a new DHCP reservation
// @Summary      Add DHCP reservation
// @Description  Adds a new DHCP reservation
// @Tags         dhcp
// @Accept       json
// @Produce      json
// @Param        reservation  body      AddReservationRequest  true  "DHCP Reservation"
// @Success      200          {object}  map[string]string
// @Failure      400          {object}  map[string]string
// @Failure      500          {object}  map[string]string
// @Router       /dhcp/reservations [post]
func (s *Server) AddReservation(c *gin.Context) {
	var json AddReservationRequest

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

// UpdateReservation updates a DHCP reservation
// @Summary      Update DHCP reservation
// @Description  Updates a DHCP reservation
// @Tags         dhcp
// @Accept       json
// @Produce      json
// @Param        reservation  body      UpdateReservationRequest  true  "DHCP Reservation Update"
// @Success      200          {object}  map[string]string
// @Failure      400          {object}  map[string]string
// @Failure      500          {object}  map[string]string
// @Router       /dhcp/reservations [put]
func (s *Server) UpdateReservation(c *gin.Context) {
	var json UpdateReservationRequest
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

// DeleteReservation deletes a DHCP reservation
// @Summary      Delete DHCP reservation
// @Description  Deletes a DHCP reservation
// @Tags         dhcp
// @Accept       json
// @Produce      json
// @Param        reservation  body      services.DHCPReservation  true  "DHCP Reservation"
// @Success      200          {object}  map[string]string
// @Failure      400          {object}  map[string]string
// @Failure      500          {object}  map[string]string
// @Router       /dhcp/reservations [delete]
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

// UpdateLease updates a DHCP lease
// @Summary      Update DHCP lease
// @Description  Updates a DHCP lease
// @Tags         dhcp
// @Accept       json
// @Produce      json
// @Param        lease  body      UpdateLeaseRequest  true  "DHCP Lease Update"
// @Success      200    {object}  map[string]string
// @Failure      400    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Router       /dhcp/leases [put]
func (s *Server) UpdateLease(c *gin.Context) {
	var json UpdateLeaseRequest
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

// DeleteLease deletes a DHCP lease
// @Summary      Delete DHCP lease
// @Description  Deletes a DHCP lease
// @Tags         dhcp
// @Accept       json
// @Produce      json
// @Param        lease  body      services.DHCPLease  true  "DHCP Lease"
// @Success      200    {object}  map[string]string
// @Failure      400    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Router       /dhcp/leases [delete]
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
