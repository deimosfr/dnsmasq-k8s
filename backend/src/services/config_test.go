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

func TestConfigService_GetConfig(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "dnsmasq.conf")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	_, err = tmpFile.WriteString("domain-needed\nbogus-priv")
	assert.NoError(t, err)

	os.Setenv("DNSMASQ_CONFIG_FILE", tmpFile.Name())
	defer os.Unsetenv("DNSMASQ_CONFIG_FILE")

	clientset := fake.NewSimpleClientset()
	configService := NewConfigService(clientset, "default")

	config, err := configService.GetConfig(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, "domain-needed\nbogus-priv", config)
}

func TestConfigService_UpdateConfig(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "dnsmasq.conf")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	os.Setenv("DNSMASQ_CONFIG_FILE", tmpFile.Name())
	defer os.Unsetenv("DNSMASQ_CONFIG_FILE")

	clientset := fake.NewSimpleClientset()
	configService := NewConfigService(clientset, "default")

	err = configService.UpdateConfig(context.Background(), "new-config")
	assert.NoError(t, err)

	content, err := ioutil.ReadFile(tmpFile.Name())
	assert.NoError(t, err)
	assert.Equal(t, "new-config", string(content))
}

func TestConfigService_SyncConfigToConfigMap(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "dnsmasq.conf")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	_, err = tmpFile.WriteString("synced-config")
	assert.NoError(t, err)

	os.Setenv("DNSMASQ_CONFIG_FILE", tmpFile.Name())
	defer os.Unsetenv("DNSMASQ_CONFIG_FILE")

	clientset := fake.NewSimpleClientset()
	configService := NewConfigService(clientset, "default")

	err = configService.SyncConfigToConfigMap(context.Background())
	assert.NoError(t, err)

	configMap, err := clientset.CoreV1().ConfigMaps("default").Get(context.Background(), "dnsmasq-config", metav1.GetOptions{})
	assert.NoError(t, err)
	assert.Equal(t, "synced-config", configMap.Data["dnsmasq.conf"])
}

func TestConfigService_RestoreConfigFromConfigMap(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "dnsmasq.conf")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	os.Setenv("DNSMASQ_CONFIG_FILE", tmpFile.Name())
	defer os.Unsetenv("DNSMASQ_CONFIG_FILE")

	clientset := fake.NewSimpleClientset(&v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "dnsmasq-config",
			Namespace: "default",
		},
		Data: map[string]string{
			"dnsmasq.conf": "restored-config",
		},
	})
	configService := NewConfigService(clientset, "default")

	err = configService.RestoreConfigFromConfigMap(context.Background())
	assert.NoError(t, err)

	content, err := ioutil.ReadFile(tmpFile.Name())
	assert.NoError(t, err)
	assert.Equal(t, "restored-config", string(content))
}
