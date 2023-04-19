package template

import (
	"okectl/pkg/util"
	"testing"
)

func TestRenderFile(t *testing.T) {
	// Write test for RenderFile function
	template := Template{
		CidrBlock:     "10.0.0.0/16",
		Random:        util.RandomString(4),
		CompartmentID: "ocid1.tenancy.oc1..aaaaaaaadskkxvb5tsienhdyofk57mcglt4hhmtcocv3ppsryxd5fxhcufes",
		Cluster: Cluster{
			Name:    "testCluster",
			Version: "1.24",
		},
		NodePool: NodePool{
			Name:    "testNP1",
			Shape:   "VM.Standard.E3.Flex",
			ImageID: "ocid1.image.oc1.iad.aaaaaaaas2zhgcuhfarrwqxow4ffdrllxbfqkm32b4y3bovmatntgjvv56va",
			Size:    "1",
			ShapeConfig: ShapeConfig{
				RAM: "8",
				CPU: "2",
			},
		},
	}

	if err := RenderFile(template); err != nil {
		t.Errorf("RenderFile() error = %v", err)
	}
}
