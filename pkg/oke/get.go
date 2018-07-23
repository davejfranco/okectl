package oke

import (
	"context"
	"fmt"
	"os"

	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/containerengine"
)

//Oke struct describing container engine
type Oke struct {
	Client       *containerengine.ContainerEngineClient
	CompartmenID string
	Ctx          context.Context
}

//GetAllClusters returns all cluster on a given compartment
func (o Oke) getAllClustersIds() ([]string, error) {

	lcr := containerengine.ListClustersRequest{
		CompartmentId: common.String(o.CompartmenID),
	}

	var clusterIds []string
	response, err := o.Client.ListClusters(o.Ctx, lcr)
	if err != nil {
		return clusterIds, err
	}

	for _, item := range response.Items {
		clusterIds = append(clusterIds, *item.Id)
	}
	return clusterIds, nil
}

//DescribeCluster returns info of a given cluster
func (o Oke) DescribeCluster(clusterID string) (containerengine.Cluster, error) {

	ceRequest := containerengine.GetClusterRequest{
		ClusterId: common.String(clusterID),
	}

	resp, err := o.Client.GetCluster(o.Ctx, ceRequest)
	if err != nil {
		return containerengine.Cluster{}, err
	}

	return resp.Cluster, nil
}

//ClustersInfo prints all cluster related info on a given compartment
func (o Oke) ClustersInfo() []containerengine.Cluster {

	allClustersIds, err := o.getAllClustersIds()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var allclusterinfo []containerengine.Cluster
	for _, id := range allClustersIds {
		info, err := o.DescribeCluster(id)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		allclusterinfo = append(allclusterinfo, info)
	}

	return allclusterinfo
}
