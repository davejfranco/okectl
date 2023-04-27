package oci

import (
	"reflect"
	"testing"

	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/containerengine"
)

// test configProvider
func TestConfigProvider(t *testing.T) {
	cfg := Config{
		Profile:        "DEFAULT",
		ConfigLocation: "",
	}
	config := cfg.configProvider()
	if reflect.TypeOf(config) != reflect.TypeOf(common.DefaultConfigProvider()) {
		t.Errorf("configProvider function failed")
	}
}

// test OKE
func TestOKE(t *testing.T) {
	cfg := Config{
		Profile:        "DEFAULT",
		ConfigLocation: "",
	}
	client, err := cfg.OKE()
	if err != nil {
		t.Errorf("OKE function failed")
	}
	if reflect.TypeOf(client) != reflect.TypeOf(containerengine.ContainerEngineClient{}) {
		t.Errorf("OKE function failed")
	}
}
