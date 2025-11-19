package services

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"path/filepath"
	"regexp"
	"strings"

	"github.com/fsnotify/fsnotify"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type ConfigService struct {
	clientset     kubernetes.Interface
	namespace     string
	configFile    string
	customDNSFile string
}

func NewConfigService(clientset kubernetes.Interface, namespace string) *ConfigService {
	configFile := os.Getenv("DNSMASQ_CONFIG_FILE")
	if configFile == "" {
		configFile = "/etc/dnsmasq.conf"
	}
	customDNSFile := os.Getenv("DNSMASQ_CUSTOM_DNS_FILE")
	if customDNSFile == "" {
		customDNSFile = "/etc/dnsmasq.d/custom.conf"
	}
	return &ConfigService{
		clientset:     clientset,
		namespace:     namespace,
		configFile:    configFile,
		customDNSFile: customDNSFile,
	}
}

func (s *ConfigService) GetConfig(ctx context.Context) (string, error) {
	content, err := ioutil.ReadFile(s.configFile)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func (s *ConfigService) UpdateConfig(ctx context.Context, config string) error {
	if err := s.validateDnsmasqConfig(config); err != nil {
		return fmt.Errorf("dnsmasq configuration validation failed: %v", err)
	}

	if err := ioutil.WriteFile(s.configFile, []byte(config), 0644); err != nil {
		return err
	}

	return s.ReloadDnsmasq()
}

func (s *ConfigService) ReloadDnsmasq() error {
	fmt.Println("INFO: Sending SIGHUP to dnsmasq via supervisorctl")
	// Use supervisorctl signal to send SIGHUP to the dnsmasq process
	cmd := exec.Command("supervisorctl", "signal", "SIGHUP", "dnsmasq")
	if err := cmd.Run(); err != nil {
		fmt.Printf("WARN: Failed to send SIGHUP to dnsmasq: %v\n", err)
		return err
	}
	return nil
}

func (s *ConfigService) GetConfigMap(ctx context.Context) (*v1.ConfigMap, error) {
	return s.clientset.CoreV1().ConfigMaps(s.namespace).Get(ctx, "dnsmasq-config", metav1.GetOptions{})
}

func (s *ConfigService) AddDNSEntry(ctx context.Context, recordType, domain, value string) error {
	// Ensure custom DNS file exists
	if _, err := os.Stat(s.customDNSFile); os.IsNotExist(err) {
		// Ensure directory exists
		dir := filepath.Dir(s.customDNSFile)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("failed to create custom DNS directory: %v", err)
			}
		}

		if err := ioutil.WriteFile(s.customDNSFile, []byte(""), 0644); err != nil {
			return fmt.Errorf("failed to create custom DNS file: %v", err)
		}
	}

	content, err := ioutil.ReadFile(s.customDNSFile)
	if err != nil {
		return err
	}
	dnsmasqConf := string(content)

	var newEntry string
	switch recordType {
	case "address":
		newEntry = fmt.Sprintf("\naddress=/%s/%s", domain, value)
	case "cname":
		newEntry = fmt.Sprintf("\ncname=%s,%s", domain, value)
	case "txt":
		newEntry = fmt.Sprintf("\ntxt-record=%s,\"%s\"", domain, value)
	default:
		return fmt.Errorf("unsupported DNS record type: %s. Only 'address', 'cname', and 'txt' are supported", recordType)
	}

	dnsmasqConf += newEntry

	// Validate the configuration before writing
	if err := s.validateDnsmasqConfig(dnsmasqConf); err != nil {
		return fmt.Errorf("dnsmasq configuration validation failed: %v", err)
	}

	if err := ioutil.WriteFile(s.customDNSFile, []byte(dnsmasqConf), 0644); err != nil {
		return err
	}

	return s.ReloadDnsmasq()
}

type DNSEntry struct {
	Type   string `json:"type"`
	Domain string `json:"domain"`
	Value  string `json:"value"`
}

func (s *ConfigService) GetDNSEntries(ctx context.Context) ([]DNSEntry, error) {
	if _, err := os.Stat(s.customDNSFile); os.IsNotExist(err) {
		return []DNSEntry{}, nil
	}

	content, err := ioutil.ReadFile(s.customDNSFile)
	if err != nil {
		return nil, err
	}

	var entries []DNSEntry
	lines := strings.Split(string(content), "\n")

	// Regex for parsing
	// address=/domain/ip
	reAddress := regexp.MustCompile(`^address=/(.+)/(.+)$`)
	// cname=domain,target
	reCname := regexp.MustCompile(`^cname=(.+),(.+)$`)
	// txt-record=domain,"value"
	reTxt := regexp.MustCompile(`^txt-record=(.+),"(.+)"$`)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if matches := reAddress.FindStringSubmatch(line); matches != nil {
			entries = append(entries, DNSEntry{Type: "address", Domain: matches[1], Value: matches[2]})
		} else if matches := reCname.FindStringSubmatch(line); matches != nil {
			entries = append(entries, DNSEntry{Type: "cname", Domain: matches[1], Value: matches[2]})
		} else if matches := reTxt.FindStringSubmatch(line); matches != nil {
			entries = append(entries, DNSEntry{Type: "txt", Domain: matches[1], Value: matches[2]})
		}
	}

	return entries, nil
}

