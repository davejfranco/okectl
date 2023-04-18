package oci

import (
	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/containerengine"
	"github.com/oracle/oci-go-sdk/resourcemanager"
)

type Config struct {
	Profile        string
	ConfigLocation string
}

func (c *Config) configProvider() common.ConfigurationProvider {

	if c.Profile != "" {
		return common.CustomProfileConfigProvider(c.Profile, "")
	}
	return common.DefaultConfigProvider()

}

// create client using config in default location
func (c *Config) OKE() (containerengine.ContainerEngineClient, error) {
	client, err := containerengine.NewContainerEngineClientWithConfigurationProvider(c.configProvider())
	if err != nil {
		return containerengine.ContainerEngineClient{}, err
	}
	return client, nil
}

func (c *Config) ResourceManager() (resourcemanager.ResourceManagerClient, error) {
	client, err := resourcemanager.NewResourceManagerClientWithConfigurationProvider(c.configProvider())
	if err != nil {
		return resourcemanager.ResourceManagerClient{}, err
	}
	return client, nil
}
