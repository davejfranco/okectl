package util

import (
	"fmt"
	"os"

	"github.com/oracle/oci-go-sdk/common"
	ce "github.com/oracle/oci-go-sdk/containerengine"
	"github.com/oracle/oci-go-sdk/core"
)

const (
	defaultConfigPath = "~/.oci/config"
)

//OciUser required connection arguments
type OciUser struct {
	Tenancy        string
	User           string
	Region         string
	Fingerprint    string
	privateKey     string
	KeyPassphrasse string
}

func (ouser OciUser) isEmpty() bool {
	if (OciUser{}) == ouser {
		return true
	}
	return false
}

//Config client to connect to tenant
type Config struct {
	Path    string
	Profile string
	OciUser
}

func (cfg Config) validPath() bool {

	if _, err := os.Stat(cfg.Path); err == nil {
		return true
	}
	return false
}

//WithFile you are providing a config file to connect to OCI tenant
func (cfg Config) WithFile() (common.ConfigurationProvider, error) {

	if cfg.validPath() {
		if cfg.Profile != "" {
			c, err := common.ConfigurationProviderFromFileWithProfile(
				cfg.Path,
				cfg.Profile,
				cfg.KeyPassphrasse,
			)
			if err != nil {
				return nil, err
			}
			return c, nil
		}
		c, err := common.ConfigurationProviderFromFile(cfg.Path, cfg.KeyPassphrasse)
		if err != nil {
			return nil, err
		}
		return c, nil
	}
	return nil, fmt.Errorf("no a valid Path")
}

//WithDefault return configuration provider with default config file
func (cfg Config) WithDefault() (common.ConfigurationProvider, error) {

	cfg.Path = defaultConfigPath
	if cfg.validPath() {
		return common.DefaultConfigProvider(), nil
	}
	return nil, fmt.Errorf("default config file was not found")
}

//VcnClient returns a client using default config on ~/.oci/config
func VcnClient() (core.VirtualNetworkClient, error) {

	//create client
	//client, err := ce.NewContainerEngineClientWithConfigurationProvider(common.DefaultConfigProvider())
	client, err := core.NewVirtualNetworkClientWithConfigurationProvider(common.DefaultConfigProvider())
	if err != nil {
		return core.VirtualNetworkClient{}, err
	}

	return client, nil
}

//OkeClient returns a client using default config on ~/.oci/config
func OkeClient() (ce.ContainerEngineClient, error) {

	c := Config{Path: "/Users/dfranco/.oci/config"}
	config, err := c.WithFile()
	if err != nil {
		panic(err)
	}

	//create client
	client, err := ce.NewContainerEngineClientWithConfigurationProvider(config)
	if err != nil {
		return ce.ContainerEngineClient{}, err
	}

	return client, nil
}
