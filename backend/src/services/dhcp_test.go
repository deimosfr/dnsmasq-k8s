package services

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestDHCPService_GetLeases(t *testing.T) {
	leaseFile, err := ioutil.TempFile("", "leases")
	assert.NoError(t, err)
	defer os.Remove(leaseFile.Name())

	_, err = leaseFile.WriteString("1677721600 00:0c:29:1c:bf:3b 192.168.1.100 my-host *\n")
	assert.NoError(t, err)

	os.Setenv("DHCP_LEASE_FILE", leaseFile.Name())
	defer os.Unsetenv("DHCP_LEASE_FILE")

	clientset := fake.NewSimpleClientset()
	dhcpService := NewDHCPService(clientset, "default", nil)

	leases, err := dhcpService.GetLeases(context.Background())
	assert.NoError(t, err)
	assert.Len(t, leases, 1)
	assert.Equal(t, "00:0c:29:1c:bf:3b", leases[0].MACAddress)
	assert.Equal(t, "192.168.1.100", leases[0].IPAddress)
	assert.Equal(t, "my-host", leases[0].Hostname)
}

func TestDHCPService_GetLeases_Empty(t *testing.T) {
	leaseFile, err := ioutil.TempFile("", "leases")
	assert.NoError(t, err)
	defer os.Remove(leaseFile.Name())

	os.Setenv("DHCP_LEASE_FILE", leaseFile.Name())
	defer os.Unsetenv("DHCP_LEASE_FILE")

	clientset := fake.NewSimpleClientset()
	dhcpService := NewDHCPService(clientset, "default", nil)

	leases, err := dhcpService.GetLeases(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, leases)
	assert.Len(t, leases, 0)
}

func TestDHCPService_SyncLeasesToConfigMap(t *testing.T) {
	leaseFile, err := ioutil.TempFile("", "leases")
	assert.NoError(t, err)
	defer os.Remove(leaseFile.Name())

	_, err = leaseFile.WriteString("1677721600 00:0c:29:1c:bf:3b 192.168.1.100 my-host *\n")
	assert.NoError(t, err)

	clientset := fake.NewSimpleClientset(&v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "dnsmasq-leases",
			Namespace: "default",
		},
		Data: map[string]string{
			"dnsmasq.leases": "",
		},
	})

	dhcpService := NewDHCPService(clientset, "default", nil)
	// Manually set lease file since we can't easily mock env var in parallel tests safely without mutex,
	// but here we are just testing the method which uses the struct field.
	// Actually NewDHCPService reads env var. Let's just set the field directly if possible or use env var.
	// Since we are in a test function, setting env var is "okay" if not running parallel.
	os.Setenv("DHCP_LEASE_FILE", leaseFile.Name())
	defer os.Unsetenv("DHCP_LEASE_FILE")

	// Re-create service to pick up env var
	dhcpService = NewDHCPService(clientset, "default", nil)

	err = dhcpService.SyncLeasesToConfigMap(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify ConfigMap content
	cm, err := clientset.CoreV1().ConfigMaps("default").Get(context.Background(), "dnsmasq-leases", metav1.GetOptions{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedContent := "1677721600 00:0c:29:1c:bf:3b 192.168.1.100 my-host *\n"
	if cm.Data["dnsmasq.leases"] != expectedContent {
		t.Errorf("expected %q, got %q", expectedContent, cm.Data["dnsmasq.leases"])
	}
}

func TestDHCPService_RestoreLeasesFromConfigMap(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	configService := NewConfigService(clientset, "default")

	// Create a ConfigMap with leases
	cm := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "dnsmasq-leases",
			Namespace: "default",
		},
		Data: map[string]string{
			"dnsmasq.leases": "restored-lease-content",
		},
	}
	_, err := clientset.CoreV1().ConfigMaps("default").Create(context.Background(), cm, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("failed to create configmap: %v", err)
	}

	// Create a temp file for leases
	tmpFile, err := ioutil.TempFile("", "dnsmasq.leases")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	os.Setenv("DHCP_LEASE_FILE", tmpFile.Name())
	defer os.Unsetenv("DHCP_LEASE_FILE")

	dhcpService := NewDHCPService(clientset, "default", configService)

	err = dhcpService.RestoreLeasesFromConfigMap(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify file content
	content, err := ioutil.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	if string(content) != "restored-lease-content" {
		t.Errorf("expected restored-lease-content, got %s", string(content))
	}
}
