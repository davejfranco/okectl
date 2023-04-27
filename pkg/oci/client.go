package oci

import (
	"context"

	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/containerengine"
	"github.com/oracle/oci-go-sdk/resourcemanager"
)

// Config
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

// Container Engine
// create client using config in default location
func (c *Config) OKE() (containerengine.ContainerEngineClient, error) {
	client, err := containerengine.NewContainerEngineClientWithConfigurationProvider(c.configProvider())
	if err != nil {
		return containerengine.ContainerEngineClient{}, err
	}
	return client, nil
}

// Resource Manager
func (c *Config) ResourceManager() (resourcemanager.ResourceManagerClient, error) {
	client, err := resourcemanager.NewResourceManagerClientWithConfigurationProvider(c.configProvider())
	if err != nil {
		return resourcemanager.ResourceManagerClient{}, err
	}
	return client, nil
}

type ResourceManager struct {
	Client resourcemanager.ResourceManagerClient
}

func (rm ResourceManager) ListStacks(ctx context.Context, request resourcemanager.ListStacksRequest) (resourcemanager.ListStacksResponse, error) {
	return rm.Client.ListStacks(ctx, request)
}

func (rm ResourceManager) CreateStack(ctx context.Context, req resourcemanager.CreateStackRequest) (resourcemanager.CreateStackResponse, error) {
	return rm.Client.CreateStack(ctx, req)
}

func (rm ResourceManager) DeleteStack(ctx context.Context, req resourcemanager.DeleteStackRequest) (resourcemanager.DeleteStackResponse, error) {
	return rm.Client.DeleteStack(ctx, req)
}
