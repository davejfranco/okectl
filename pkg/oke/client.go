package oke

import (
	"context"

	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/containerengine"
)

//Conn struct to specify config and context for oracle's services
type Conn struct {
	Config  common.ConfigurationProvider
	Context *context.Context
}

//Client creates a client connection to oracle container service
func Client(config common.ConfigurationProvider) (containerengine.ContainerEngineClient, error) {

	client, err := containerengine.NewContainerEngineClientWithConfigurationProvider(config)
	if err != nil {
		return containerengine.ContainerEngineClient{}, err
	}
	return client, nil
}
