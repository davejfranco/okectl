package ctl

import (
	"fmt"
	"okectl/pkg/oci"
	"okectl/pkg/template"
	"okectl/pkg/util"
	"os"
)

func NewCluster(clusterName, compartmentID, k8sVersion, region string) {

	//New Config
	config := oci.NewConfigProvider("", "")
	//Check kubernetes version
	availableVersions := oci.GetKubernetesVersion(config)
	//Don't know if I'll leave this here
	var counter int = 0
	for _, v := range availableVersions {
		if v == k8sVersion {
			break
		}
		counter++
	}
	if counter == len(availableVersions) {
		k8sVersion = availableVersions[len(availableVersions)-1] //return latest available version
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

	/*if err := template.RenderFile(t); err != nil {
		panic(err)
	}*/
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

}

func NewNodePool() {

	fmt.Println("Create NodePool called")
}
