package services

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fsnotify/fsnotify"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type DHCPService struct {
	clientset        kubernetes.Interface
	namespace        string
	leaseFile        string
	reservationsFile string
	configService    *ConfigService
}

type DHCPLease struct {
	MACAddress string `json:"mac_address"`
	IPAddress  string `json:"ip_address"`
	Hostname   string `json:"hostname"`
	ExpiryTime int64  `json:"expiry_time"`
}

type DHCPReservation struct {
	MACAddress string `json:"mac_address"`
	IPAddress  string `json:"ip_address"`
	Hostname   string `json:"hostname"`
}

func NewDHCPService(clientset kubernetes.Interface, namespace string, configService *ConfigService) *DHCPService {
	leaseFile := os.Getenv("DHCP_LEASE_FILE")
	if leaseFile == "" {
		leaseFile = "/var/lib/misc/dnsmasq.leases"
	}
	reservationsFile := os.Getenv("DHCP_RESERVATIONS_FILE")
	if reservationsFile == "" {
		reservationsFile = "/etc/dnsmasq.d/reservations.conf"
	}
	return &DHCPService{
		clientset:        clientset,
		namespace:        namespace,
		leaseFile:        leaseFile,
		reservationsFile: reservationsFile,
		configService:    configService,
	}
}

func (s *DHCPService) GetLeases(ctx context.Context) ([]DHCPLease, error) {
	file, err := os.Open(s.leaseFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	leases := []DHCPLease{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) >= 4 {
			expiry, _ := strconv.ParseInt(parts[0], 10, 64)
			leases = append(leases, DHCPLease{
				MACAddress: parts[1],
				IPAddress:  parts[2],
				Hostname:   parts[3],
				ExpiryTime: expiry,
			})
		}
	}

	return leases, scanner.Err()
}

func (s *DHCPService) GetReservations(ctx context.Context) ([]DHCPReservation, error) {
	if _, err := os.Stat(s.reservationsFile); os.IsNotExist(err) {
		return []DHCPReservation{}, nil
	}

	content, err := os.ReadFile(s.reservationsFile)
	if err != nil {
		return nil, err
	}

	reservations := []DHCPReservation{}
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// dhcp-host=mac,ip,hostname OR dhcp-host=hostname,mac,ip
		if strings.HasPrefix(line, "dhcp-host=") {
			parts := strings.Split(strings.TrimPrefix(line, "dhcp-host="), ",")
			if len(parts) >= 3 {
				res := DHCPReservation{}
				for _, part := range parts {
					part = strings.TrimSpace(part)
					if strings.Contains(part, ":") {
						res.MACAddress = part
					} else if strings.Contains(part, ".") {
						res.IPAddress = part
					} else {
						res.Hostname = part
					}
				}
				if res.MACAddress != "" && res.IPAddress != "" {
					reservations = append(reservations, res)
				}
			}
		}
	}
	return reservations, nil
}

func (s *DHCPService) AddReservation(ctx context.Context, macAddress, ipAddress, hostname string) error {
	// Ensure directory exists
	dir := filepath.Dir(s.reservationsFile)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create reservations directory: %v", err)
		}
	}

	// Append to file
	f, err := os.OpenFile(s.reservationsFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	// defer f.Close() // We close explicitly

	// Write in format: dhcp-host=hostname,mac,ip
	if _, err := f.WriteString(fmt.Sprintf("\ndhcp-host=%s,%s,%s", hostname, macAddress, ipAddress)); err != nil {
		return err
	}

	// We need to close the file before reloading, but defer handles it at end of function.
	// However, ReloadDnsmasq is external command, so it should be fine if file is flushed.
	// But to be safe, maybe we should close it explicitly or just rely on OS.
	// Actually, defer runs after return, so we should probably close it before calling ReloadDnsmasq if we want to be sure.
	// But wait, defer runs when function returns. So if we call ReloadDnsmasq here, the file is still open.
	// Let's close it first.
	f.Close()

	return s.configService.ReloadDnsmasq()
}

func (s *DHCPService) UpdateReservation(ctx context.Context, oldRes, newRes DHCPReservation) error {
	return s.modifyReservation(ctx, oldRes, newRes, false)
}

func (s *DHCPService) DeleteReservation(ctx context.Context, res DHCPReservation) error {
	return s.modifyReservation(ctx, res, DHCPReservation{}, true)
}

