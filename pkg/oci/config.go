package oci

import (
	"errors"
	"fmt"
	"os"

	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/containerengine"
	"github.com/oracle/oci-go-sdk/core"
)

//OciUser required connection arguments
type User struct {
	Tenancy        string
	User           string
	Region         string
	Fingerprint    string
	PrivateKey     string
	KeyPassphrasse string
}

func (ouser User) isEmpty() bool {
	return User{} == ouser
}

//Config client to connect to tenant
type Config struct {
	Path    string
	Profile string
}

func (cfg Config) validPath() bool {

	if _, err := os.Stat(cfg.Path); err == nil {
		return true
	}
	return false
}

//Load creates a condig object to connect to OCI
func (cfg Config) Load() (common.ConfigurationProvider, error) {

	if cfg.Path != "" || cfg.Profile != "" {
		fmt.Println(cfg.Profile)
		if !cfg.validPath() {
			return nil, errors.New("invalid Path to config file")
		}
		c := common.CustomProfileConfigProvider(
			cfg.Path,
			cfg.Profile,
		)
		return c, nil
	}

	return common.DefaultConfigProvider(), nil

}

//OkeClient provides connection to Container Engine service
func OkeClient(config common.ConfigurationProvider) (containerengine.ContainerEngineClient, error) {

	client, err := containerengine.NewContainerEngineClientWithConfigurationProvider(config)
	if err != nil {
		return containerengine.ContainerEngineClient{}, err
	}
	return client, nil
}

//ComputeClient given a config provider returns a client for computeclient
func ComputeClient(config common.ConfigurationProvider) (core.ComputeClient, error) {

	client, err := core.NewComputeClientWithConfigurationProvider(config)
	if err != nil {
		return core.ComputeClient{}, err
	}
	return client, nil
}
