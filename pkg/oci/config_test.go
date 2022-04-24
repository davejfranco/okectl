package oci

import (
	"os"
	"reflect"
	"testing"

	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/containerengine"
	"github.com/oracle/oci-go-sdk/core"
)

func TestUserIsEmpty(t *testing.T) {

	ouser := User{}
	if !ouser.isEmpty() {
		t.Error("Error: user previously created is empty")
	}
}

func TestUserIsNotEmpty(t *testing.T) {

	ouser := User{
		User:   "dfranco",
		Region: "us-ashburn-1",
	}

	if ouser.isEmpty() {
		t.Error("Error: user previously created is not empty")
	}
}

func TestIsValidPath(t *testing.T) {

	cfg := Config{Path: os.Getenv("HOME")}
	if !cfg.validPath() {
		t.Error("Error: Home is always a valid path ;)")
	}
}

func TestIsNotValidPath(t *testing.T) {
	cfg := Config{Path: "/whateverpath"}
	if cfg.validPath() {
		t.Error("Error: /whateverpath is no a valid path unless you like to named your path like this")
	}

}

func TestLoadWithDefault(t *testing.T) {
	cfg := Config{}
	config, err := cfg.Load()
	if err != nil {
		t.Error(err)
	}

	if reflect.TypeOf(config) != reflect.TypeOf(common.DefaultConfigProvider()) {
		t.Error("Error: should be type common.ConfigProvider")
	}
}

func TestLoadWithFileNoProfile(t *testing.T) {
	configLocation := os.Getenv("HOME") + "/.oci/config"
	commonConfig := common.CustomProfileConfigProvider(configLocation, "")

	cfg := Config{Path: configLocation}
	config, err := cfg.Load()
	if err != nil {
		t.Error(err)
	}
	if reflect.TypeOf(config) != reflect.TypeOf(commonConfig) {
		t.Errorf("Error: should be type %s", reflect.TypeOf(commonConfig).String())
	}

}

func TestLoadWithFileProfile(t *testing.T) {
	configLocation := os.Getenv("HOME") + "/.oci/config"
	commonConfig := common.CustomProfileConfigProvider(configLocation, "")

	cfg := Config{Path: configLocation, Profile: "DEFAULT"}
	config, err := cfg.Load()
	if err != nil {
		t.Error(err)
	}
	if reflect.TypeOf(config) != reflect.TypeOf(commonConfig) {
		t.Errorf("Error: should be type %s", reflect.TypeOf(commonConfig).String())
	}

}

func TestOkeClient(t *testing.T) {
	cfg_path := Config{} //config using default path and profile

	cfg, err := cfg_path.Load()
	if err != nil {
		t.Error(err)
	}
	oke, err := OkeClient(cfg)
	if err != nil {
		t.Error(err)
	}
	if reflect.TypeOf(oke) != reflect.TypeOf(containerengine.ContainerEngineClient{}) {
		t.Errorf("Error: should be type %s", reflect.TypeOf(containerengine.ContainerEngineClient{}).String())
	}
}

/*func TestVcnClient(t *testing.T) {
	cfg_path := Config{} //config using default path and profile

	cfg, err := cfg_path.Load()
	if err != nil {
		t.Error(err)
	}
	vcn, err := VcnClient(cfg)
	if err != nil {
		t.Error(err)
	}
	if reflect.TypeOf(vcn) != reflect.TypeOf(core.VirtualNetworkClient{}) {
		t.Errorf("Error: should be type %s", reflect.TypeOf(core.VirtualNetworkClient{}).String())
	}
}*/

func TestComputeClient(t *testing.T) {
	cfg_path := Config{} //config using default path and profile

	cfg, err := cfg_path.Load()
	if err != nil {
		t.Error(err)
	}
	srv, err := ComputeClient(cfg)
	if err != nil {
		t.Error(err)
	}
	if reflect.TypeOf(srv) != reflect.TypeOf(core.ComputeClient{}) {
		t.Errorf("Error: should be type %s", reflect.TypeOf(core.ComputeClient{}).String())
	}
}
