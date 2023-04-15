package resourcemanager

import (
	"os"
	"text/template"
)

const (
	cidrBlock        = "10.0.0.0/16"
	templateLocation = "pkg/resourcemanager/files/main.tf.tmpl"
)

type Cluster struct {
	Name    string
	Version string
}

type NodePool struct {
	Name  string
	Shape string
	Image string
	Size  string
}

// Template struct
type Template struct {
	CidrBlock     string
	Random        string
	CompartmentID string
	Cluster       Cluster
	NodePool      NodePool
}

// RenderFile renders a template file located in the files/main.tf.tmpl
func RenderFile(t Template) error {
	//Open the template file
	templateFile, err := os.Open(templateLocation)
	if err != nil {
		panic(err)
	}

	defer templateFile.Close()

	// Parse the template file
	tmpl, err := template.ParseFiles(templateLocation)
	if err != nil {
		panic(err)
	}

	// Create a new file to write the rendered Terraform code to
	file, err := os.Create("main.tf")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Render the template
	err = tmpl.Execute(file, t)
	if err != nil {
		panic(err)
	}

	return nil
}