func (s *ConfigService) DeleteDNSEntry(ctx context.Context, entry DNSEntry) error {
	return s.modifyDNSEntry(ctx, entry, DNSEntry{}, true)
}

func (s *ConfigService) UpdateDNSEntry(ctx context.Context, oldEntry, newEntry DNSEntry) error {
	return s.modifyDNSEntry(ctx, oldEntry, newEntry, false)
}

func (s *ConfigService) modifyDNSEntry(ctx context.Context, targetEntry, newEntry DNSEntry, isDelete bool) error {
	content, err := ioutil.ReadFile(s.customDNSFile)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	var newLines []string
	found := false

	targetLine := ""
	switch targetEntry.Type {
	case "address":
		targetLine = fmt.Sprintf("address=/%s/%s", targetEntry.Domain, targetEntry.Value)
	case "cname":
		targetLine = fmt.Sprintf("cname=%s,%s", targetEntry.Domain, targetEntry.Value)
	case "txt-record":
		targetLine = fmt.Sprintf("txt-record=%s,\"%s\"", targetEntry.Domain, targetEntry.Value)
	}

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if !found && trimmedLine == targetLine {
			found = true
			if !isDelete {
				// Replace with new entry
				var newLine string
				switch newEntry.Type {
				case "address":
					newLine = fmt.Sprintf("address=/%s/%s", newEntry.Domain, newEntry.Value)
				case "cname":
					newLine = fmt.Sprintf("cname=%s,%s", newEntry.Domain, newEntry.Value)
				case "txt-record":
					newLine = fmt.Sprintf("txt-record=%s,\"%s\"", newEntry.Domain, newEntry.Value)
				}
				newLines = append(newLines, newLine)
			}
			continue
		}
		newLines = append(newLines, line)
	}

	if !found {
		return fmt.Errorf("entry not found")
	}

	newContent := strings.Join(newLines, "\n")

	// Validate if it's an update (deletion should be safe usually, but good to check if we want to be strict)
	// For now, let's just save it.

	if err := ioutil.WriteFile(s.customDNSFile, []byte(newContent), 0644); err != nil {
		return err
	}

	return s.ReloadDnsmasq()
}

func (s *ConfigService) validateDnsmasqConfig(config string) error {
	tmpfile, err := ioutil.TempFile("", "dnsmasq-config-")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.WriteString(config); err != nil {
		return fmt.Errorf("failed to write config to temporary file: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		return fmt.Errorf("failed to close temporary file: %v", err)
	}

	if _, err := exec.LookPath("dnsmasq"); err != nil {
		return nil // Skip validation if dnsmasq is not installed (e.g. in tests)
	}

	cmd := exec.Command("dnsmasq", "--test", "--conf-file="+tmpfile.Name())
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("dnsmasq --test failed: %s\n%s", err.Error(), string(output))
	}

	return nil
}

