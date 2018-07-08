package oke

import (
	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/containerengine"
)

//NewDefaultClient returns a client using default config on ~/.oci/config
func NewDefaultClient() (containerengine.ContainerEngineClient, error) {

	//ctx := context.Background()

	//create client
	client, err := containerengine.NewContainerEngineClientWithConfigurationProvider(common.DefaultConfigProvider())
	if err != nil {
		return containerengine.ContainerEngineClient{}, err
	}

	return client, nil
}

//NewClient creates a client to perform ops on Oracle's container engine service
type NewClient struct {
	Client containerengine.ContainerEngineClient
}
