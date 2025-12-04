package api

import (
	"backend/src/services"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UpdateConfigRequest struct {
	Config string `json:"config"`
}

type AddDNSEntryRequest struct {
	Type    string `json:"type"`
	Domain  string `json:"domain"`
	Value   string `json:"value"`
	Comment string `json:"comment"`
}

type UpdateDNSEntryRequest struct {
	Old services.DNSEntry `json:"old"`
	New services.DNSEntry `json:"new"`
}

// GetConfig returns the current dnsmasq configuration
// @Summary      Get configuration
// @Description  Returns the current dnsmasq configuration
// @Tags         config
// @Produce      text/plain
// @Success      200  {string}  string
// @Failure      500  {object}  map[string]string
// @Router       /config [get]
func (s *Server) GetConfig(c *gin.Context) {
	config, err := s.configService.GetConfig(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.String(http.StatusOK, config)
}

// GetTags returns all available tags from the configuration
// @Summary      Get tags
// @Description  Returns all available tags from the configuration
// @Tags         config
// @Produce      json
// @Success      200  {array}   string
// @Failure      500  {object}  map[string]string
// @Router       /config/tags [get]
func (s *Server) GetTags(c *gin.Context) {
	tags, err := s.configService.GetTags(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tags": tags})
}

// UpdateConfig updates the dnsmasq configuration
// @Summary      Update configuration
// @Description  Updates the dnsmasq configuration
// @Tags         config
// @Accept       json
// @Produce      json
// @Param        config  body      UpdateConfigRequest  true  "Configuration"
// @Success      200     {object}  map[string]string
// @Failure      400     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /config [put]
func (s *Server) UpdateConfig(c *gin.Context) {
	var json UpdateConfigRequest

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

// AddDNSEntry adds a new DNS entry
// @Summary      Add DNS entry
// @Description  Adds a new DNS entry (A, CNAME, or TXT)
// @Tags         dns
// @Accept       json
// @Produce      json
// @Param        entry  body      AddDNSEntryRequest  true  "DNS Entry"
// @Success      200    {object}  map[string]string
// @Failure      400    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Router       /dns/entries [post]
func (s *Server) AddDNSEntry(c *gin.Context) {
	var json AddDNSEntryRequest

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

	err := s.configService.AddDNSEntry(c.Request.Context(), json.Type, json.Domain, json.Value, json.Comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// GetDNSEntries returns all DNS entries
// @Summary      Get DNS entries
// @Description  Returns all DNS entries
// @Tags         dns
// @Produce      json
// @Success      200  {array}   services.DNSEntry
// @Failure      500  {object}  map[string]string
// @Router       /dns/entries [get]
func (s *Server) GetDNSEntries(c *gin.Context) {
	entries, err := s.configService.GetDNSEntries(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, entries)
}

// DeleteDNSEntry deletes a DNS entry
// @Summary      Delete DNS entry
// @Description  Deletes a DNS entry
// @Tags         dns
// @Accept       json
// @Produce      json
// @Param        entry  body      services.DNSEntry  true  "DNS Entry"
// @Success      200    {object}  map[string]string
// @Failure      400    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Router       /dns/entries [delete]
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

// UpdateDNSEntry updates a DNS entry
// @Summary      Update DNS entry
// @Description  Updates a DNS entry
// @Tags         dns
// @Accept       json
// @Produce      json
// @Param        entry  body      UpdateDNSEntryRequest  true  "DNS Entry Update"
// @Success      200    {object}  map[string]string
// @Failure      400    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Router       /dns/entries [put]
func (s *Server) UpdateDNSEntry(c *gin.Context) {
	var json UpdateDNSEntryRequest
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