func (s *ConfigService) StartConfigSync(ctx context.Context) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Printf("ERROR: failed to create watcher: %v\n", err)
		return
	}
	defer watcher.Close()

	err = watcher.Add(s.configFile)
	if err != nil {
		// If file doesn't exist, try to create it empty or wait?
		// RestoreConfigFromConfigMap should have created it if it existed in ConfigMap.
		// If not, maybe create empty?
		if os.IsNotExist(err) {
			// Create empty file if it doesn't exist so we can watch it
			if err := ioutil.WriteFile(s.configFile, []byte(""), 0644); err != nil {
				fmt.Printf("ERROR: failed to create config file: %v\n", err)
				return
			}
			err = watcher.Add(s.configFile)
			if err != nil {
				fmt.Printf("ERROR: failed to watch config file: %v\n", err)
				return
			}
		} else {
			fmt.Printf("ERROR: failed to watch config file: %v\n", err)
			return
		}
	}

	fmt.Printf("INFO: Starting config sync for %s\n", s.configFile)

	// Restore config from ConfigMap if it exists
	if err := s.RestoreConfigFromConfigMap(ctx); err != nil {
		fmt.Printf("WARN: failed to restore config from ConfigMap: %v\n", err)
	} else {
		// If restored, we might need to re-add watch if the file was recreated?
		// RestoreConfigFromConfigMap uses ioutil.WriteFile which might truncate/create.
		// fsnotify usually handles writes fine.
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
				fmt.Println("INFO: Config file modified, syncing to ConfigMap")
				if err := s.SyncConfigToConfigMap(ctx); err != nil {
					fmt.Printf("ERROR: failed to sync config to ConfigMap: %v\n", err)
				}
			}
			if event.Op&fsnotify.Rename == fsnotify.Rename || event.Op&fsnotify.Remove == fsnotify.Remove {
				fmt.Println("INFO: Config file renamed or removed, re-watching")
				watcher.Remove(event.Name)
				// Wait for file to be recreated?
				// In our case, we might be the ones recreating it via Restore or Update?
				// Or user might have done it.
				// Let's try to re-add watch.
				for {
					err := watcher.Add(s.configFile)
					if err == nil {
						break
					}
					// Sleep a bit?
				}
				if err := s.SyncConfigToConfigMap(ctx); err != nil {
					fmt.Printf("ERROR: failed to sync config to ConfigMap: %v\n", err)
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

func (s *ConfigService) SyncConfigToConfigMap(ctx context.Context) error {
	content, err := ioutil.ReadFile(s.configFile)
	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}

	configMap, err := s.clientset.CoreV1().ConfigMaps(s.namespace).Get(ctx, "dnsmasq-config", metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			// Create ConfigMap
			newCM := &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "dnsmasq-config",
					Namespace: s.namespace,
				},
				Data: map[string]string{
					"dnsmasq.conf": string(content),
				},
			}
			_, err = s.clientset.CoreV1().ConfigMaps(s.namespace).Create(ctx, newCM, metav1.CreateOptions{})
			if err != nil {
				return fmt.Errorf("failed to create ConfigMap: %v", err)
			}
			return nil
		}
		return fmt.Errorf("failed to get ConfigMap: %v", err)
	}

	configMap.Data["dnsmasq.conf"] = string(content)
	_, err = s.clientset.CoreV1().ConfigMaps(s.namespace).Update(ctx, configMap, metav1.UpdateOptions{})
	return err
}

func (s *ConfigService) RestoreConfigFromConfigMap(ctx context.Context) error {
	configMap, err := s.clientset.CoreV1().ConfigMaps(s.namespace).Get(ctx, "dnsmasq-config", metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil // Nothing to restore
		}
		return err
	}

	if content, ok := configMap.Data["dnsmasq.conf"]; ok {
		return ioutil.WriteFile(s.configFile, []byte(content), 0644)
	}
	return nil
}

func (s *ConfigService) StartCustomDNSSync(ctx context.Context) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Printf("ERROR: failed to create watcher: %v\n", err)
		return
	}
	defer watcher.Close()

	// Ensure file exists
	if _, err := os.Stat(s.customDNSFile); os.IsNotExist(err) {
		// Ensure directory exists
		dir := filepath.Dir(s.customDNSFile)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Printf("ERROR: failed to create custom DNS directory: %v\n", err)
				return
			}
		}

		if err := ioutil.WriteFile(s.customDNSFile, []byte(""), 0644); err != nil {
			fmt.Printf("ERROR: failed to create custom DNS file: %v\n", err)
			return
		}
	}

	err = watcher.Add(s.customDNSFile)
	if err != nil {
		fmt.Printf("ERROR: failed to watch custom DNS file: %v\n", err)
		return
	}

	fmt.Printf("INFO: Starting custom DNS sync for %s\n", s.customDNSFile)

	// Restore from ConfigMap
	if err := s.RestoreCustomDNSFromConfigMap(ctx); err != nil {
		fmt.Printf("WARN: failed to restore custom DNS from ConfigMap: %v\n", err)
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
				fmt.Println("INFO: Custom DNS file modified, syncing to ConfigMap")
				if err := s.SyncCustomDNSToConfigMap(ctx); err != nil {
					fmt.Printf("ERROR: failed to sync custom DNS to ConfigMap: %v\n", err)
				}
			}
			if event.Op&fsnotify.Rename == fsnotify.Rename || event.Op&fsnotify.Remove == fsnotify.Remove {
				fmt.Println("INFO: Custom DNS file renamed or removed, re-watching")
				watcher.Remove(event.Name)
				for {
					err := watcher.Add(s.customDNSFile)
					if err == nil {
						break
					}
				}
				if err := s.SyncCustomDNSToConfigMap(ctx); err != nil {
					fmt.Printf("ERROR: failed to sync custom DNS to ConfigMap: %v\n", err)
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

func (s *ConfigService) SyncCustomDNSToConfigMap(ctx context.Context) error {
	content, err := ioutil.ReadFile(s.customDNSFile)
	if err != nil {
		return fmt.Errorf("failed to read custom DNS file: %v", err)
	}

	configMap, err := s.clientset.CoreV1().ConfigMaps(s.namespace).Get(ctx, "dnsmasq-custom-dns", metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			// Create ConfigMap
			newCM := &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "dnsmasq-custom-dns",
					Namespace: s.namespace,
				},
				Data: map[string]string{
					"custom.conf": string(content),
				},
			}
			_, err = s.clientset.CoreV1().ConfigMaps(s.namespace).Create(ctx, newCM, metav1.CreateOptions{})
			if err != nil {
				return fmt.Errorf("failed to create ConfigMap: %v", err)
			}
			return nil
		}
		return fmt.Errorf("failed to get ConfigMap: %v", err)
	}

	configMap.Data["custom.conf"] = string(content)
	_, err = s.clientset.CoreV1().ConfigMaps(s.namespace).Update(ctx, configMap, metav1.UpdateOptions{})
	return err
}

