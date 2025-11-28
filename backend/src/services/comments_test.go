package services

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes/fake"
)

func TestDHCPReservationComments(t *testing.T) {
	// Setup temporary directory
	tmpDir, err := os.MkdirTemp("", "dhcp-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	reservationsFile := filepath.Join(tmpDir, "reservations.conf")
	os.Setenv("DHCP_RESERVATIONS_FILE", reservationsFile)
	defer os.Unsetenv("DHCP_RESERVATIONS_FILE")

	// Create dummy file with comments
	content := `
dhcp-host=AA:BB:CC:DD:EE:FF,192.168.1.10,host1 # This is a comment
dhcp-host=host2,11:22:33:44:55:66,192.168.1.11 # Another comment
dhcp-host=host3,11:22:33:44:55:77,192.168.1.12 # No mac
dhcp-host=AA:BB:CC:DD:EE:00,192.168.1.13,host4
`
	err = os.WriteFile(reservationsFile, []byte(content), 0644)
	if err != nil {
		t.Fatal(err)
	}

	clientset := fake.NewSimpleClientset()
	configService := NewConfigService(clientset, "default")
	service := NewDHCPService(clientset, "default", configService)

	// Test GetReservations
	reservations, err := service.GetReservations(context.Background())
	assert.NoError(t, err)
	assert.Len(t, reservations, 4)

	assert.Equal(t, "This is a comment", reservations[0].Comment)
	assert.Equal(t, "Another comment", reservations[1].Comment)
	assert.Equal(t, "No mac", reservations[2].Comment)
	assert.Equal(t, "", reservations[3].Comment)

	// Test AddReservation with comment
	err = service.AddReservation(context.Background(), "00:11:22:33:44:55", "192.168.1.20", "newhost", "New comment")
	assert.NoError(t, err)

	reservations, err = service.GetReservations(context.Background())
	assert.NoError(t, err)
	assert.Len(t, reservations, 5)
	assert.Equal(t, "New comment", reservations[4].Comment)

	// Verify file content
	newContent, err := os.ReadFile(reservationsFile)
	assert.NoError(t, err)
	assert.Contains(t, string(newContent), "dhcp-host=newhost,00:11:22:33:44:55,192.168.1.20 # New comment")
}

func TestDNSEntryComments(t *testing.T) {
	// Setup temporary directory
	tmpDir, err := os.MkdirTemp("", "dns-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	customDNSFile := filepath.Join(tmpDir, "custom.conf")
	os.Setenv("DNSMASQ_CUSTOM_DNS_FILE", customDNSFile)
	defer os.Unsetenv("DNSMASQ_CUSTOM_DNS_FILE")

	// Create dummy file with comments
	content := `
address=/domain1.com/1.2.3.4 # Comment 1
cname=domain2.com,target.com # Comment 2
txt-record=domain3.com,"some text" # Comment 3
address=/domain4.com/5.6.7.8
`
	err = os.WriteFile(customDNSFile, []byte(content), 0644)
	if err != nil {
		t.Fatal(err)
	}

	clientset := fake.NewSimpleClientset()
	service := NewConfigService(clientset, "default")

	// Test GetDNSEntries
	entries, err := service.GetDNSEntries(context.Background())
	assert.NoError(t, err)
	assert.Len(t, entries, 4)

	assert.Equal(t, "Comment 1", entries[0].Comment)
	assert.Equal(t, "Comment 2", entries[1].Comment)
	assert.Equal(t, "Comment 3", entries[2].Comment)
	assert.Equal(t, "", entries[3].Comment)

	// Test AddDNSEntry with comment
	err = service.AddDNSEntry(context.Background(), "address", "new.com", "9.9.9.9", "New DNS comment")
	assert.NoError(t, err)

	entries, err = service.GetDNSEntries(context.Background())
	assert.NoError(t, err)
	assert.Len(t, entries, 5)
	assert.Equal(t, "New DNS comment", entries[4].Comment)

	// Verify file content
	newContent, err := os.ReadFile(customDNSFile)
	assert.NoError(t, err)
	assert.Contains(t, string(newContent), "address=/new.com/9.9.9.9 # New DNS comment")
}
