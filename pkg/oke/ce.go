package oke

import (
	"context"
	"fmt"

	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/containerengine"
)

//Oke struct describing container engine
type Oke struct {
	Client        *containerengine.ContainerEngineClient
	CompartmentID string
	Ctx           context.Context
}

//ClusterInfo is the struct to represent cluster on the clt
type Cluster struct {
	Name        string
	ID          string
	Status      string
	NodePools   []nodePool
	VcnID       string
	KubeVersion string
	Created     string
}

//Check if kubernetes version si s valid one
func (o Oke) validKubeVersion(version string) bool {

	validVersions, err := o.clusterAvailableVersions()
	if err != nil {
		return false
	}
	for _, v := range validVersions {

		if version == v {
			return true
		}
	}
	return false
}

//return list of kubernetes available versions
func (o Oke) clusterAvailableVersions() ([]string, error) {

	req := containerengine.GetClusterOptionsRequest{ClusterOptionId: common.String("all")}
	options, err := o.Client.GetClusterOptions(o.Ctx, req)
	if err != nil {
		return []string{}, err
	}
	return options.KubernetesVersions, nil
}

//NodePool Operations
func (o Oke) listClusterNodePools(clusterID string) ([]nodePool, error) {

	lreq := containerengine.ListNodePoolsRequest{
		CompartmentId: common.String(o.CompartmentID),
		ClusterId:     common.String(clusterID),
	}

	response, err := o.Client.ListNodePools(o.Ctx, lreq)
	if err != nil {
		return []nodePool{}, err
	}

	var pools []nodePool
	for _, np := range response.Items {
		p := nodePool{
			id:          *np.Id,
			name:        *np.Name,
			imageID:     *np.NodeImageId,
			nodeShape:   *np.NodeShape,
			kubeVersion: *np.KubernetesVersion,
			subnetIds:   np.SubnetIds,
		}

		pools = append(pools, p)
	}
	return pools, nil
}

func (o Oke) getClusterByName(clusterName string) ([]containerengine.ClusterSummary, error) {

	lcr := containerengine.ListClustersRequest{
		CompartmentId: common.String(o.CompartmentID),
		Name:          common.String(clusterName),
	}

	response, err := o.Client.ListClusters(o.Ctx, lcr)
	switch {
	case err != nil:
		return []containerengine.ClusterSummary{}, err
	case len(response.Items) == 0:
		return []containerengine.ClusterSummary{}, fmt.Errorf("No cluster found with name: %s", clusterName)
	}

	return response.Items, nil
}

//GetAllClusters returns all cluster in a given compartment
func (o Oke) GetAllClusters() ([]Cluster, error) {

	lcr := containerengine.ListClustersRequest{
		CompartmentId: common.String(o.CompartmentID),
	}

	response, err := o.Client.ListClusters(o.Ctx, lcr)
	if err != nil {
		return []Cluster{}, err
	}

	var output []Cluster
	for _, c := range response.Items {
		//not sure what to do when getting an error here
		np, _ := o.listClusterNodePools(*c.Id)

		cluster := Cluster{
			Name:        *c.Name,
			ID:          *c.Id,
			Status:      fmt.Sprintln(c.LifecycleState),
			NodePools:   np,
			VcnID:       *c.VcnId,
			KubeVersion: *c.KubernetesVersion,
			Created:     c.Metadata.TimeCreated.Format("02-01-2006"),
		}

		output = append(output, cluster)
	}
	return output, nil
}

//DeleteCluster to remove any cluster on a given compartment
func (o Oke) DeleteCluster(clusterName string) error {

	cluster, err := o.getClusterByName(clusterName)
	if err != nil {
		panic(err)
	}

	resp, err := o.Client.DeleteCluster(o.Ctx, containerengine.DeleteClusterRequest{ClusterId: cluster[0].Id})
	if err != nil {
		panic(err)
	}

	if resp.RawResponse.Status != "200" {
		return fmt.Errorf("Unable to delete selected cluster")
	}
	return nil
}