func (s *DHCPService) modifyReservation(ctx context.Context, target, newRes DHCPReservation, isDelete bool) error {
	content, err := os.ReadFile(s.reservationsFile)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	var newLines []string
	found := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if !found && strings.HasPrefix(trimmed, "dhcp-host=") {
			parts := strings.Split(strings.TrimPrefix(trimmed, "dhcp-host="), ",")
			if len(parts) >= 3 {
				currentRes := DHCPReservation{}
				for _, part := range parts {
					part = strings.TrimSpace(part)
					if strings.Contains(part, ":") {
						currentRes.MACAddress = part
					} else if strings.Contains(part, ".") {
						currentRes.IPAddress = part
					} else {
						currentRes.Hostname = part
					}
				}

				if currentRes.MACAddress == target.MACAddress && currentRes.IPAddress == target.IPAddress && currentRes.Hostname == target.Hostname {
					found = true
					if !isDelete {
						// Write in format: dhcp-host=hostname,mac,ip
						newLines = append(newLines, fmt.Sprintf("dhcp-host=%s,%s,%s", newRes.Hostname, newRes.MACAddress, newRes.IPAddress))
					}
					continue
				}
			}
		}
		newLines = append(newLines, line)
	}

	if !found {
		return fmt.Errorf("reservation not found")
	}

	if err := os.WriteFile(s.reservationsFile, []byte(strings.Join(newLines, "\n")), 0644); err != nil {
		return err
	}

	return s.configService.ReloadDnsmasq()
}

func (s *DHCPService) UpdateLease(ctx context.Context, oldLease, newLease DHCPLease) error {
	// Leases file format: timestamp mac ip hostname clientid
	// We only care about mac, ip, hostname for now. Timestamp is usually first.
	// This is tricky because dnsmasq manages this file.
	// Editing it while dnsmasq is running might be overwritten.
	// However, for the sake of this UI, we will try to edit it.
	// Ideally, we should use dhcp_release or similar tools, but we are editing the file directly.

	content, err := os.ReadFile(s.leaseFile)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	var newLines []string
	found := false

	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) >= 4 {
			// parts[1] is mac, parts[2] is ip, parts[3] is hostname
			if parts[1] == oldLease.MACAddress && parts[2] == oldLease.IPAddress && parts[3] == oldLease.Hostname {
				found = true
				// Preserve timestamp (parts[0]) and clientid (parts[4] if exists)
				newLine := fmt.Sprintf("%s %s %s %s", parts[0], newLease.MACAddress, newLease.IPAddress, newLease.Hostname)
				if len(parts) > 4 {
					newLine += " " + strings.Join(parts[4:], " ")
				}
				newLines = append(newLines, newLine)
				continue
			}
		}
		newLines = append(newLines, line)
	}

	if !found {
		return fmt.Errorf("lease not found")
	}

	if err := os.WriteFile(s.leaseFile, []byte(strings.Join(newLines, "\n")), 0644); err != nil {
		return err
	}

	return s.configService.ReloadDnsmasq()
}

func (s *DHCPService) DeleteLease(ctx context.Context, lease DHCPLease) error {
	content, err := os.ReadFile(s.leaseFile)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	var newLines []string
	found := false

	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) >= 4 {
			if parts[1] == lease.MACAddress && parts[2] == lease.IPAddress && parts[3] == lease.Hostname {
				found = true
				continue // Skip adding this line
			}
		}
		newLines = append(newLines, line)
	}

	if !found {
		return fmt.Errorf("lease not found")
	}

	if err := os.WriteFile(s.leaseFile, []byte(strings.Join(newLines, "\n")), 0644); err != nil {
		return err
	}

	return s.configService.ReloadDnsmasq()
}

func (s *DHCPService) StartLeaseSync(ctx context.Context) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Printf("ERROR: failed to create watcher: %v\n", err)
		return
	}
	defer watcher.Close()

	// Ensure lease file exists
	if _, err := os.Stat(s.leaseFile); os.IsNotExist(err) {
		file, err := os.Create(s.leaseFile)
		if err != nil {
			fmt.Printf("ERROR: failed to create lease file: %v\n", err)
			return
		}
		file.Close()
	}

	err = watcher.Add(s.leaseFile)
	if err != nil {
		fmt.Printf("ERROR: failed to watch lease file: %v\n", err)
		return
	}

	fmt.Printf("INFO: Starting lease sync for %s\n", s.leaseFile)

	// Restore leases from ConfigMap if they exist
	if err := s.RestoreLeasesFromConfigMap(ctx); err != nil {
		fmt.Printf("WARN: failed to restore leases from ConfigMap: %v\n", err)
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
				fmt.Println("INFO: lease file modified/created, syncing to ConfigMap")
				if err := s.SyncLeasesToConfigMap(ctx); err != nil {
					fmt.Printf("ERROR: failed to sync leases: %v\n", err)
				}
			}
			// Handle atomic updates (Rename/Remove)
			if event.Op&fsnotify.Rename == fsnotify.Rename || event.Op&fsnotify.Remove == fsnotify.Remove {
				fmt.Println("INFO: lease file renamed/removed, re-watching")
				// Wait a bit for the file to be recreated
				// In a real scenario, we might want a retry loop here
				watcher.Remove(s.leaseFile)
				watcher.Add(s.leaseFile)
				// Trigger sync just in case
				if err := s.SyncLeasesToConfigMap(ctx); err != nil {
					fmt.Printf("ERROR: failed to sync leases after recreate: %v\n", err)
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			fmt.Printf("ERROR: watcher error: %v\n", err)
		case <-ctx.Done():
			return
		}
	}
}

