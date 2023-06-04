package oci

import (
	"context"
	"errors"
	"fmt"
	"okectl/pkg/util"
	"strings"

	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/containerengine"
)

// GetKubernetesVersion returns a list of available Kubernetes versions
func GetKubernetesVersion() []string {

	config := Config{
		Profile: "DEFAULT",
	}

	client, err := config.Oke()
	util.FatalIfError(err)

	getClusterOptionsReq := containerengine.GetClusterOptionsRequest{
		ClusterOptionId: common.String("all"),
	}
	getClusterOptionsResp, err := client.GetClusterOptions(context.Background(), getClusterOptionsReq)
	util.FatalIfError(err)
	kubernetesVersion := getClusterOptionsResp.KubernetesVersions

	if len(kubernetesVersion) < 1 {
		util.FatalIfError(errors.New("error returning kubernetes versions"))
	}

	return kubernetesVersion
}

// Return a list of available shapes and image ocid
func GetNodePoolOptions(client containerengine.ContainerEngineClient) containerengine.GetNodePoolOptionsResponse {
	getNodePoolOptionsReq := containerengine.GetNodePoolOptionsRequest{
		NodePoolOptionId: common.String("all"),
	}
	getNodePoolOptionsResp, err := client.GetNodePoolOptions(context.Background(), getNodePoolOptionsReq)
	util.FatalIfError(err)

	return getNodePoolOptionsResp
}

func GetNodePoolImageID(client containerengine.ContainerEngineClient, linuxVersion, kubernetesVersion string) string {
	np := GetNodePoolOptions(client)
	searchStr := fmt.Sprintf("%s-%s-%s", "Oracle-Linux", linuxVersion, currentYear)
	for _, src := range np.Sources {
		if strings.Contains(*src.GetSourceName(), searchStr) &&
			strings.Contains(*src.GetSourceName(), kubernetesVersion) {
			return *src.(containerengine.NodeSourceViaImageOption).ImageId
		}
	}
	return ""
}