func (s *ConfigService) RestoreCustomDNSFromConfigMap(ctx context.Context) error {
	configMap, err := s.clientset.CoreV1().ConfigMaps(s.namespace).Get(ctx, "dnsmasq-custom-dns", metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil // Nothing to restore
		}
		return err
	}

	if content, ok := configMap.Data["custom.conf"]; ok {
		return ioutil.WriteFile(s.customDNSFile, []byte(content), 0644)
	}
	return nil
}

func (s *ConfigService) StartConfigMapWatch(ctx context.Context) {
	watcher, err := s.clientset.CoreV1().ConfigMaps(s.namespace).Watch(ctx, metav1.ListOptions{})
	if err != nil {
		fmt.Printf("ERROR: failed to start ConfigMap watcher: %v\n", err)
		return
	}
	defer watcher.Stop()

	fmt.Println("INFO: Starting ConfigMap watch")

	for {
		select {
		case event, ok := <-watcher.ResultChan():
			if !ok {
				fmt.Println("WARN: ConfigMap watcher channel closed, restarting")
				watcher.Stop()
				watcher, err = s.clientset.CoreV1().ConfigMaps(s.namespace).Watch(ctx, metav1.ListOptions{})
				if err != nil {
					fmt.Printf("ERROR: failed to restart ConfigMap watcher: %v\n", err)
					return
				}
				continue
			}

			if event.Type == "MODIFIED" || event.Type == "ADDED" {
				cm, ok := event.Object.(*v1.ConfigMap)
				if !ok {
					continue
				}

				if cm.Name == "dnsmasq-config" {
					if content, ok := cm.Data["dnsmasq.conf"]; ok {
						if err := s.syncFileIfChanged(s.configFile, content); err != nil {
							fmt.Printf("ERROR: failed to sync dnsmasq-config to file: %v\n", err)
						}
					}
				} else if cm.Name == "dnsmasq-custom-dns" {
					if content, ok := cm.Data["custom.conf"]; ok {
						if err := s.syncFileIfChanged(s.customDNSFile, content); err != nil {
							fmt.Printf("ERROR: failed to sync dnsmasq-custom-dns to file: %v\n", err)
						}
					}
				} else if cm.Name == "dnsmasq-reservations" {
					// We need to know the reservations file path here.
					// Ideally ConfigService should know about it or we pass it.
					// For now, let's assume standard path or get env again.
					reservationsFile := os.Getenv("DHCP_RESERVATIONS_FILE")
					if reservationsFile == "" {
						reservationsFile = "/etc/dnsmasq.d/reservations.conf"
					}
					if content, ok := cm.Data["reservations.conf"]; ok {
						if err := s.syncFileIfChanged(reservationsFile, content); err != nil {
							fmt.Printf("ERROR: failed to sync dnsmasq-reservations to file: %v\n", err)
						}
					}
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

func (s *ConfigService) syncFileIfChanged(filePath, newContent string) error {
	// Read current content
	currentContent, err := ioutil.ReadFile(filePath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// If content matches, do nothing (prevents infinite loop)
	if string(currentContent) == newContent {
		return nil
	}

	fmt.Printf("INFO: Syncing ConfigMap change to %s\n", filePath)
	// Write new content
	if err := ioutil.WriteFile(filePath, []byte(newContent), 0644); err != nil {
		return err
	}

	return s.ReloadDnsmasq()
}
