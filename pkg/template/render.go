package template

import (
	"fmt"
	"os"
	"text/template"
)

const (
	CidrBlock string = "10.0.0.0/16"
)

type Cluster struct {
	Name    string
	Version string
}

type ShapeConfig struct {
	RAM string
	CPU string
}

type NodePool struct {
	Name        string
	Shape       string
	ImageID     string
	Size        string
	ShapeConfig ShapeConfig
}

// Template struct
type Template struct {
	CidrBlock     string
	Random        string
	CompartmentID string
	Cluster       Cluster
	NodePool      NodePool
}

func (t Template) Generate() Template {

	if t.CidrBlock == "" {
		t.CidrBlock = CidrBlock
	}

	if t.Cluster.Name == "" {
		t.Cluster.Name = "okectl"
	}

	if t.Cluster.Version == "" {
		t.Cluster.Version = "1.24" //TODO: Get the latest version from the API
	}

	return t
}

// RenderFile renders a template file located in the files/main.tf.tmpl
func RenderFile(t Template) error {

	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	templateLocation := fmt.Sprintf("%s/%s", currentDir, "files/main.tf.tmpl") //"files/main.tf.tmpl"
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
