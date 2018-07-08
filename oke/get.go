package oke

import (
	"context"

	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/containerengine"
)

//GetAllClusters returns all cluster on a given compartment
func (client *NewClient) GetAllClusters(compartmentID string) (containerengine.ListClustersResponse, error) {

	// client, err := NewDefaultClient()
	// if err != nil {
	// 	return containerengine.ListClustersResponse{}, err
	// }
	//context
	ctx := context.Background()

	lcr := containerengine.ListClustersRequest{
		CompartmentId: common.String(compartmentID),
	}

	response, err := client.Client.ListClusters(ctx, lcr)
	if err != nil {
		return containerengine.ListClustersResponse{}, err
	}
	return response, nil
}
