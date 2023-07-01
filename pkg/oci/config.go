package oci

import (
	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/containerengine"
	"github.com/oracle/oci-go-sdk/resourcemanager"
)

const (
	Profile        string = "DEFAULT"
	ConfigLocation string = "~/.oci/config"
)

// Config
type Config struct {
	Profile        string
	ConfigLocation string
	provider       common.ConfigurationProvider
}

func NewConfigProvider(profile, configLocation string) *Config {
	if profile == "" {
		profile = Profile
	}

	if configLocation == "" {
		configLocation = ConfigLocation
	}

	return &Config{
		Profile:        profile,
		ConfigLocation: configLocation,
		provider:       common.CustomProfileConfigProvider(configLocation, profile),
	}
}

func (c *Config) ConfigDetails() map[string]string {

	tenancyOCID, err := c.provider.TenancyOCID()
	if err != nil {
		panic(err)
	}

	userOCID, err := c.provider.UserOCID()
	if err != nil {
		panic(err)
	}

	region, err := c.provider.Region()
	if err != nil {
		panic(err)
	}

	provider := make(map[string]string)
	provider["tenancy"] = tenancyOCID
	provider["user"] = userOCID
	provider["region"] = region

	return provider
}

// Resource Manager Client
func (c *Config) ResourceManager() (resourcemanager.ResourceManagerClient, error) {
	client, err := resourcemanager.NewResourceManagerClientWithConfigurationProvider(c.provider)
	if err != nil {
		return resourcemanager.ResourceManagerClient{}, err
	}
	return client, nil
}

// Container Engine
// create client using config in default location
func (c *Config) Oke() (containerengine.ContainerEngineClient, error) {
	client, err := containerengine.NewContainerEngineClientWithConfigurationProvider(c.provider)
	if err != nil {
		return containerengine.ContainerEngineClient{}, err
	}
	return client, nil
}
