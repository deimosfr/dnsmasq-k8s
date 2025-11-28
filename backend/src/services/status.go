package services

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type StatusService struct {
	startTime time.Time
}

type Status struct {
	API                bool     `json:"api"`
	DNS                bool     `json:"dns"`
	DHCP               bool     `json:"dhcp"`
	Uptime             string   `json:"uptime"`
	DnsmasqPID         string   `json:"dnsmasq_pid"`
	SupervisorServices []string `json:"supervisor_services"`
}

func NewStatusService() *StatusService {
	return &StatusService{
		startTime: time.Now(),
	}
}

func (s *StatusService) GetStatus() *Status {
	uptime := time.Since(s.startTime)
	uptimeStr := formatDuration(uptime)

	dnsmasqPID := getDnsmasqPID()

	return &Status{
		API:                true,
		DNS:                os.Getenv("DNS_ENABLED") == "true",
		DHCP:               os.Getenv("DHCP_ENABLED") == "true",
		Uptime:             uptimeStr,
		DnsmasqPID:         dnsmasqPID,
		SupervisorServices: getSupervisorStatus(),
	}
}

func formatDuration(d time.Duration) string {
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm %ds", days, hours, minutes, seconds)
	} else if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
	} else if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}

func getDnsmasqPID() string {
	// Try pidof first
	cmd := exec.Command("pidof", "dnsmasq")
	output, err := cmd.Output()
	if err == nil && len(output) > 0 {
		// pidof may return multiple PIDs, take the first one
		pids := strings.Fields(string(output))
		if len(pids) > 0 {
			return pids[0]
		}
	}

	// Fallback: try pgrep
	cmd = exec.Command("pgrep", "dnsmasq")
	output, err = cmd.Output()
	if err == nil && len(output) > 0 {
		pids := strings.Fields(string(output))
		if len(pids) > 0 {
			return pids[0]
		}
	}

	return "N/A"
}

func getSupervisorStatus() []string {
	cmd := exec.Command("supervisorctl", "status")
	output, err := cmd.Output()

	// supervisorctl returns non-zero exit code if any service is not running (e.g. STOPPED, EXITED, FATAL)
	// We should still try to parse the output if we got any.
	if err != nil {
		// If we have no output, then it's a real error (e.g. supervisorctl not found)
		if len(output) == 0 {
			return []string{"Error fetching status"}
		}
		// If we have output, we proceed to parse it, ignoring the exit code error
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var services []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			services = append(services, line)
		}
	}
	return services
}
