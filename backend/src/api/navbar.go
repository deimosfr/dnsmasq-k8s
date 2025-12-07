package api

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type NavbarItem struct {
	Label        string `json:"label"`
	Link         string `json:"link"`
	ActivePageID string `json:"activePageId"`
	ID           string `json:"id"`
}

// GetNavbar returns the list of navbar items based on configuration
// @Summary      Get navbar items
// @Description  Returns the dynamic list of navbar items based on enabled features
// @Tags         navbar
// @Produce      json
// @Success      200  {array}   NavbarItem
// @Router       /navbar [get]
func (s *Server) GetNavbar(c *gin.Context) {
	items := []NavbarItem{
		{
			Label:        "Home",
			Link:         "/static/index.html",
			ActivePageID: "home",
		},
		{
			Label:        "Config",
			Link:         "/static/pages/config.html",
			ActivePageID: "config",
		},
	}

	if os.Getenv("DNS_ENABLED") == "true" {
		items = append(items, NavbarItem{
			Label:        "DNS",
			Link:         "/static/pages/dns.html",
			ActivePageID: "dns",
			ID:           "nav-item-dns",
		})
	}

	if os.Getenv("DHCP_ENABLED") == "true" {
		items = append(items, NavbarItem{
			Label:        "DHCP",
			Link:         "/static/pages/dhcp.html",
			ActivePageID: "dhcp",
			ID:           "nav-item-dhcp",
		})
	}

	items = append(items, NavbarItem{
		Label:        "API",
		Link:         "/static/pages/api.html",
		ActivePageID: "api",
	})

	c.JSON(http.StatusOK, items)
}
