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

//ClusterDetail is the struct to represent cluster on the clt
type ClusterDetail struct {
	Name        string
	ID          string
	Status      string
	NodePools   []nodePool
	VcnID       string
	KubeVersion string
	Created     string
}

type nodePool struct {
	id          string
	name        string
	imageID     string
	nodeShape   string
	kubeVersion string
	subnetIds   []string
}

func (o Oke) rawAllClusters() ([]containerengine.ClusterSummary, error) {

	lcr := containerengine.ListClustersRequest{
		CompartmentId: common.String(o.CompartmentID),
	}

	response, err := o.Client.ListClusters(o.Ctx, lcr)
	if err != nil {
		return []containerengine.ClusterSummary{}, err
	}

	return response.Items, nil
}

func (o Oke) getClusterByName(clusterName string) (containerengine.ClusterSummary, error) {

	lcr := containerengine.ListClustersRequest{
		CompartmentId: common.String(o.CompartmentID),
		Name:          common.String(clusterName),
	}

	response, err := o.Client.ListClusters(o.Ctx, lcr)
	if err != nil {
		return containerengine.ClusterSummary{}, err
	}

	return response.Items[0], nil
}

func (o Oke) clusterNodePools(clusterID string) ([]nodePool, error) {

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

func (o Oke) getClusterByID(clusterID string) (containerengine.Cluster, error) {

	ceRequest := containerengine.GetClusterRequest{
		ClusterId: common.String(clusterID),
	}

	resp, err := o.Client.GetCluster(o.Ctx, ceRequest)
	if err != nil {
		return containerengine.Cluster{}, err
	}

	return resp.Cluster, nil
}

//getAllClusters returns all cluster in a given compartment
func (o Oke) getAllClusters() ([]ClusterDetail, error) {

	all, err := o.rawAllClusters()
	if err != nil {
		return []ClusterDetail{}, nil
	}

	var output []ClusterDetail
	for _, c := range all {
		//not sure what to do when getting an error here
		np, _ := o.clusterNodePools(*c.Id)

		cluster := ClusterDetail{
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
