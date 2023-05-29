package oci

import (
	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/containerengine"
	"github.com/oracle/oci-go-sdk/resourcemanager"
)

// Config
type Config struct {
	Profile        string
	ConfigLocation string
}

func (c Config) configProvider() common.ConfigurationProvider {

	if c.Profile != "" {
		return common.CustomProfileConfigProvider(c.ConfigLocation, c.Profile)
	}
	return common.DefaultConfigProvider()
}

// Resource Manager Client
func (c Config) ResourceManager() (resourcemanager.ResourceManagerClient, error) {
	client, err := resourcemanager.NewResourceManagerClientWithConfigurationProvider(c.configProvider())
	if err != nil {
		return resourcemanager.ResourceManagerClient{}, err
	}
	return client, nil
}

// Container Engine
// create client using config in default location
func (c Config) Oke() (containerengine.ContainerEngineClient, error) {
	client, err := containerengine.NewContainerEngineClientWithConfigurationProvider(c.configProvider())
	if err != nil {
		return containerengine.ContainerEngineClient{}, err
	}
	return client, nil
}