func (s *DHCPService) SyncLeasesToConfigMap(ctx context.Context) error {
	content, err := os.ReadFile(s.leaseFile)
	if err != nil {
		return err
	}

	configMap, err := s.clientset.CoreV1().ConfigMaps(s.namespace).Get(ctx, "dnsmasq-leases", metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			fmt.Println("INFO: ConfigMap not found, creating it")
			// Create if not exists
			newCM := &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "dnsmasq-leases",
					Namespace: s.namespace,
				},
				Data: map[string]string{
					"dnsmasq.leases": string(content),
				},
			}
			_, err = s.clientset.CoreV1().ConfigMaps(s.namespace).Create(ctx, newCM, metav1.CreateOptions{})
			if err != nil {
				fmt.Printf("ERROR: failed to create ConfigMap: %v\n", err)
			}
			return err
		}
		fmt.Printf("ERROR: failed to get ConfigMap: %v\n", err)
		return err
	}

	configMap.Data["dnsmasq.leases"] = string(content)

	_, err = s.clientset.CoreV1().ConfigMaps(s.namespace).Update(ctx, configMap, metav1.UpdateOptions{})
	return err
}

func (s *DHCPService) RestoreLeasesFromConfigMap(ctx context.Context) error {
	configMap, err := s.clientset.CoreV1().ConfigMaps(s.namespace).Get(ctx, "dnsmasq-leases", metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil // Nothing to restore
		}
		return err
	}

	if content, ok := configMap.Data["dnsmasq.leases"]; ok {
		if err := os.WriteFile(s.leaseFile, []byte(content), 0644); err != nil {
			return err
		}
		return s.configService.ReloadDnsmasq()
	}
	return nil
}

func (s *DHCPService) StartReservationsSync(ctx context.Context) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Printf("ERROR: failed to create watcher: %v\n", err)
		return
	}
	defer watcher.Close()

	// Ensure file exists
	if _, err := os.Stat(s.reservationsFile); os.IsNotExist(err) {
		dir := filepath.Dir(s.reservationsFile)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			os.MkdirAll(dir, 0755)
		}
		os.WriteFile(s.reservationsFile, []byte(""), 0644)
	}

	err = watcher.Add(s.reservationsFile)
	if err != nil {
		fmt.Printf("ERROR: failed to watch reservations file: %v\n", err)
		return
	}

	fmt.Printf("INFO: Starting reservations sync for %s\n", s.reservationsFile)

	if err := s.RestoreReservationsFromConfigMap(ctx); err != nil {
		fmt.Printf("WARN: failed to restore reservations from ConfigMap: %v\n", err)
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
				fmt.Println("INFO: Reservations file modified, syncing to ConfigMap")
				if err := s.SyncReservationsToConfigMap(ctx); err != nil {
					fmt.Printf("ERROR: failed to sync reservations to ConfigMap: %v\n", err)
				}
			}
			if event.Op&fsnotify.Rename == fsnotify.Rename || event.Op&fsnotify.Remove == fsnotify.Remove {
				fmt.Println("INFO: Reservations file renamed/removed, re-watching")
				watcher.Remove(event.Name)
				for {
					err := watcher.Add(s.reservationsFile)
					if err == nil {
						break
					}
				}
				if err := s.SyncReservationsToConfigMap(ctx); err != nil {
					fmt.Printf("ERROR: failed to sync reservations to ConfigMap: %v\n", err)
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			fmt.Printf("ERROR: watcher error: %v\n", err)
		case <-ctx.Done():
			return
		}
	}
}

func (s *DHCPService) SyncReservationsToConfigMap(ctx context.Context) error {
	content, err := os.ReadFile(s.reservationsFile)
	if err != nil {
		return fmt.Errorf("failed to read reservations file: %v", err)
	}

	configMap, err := s.clientset.CoreV1().ConfigMaps(s.namespace).Get(ctx, "dnsmasq-reservations", metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			newCM := &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "dnsmasq-reservations",
					Namespace: s.namespace,
				},
				Data: map[string]string{
					"reservations.conf": string(content),
				},
			}
			_, err = s.clientset.CoreV1().ConfigMaps(s.namespace).Create(ctx, newCM, metav1.CreateOptions{})
			return err
		}
		return err
	}

	configMap.Data["reservations.conf"] = string(content)
	_, err = s.clientset.CoreV1().ConfigMaps(s.namespace).Update(ctx, configMap, metav1.UpdateOptions{})
	return err
}

func (s *DHCPService) RestoreReservationsFromConfigMap(ctx context.Context) error {
	configMap, err := s.clientset.CoreV1().ConfigMaps(s.namespace).Get(ctx, "dnsmasq-reservations", metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	if content, ok := configMap.Data["reservations.conf"]; ok {
		if err := os.WriteFile(s.reservationsFile, []byte(content), 0644); err != nil {
			return err
		}
		return s.configService.ReloadDnsmasq()
	}
	return nil
}
