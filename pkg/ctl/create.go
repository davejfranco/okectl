package ctl

import (
	"errors"
	"fmt"
	"okectl/pkg/oci"
	"okectl/pkg/template"
	"okectl/pkg/util"
	"os"
)

// New Config
var config = oci.NewConfigProvider("", "")

func NewCluster(clusterName, compartmentID, k8sVersion, region string) error {

	okeClient, err := config.Oke()
	if err != nil {
		return err
	}
	//Check if cluster version is valid
	if !oci.IsValidKubernetesVersion(k8sVersion, &okeClient) {
		return errors.New("invalid kubernetes version") //TODO: Return error
	}

	if region == "" {
		region = config.ConfigDetails()["region"]
	}

	//Render file for Resource Manager
	t := template.Template{
		CidrBlock:     template.CidrBlock,
		Random:        util.RandomString(4),
		Region:        region,
		CompartmentID: compartmentID,
		Cluster: template.Cluster{
			Name:    clusterName,
			Version: k8sVersion,
		},
		NodePool: template.NodePool{
			Name:    template.NodePoolName,
			Shape:   template.NodePoolShape,
			ImageID: template.NodePoolImageID,
			Size:    template.NodePoolSize,
			ShapeConfig: template.ShapeConfig{
				RAM: template.NodePoolRAM,
				CPU: template.NodePoolCPU,
			},
		},
	}

	//Generate the main.tf file
	if err := template.RenderFile(t); err != nil {
		return err
	}

	zipFile, err := template.ZipAndEncodeTemplate(t)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//Time to call Resource Manager
	//Resource Manager client
	rmClient, err := config.ResourceManager()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	stack := oci.NewStack(rmClient)
	stack.Name = clusterName
	stack.CompartmentID = compartmentID

	_, err = stack.Create(zipFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//Apply the stack
	_, err = stack.Job("apply")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return nil
}

func NewNodePool() {

	fmt.Println("Create NodePool called")
}
