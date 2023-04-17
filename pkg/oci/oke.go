package oci

import (
	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/containerengine"
)

const (
	// This year
	currentYear = "2023"
)

type Options struct {
	Profile        string
	ConfigLocation string
}

// create client using config in default location
func NewOKEClient(options ...*Options) (containerengine.ContainerEngineClient, error) {

	if options != nil {
		conf := common.CustomProfileConfigProvider(options[0].Profile, options[0].ConfigLocation)
		return containerengine.NewContainerEngineClientWithConfigurationProvider(conf)
	}
	// Return client with specified profile
	return containerengine.NewContainerEngineClientWithConfigurationProvider(common.DefaultConfigProvider())
}
