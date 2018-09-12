package oke

import (
	"context"

	"github.com/oracle/oci-go-sdk/core"
)

//Compute struct
type Compute struct {
	Client        core.ComputeClient
	Ctx           context.Context
	CompartmentID string
}

// func (c Compute) validShape() bool {

// 	for _, shape := range validNodeShapes {
// 		if n.nodeShape == shape {
// 			return true
// 		}
// 	}
// 	return false
// }

type nodePool struct {
	id          string
	name        string
	imageID     string
	nodeShape   string
	kubeVersion string
	subnetIds   []string
}

func (n nodePool) validImageID() bool {
	if n.imageID != nodeDefaultImageOs {
		return false
	}
	return true
}
