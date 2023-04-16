package oke

import (
	"context"
	"fmt"
	"okectl/pkg/util"

	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/containerengine"
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

// GetKubernetesVersion returns a list of available Kubernetes versions
func GetKubernetesVersion(client containerengine.ContainerEngineClient) []string {
	getClusterOptionsReq := containerengine.GetClusterOptionsRequest{
		ClusterOptionId: common.String("all"),
	}
	getClusterOptionsResp, err := client.GetClusterOptions(context.Background(), getClusterOptionsReq)
	util.FatalIfError(err)
	kubernetesVersion := getClusterOptionsResp.KubernetesVersions

	if len(kubernetesVersion) < 1 {
		fmt.Println("Kubernetes version not available")
	}

	return kubernetesVersion
}
