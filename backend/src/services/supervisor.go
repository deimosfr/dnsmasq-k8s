package services

import (
	"fmt"
	"os/exec"
)

type SupervisorService struct{}

func NewSupervisorService() *SupervisorService {
	return &SupervisorService{}
}

func (s *SupervisorService) StartService(serviceName string) error {
	fmt.Printf("INFO: Starting supervisor service: %s\n", serviceName)
	cmd := exec.Command("supervisorctl", "start", serviceName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("ERROR: Failed to start service %s: %v, output: %s\n", serviceName, err, string(output))
		return fmt.Errorf("failed to start service: %v", err)
	}
	fmt.Printf("INFO: Service %s started successfully\n", serviceName)
	return nil
}

func (s *SupervisorService) StopService(serviceName string) error {
	fmt.Printf("INFO: Stopping supervisor service: %s\n", serviceName)
	cmd := exec.Command("supervisorctl", "stop", serviceName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("ERROR: Failed to stop service %s: %v, output: %s\n", serviceName, err, string(output))
		return fmt.Errorf("failed to stop service: %v", err)
	}
	fmt.Printf("INFO: Service %s stopped successfully\n", serviceName)
	return nil
}

func (s *SupervisorService) RestartService(serviceName string) error {
	fmt.Printf("INFO: Restarting supervisor service: %s\n", serviceName)
	cmd := exec.Command("supervisorctl", "restart", serviceName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("ERROR: Failed to restart service %s: %v, output: %s\n", serviceName, err, string(output))
		return fmt.Errorf("failed to restart service: %v", err)
	}
	fmt.Printf("INFO: Service %s restarted successfully\n", serviceName)
	return nil
}
