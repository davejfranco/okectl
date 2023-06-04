package ctl

import (
	"fmt"
)

func CreateCluster(clusterName, compartmentID, k8sVersion string) {
	fmt.Printf("Create Cluster called with name: %s\nCompartmentID: %s\nKubernetes version: %s\n", clusterName, compartmentID, k8sVersion)
}

func CreateClusterFromFile(filePath string) {
	fmt.Printf("Create Cluster from file called with file %s\n", filePath)
}

func CreateNodePool() {

	fmt.Println("Create NodePool called")
}
