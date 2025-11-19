package api

import (
	"backend/src/services"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes/fake"
)

func TestGetConfig(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "dnsmasq.conf")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	_, err = tmpFile.WriteString("domain-needed\nbogus-priv")
	assert.NoError(t, err)

	os.Setenv("DNSMASQ_CONFIG_FILE", tmpFile.Name())
	defer os.Unsetenv("DNSMASQ_CONFIG_FILE")

	clientset := fake.NewSimpleClientset()
	configService := services.NewConfigService(clientset, "default")
	dhcpService := services.NewDHCPService(clientset, "default", configService)
	statusService := services.NewStatusService()
	supervisorService := services.NewSupervisorService()
	server := NewServer(configService, dhcpService, statusService, supervisorService)

	r := gin.Default()
	r.GET("/config", server.GetConfig)

	req, _ := http.NewRequest("GET", "/config", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateConfig(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "dnsmasq.conf")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	os.Setenv("DNSMASQ_CONFIG_FILE", tmpFile.Name())
	defer os.Unsetenv("DNSMASQ_CONFIG_FILE")

	clientset := fake.NewSimpleClientset()
	configService := services.NewConfigService(clientset, "default")
	dhcpService := services.NewDHCPService(clientset, "default", configService)
	statusService := services.NewStatusService()
	supervisorService := services.NewSupervisorService()
	server := NewServer(configService, dhcpService, statusService, supervisorService)

	r := gin.Default()
	r.PUT("/config", server.UpdateConfig)

	req, _ := http.NewRequest("PUT", "/config", strings.NewReader(`{"config": "new-config"}`))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetLeases(t *testing.T) {
	leaseFile, err := ioutil.TempFile("", "leases")
	assert.NoError(t, err)
	defer os.Remove(leaseFile.Name())

	os.Setenv("DHCP_LEASE_FILE", leaseFile.Name())
	defer os.Unsetenv("DHCP_LEASE_FILE")

	tmpConfigFile, err := ioutil.TempFile("", "dnsmasq.conf")
	assert.NoError(t, err)
	defer os.Remove(tmpConfigFile.Name())
	defer tmpConfigFile.Close()

	os.Setenv("DNSMASQ_CONFIG_FILE", tmpConfigFile.Name())
	defer os.Unsetenv("DNSMASQ_CONFIG_FILE")

	clientset := fake.NewSimpleClientset()
	configService := services.NewConfigService(clientset, "default")
	dhcpService := services.NewDHCPService(clientset, "default", configService)
	statusService := services.NewStatusService()
	supervisorService := services.NewSupervisorService()
	server := NewServer(configService, dhcpService, statusService, supervisorService)

	r := gin.Default()
	r.GET("/dhcp/leases", server.GetLeases)

	req, _ := http.NewRequest("GET", "/dhcp/leases", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
