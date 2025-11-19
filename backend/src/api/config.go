package api

import (
	"backend/src/services"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) GetConfig(c *gin.Context) {
	config, err := s.configService.GetConfig(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.String(http.StatusOK, config)
}

func (s *Server) UpdateConfig(c *gin.Context) {
	var json struct {
		Config string `json:"config"`
	}

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := s.configService.UpdateConfig(c.Request.Context(), json.Config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (s *Server) AddDNSEntry(c *gin.Context) {
	var json struct {
		Type   string `json:"type"`
		Domain string `json:"domain"`
		Value  string `json:"value"`
	}

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate record type - only support A, CNAME, TXT
	validTypes := map[string]bool{
		"address": true,
		"cname":   true,
		"txt":     true,
	}
	if !validTypes[json.Type] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported DNS record type. Only A, CNAME, and TXT records are supported."})
		return
	}

	// Validate IP address for A records
	if json.Type == "address" {
		if net.ParseIP(json.Value) == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid IP address"})
			return
		}
	}

	// Validate domain name for CNAME records
	if json.Type == "cname" {
		if json.Value == "" || len(json.Value) > 253 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid domain name for CNAME record"})
			return
		}
	}

	err := s.configService.AddDNSEntry(c.Request.Context(), json.Type, json.Domain, json.Value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (s *Server) GetDNSEntries(c *gin.Context) {
	entries, err := s.configService.GetDNSEntries(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, entries)
}

func (s *Server) DeleteDNSEntry(c *gin.Context) {
	var json services.DNSEntry
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := s.configService.DeleteDNSEntry(c.Request.Context(), json)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (s *Server) UpdateDNSEntry(c *gin.Context) {
	var json struct {
		Old services.DNSEntry `json:"old"`
		New services.DNSEntry `json:"new"`
	}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := s.configService.UpdateDNSEntry(c.Request.Context(), json.Old, json.New)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
