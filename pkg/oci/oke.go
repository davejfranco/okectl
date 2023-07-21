package oci

import (
	"context"
	"errors"
	"okectl/pkg/util"
	"strings"
	"syscall"

	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/containerengine"
)

// Container Engine
// create client using config in default location
func (c *Config) Oke() (containerengine.ContainerEngineClient, error) {
	client, err := containerengine.NewContainerEngineClientWithConfigurationProvider(c.provider)
	if err != nil {
		return containerengine.ContainerEngineClient{}, err
	}
	return client, nil
}

// GetKubernetesVersion returns a list of available Kubernetes versions
func getKubernetesVersion(client *containerengine.ContainerEngineClient) []string {

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

// isValidKubernetesVersion checks if the kubernetes version is valid
func IsValidKubernetesVersion(k8sVersion string, client *containerengine.ContainerEngineClient) bool {

	availableVersions := getKubernetesVersion(client)
	for _, v := range availableVersions {
		//check if the version is in the list of available versions, doesn't matter if it's v1.19.10 or 1.19.10
		if strings.Contains(v, k8sVersion) {
			return true
		}
	}
	return false
}

// Return a list of available shapes and image ocid
func getNodePoolOptions(client *containerengine.ContainerEngineClient) containerengine.GetNodePoolOptionsResponse {
	getNodePoolOptionsReq := containerengine.GetNodePoolOptionsRequest{
		NodePoolOptionId: common.String("all"),
	}
	getNodePoolOptionsResp, err := client.GetNodePoolOptions(context.Background(), getNodePoolOptionsReq)
	util.FatalIfError(err)

	return getNodePoolOptionsResp
}

func GetNodePoolImageID(client *containerengine.ContainerEngineClient, kubernetesVersion string) string {
	if !IsValidKubernetesVersion(kubernetesVersion, client) {
		util.FatalIfError(errors.New("invalid kubernetes version"))
		syscall.Exit(0)
	}
	np := getNodePoolOptions(client)
	//fmt.Println(np)

	/*
		Oracle-Linux-8.7-aarch64-2023.05.24-0-OKE-1.26.2-625 ImageId=ocid1.image.oc1.iad.aaaaaaaawcrzxwwinxasinuhadfrs5qpc7yqibbc5facc54fqdpahbsfnntq
		Oracle-Linux-8.7-Gen2-GPU-2023.04.27-0-OKE-1.26.2-607 ImageId=ocid1.image.oc1.iad.aaaaaaaa5hrle4c3pfwytllp5vgmlygchedb34nf52czrqovtrqdqnir5gaq
		Oracle-Linux-8.7-2023.05.24-0-OKE-1.26.2-625 ImageId=ocid1.image.oc1.iad.aaaaaaaaqvn4ubp2zfm5xagjaelgeg6vwrbru6hfpocmqrqxiidkp5tstqiq
		[Oracle Linux Version]-[date]-OKE-[k8s Version]

	*/
	//searchStr := fmt.Sprintf("%s-%s-%s", "Oracle-Linux", linuxVersion, currentYear)
	var sources []containerengine.NodeSourceOption
	for _, src := range np.Sources {
		/*if strings.Contains(*src.GetSourceName(), searchStr) &&
			strings.Contains(*src.GetSourceName(), kubernetesVersion) {
			return *src.(containerengine.NodeSourceViaImageOption).ImageId
		}*/

		if strings.Contains(*src.GetSourceName(), strings.Trim(kubernetesVersion, "v")) &&
			!strings.Contains(*src.GetSourceName(), "GPU") && !strings.Contains(*src.GetSourceName(), "aarch64") {
			sources = append(sources, src)
		}
	}

	//sort sources by name
	//sort.Slice(sources, func(i, j int) bool {
	//	return *sources[i].GetSourceName() > *sources[j].GetSourceName()
	//})

	return *sources[0].(containerengine.NodeSourceViaImageOption).ImageId
}
