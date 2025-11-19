package api

import (
	"backend/src/services"
	"context"
)

type Server struct {
	configService     *services.ConfigService
	dhcpService       *services.DHCPService
	statusService     *services.StatusService
	supervisorService *services.SupervisorService
}

func NewServer(configService *services.ConfigService, dhcpService *services.DHCPService, statusService *services.StatusService, supervisorService *services.SupervisorService) *Server {
	server := &Server{
		configService:     configService,
		dhcpService:       dhcpService,
		statusService:     statusService,
		supervisorService: supervisorService,
	}

	go dhcpService.StartLeaseSync(context.Background())
	go configService.StartConfigSync(context.Background())
	go configService.StartCustomDNSSync(context.Background())
	go configService.StartConfigMapWatch(context.Background())
	go dhcpService.StartReservationsSync(context.Background())

	return server
}
