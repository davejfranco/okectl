package template

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"text/template"
)

const (
	CidrBlock       string = "10.0.0.0/16"
	NodePoolName    string = "NP1"
	NodePoolShape   string = "VM.Standard.E3.Flex"
	NodePoolSize    string = "1"
	NodePoolImageID string = "ocid1.image.oc1.iad.aaaaaaaaqvn4ubp2zfm5xagjaelgeg6vwrbru6hfpocmqrqxiidkp5tstqiq"
	NodePoolRAM     string = "2"
	NodePoolCPU     string = "1"

	DefaultRegion          string = "us-ashburn-1"
	DefultTemplateLocation string = "pkg/template/files/main.tf.tmpl"
)

type Cluster struct {
	Name          string
	Type          string //Either BASIC_CLUSTER or ENHANCED_CLUSTER
	Version       string
	CompartmentID string
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

type VCN struct {
	ID            string
	Name          string
	CidrBlock     string
	CompartmentID string
}

// Template struct
type Template struct {
	Vcn              VCN
	Cluster          Cluster
	NodePool         NodePool
	CompartmentID    string
	TemplateLocation string
}

func (t Template) Generate() Template {

	if t.Vcn.CidrBlock == "" {
		t.Vcn.CidrBlock = CidrBlock
	}

	if t.Cluster.Name == "" {
		t.Cluster.Name = "okectl"
	}

	if t.Cluster.Version == "" {
		t.Cluster.Version = "1.24" //TODO: Get the latest version from the API
	}

	return t
}

// RenderToFile renders the template to a file, I'll leave this here for now
func (t Template) RenderToFile() error {

	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	dir := okectlDir() //Check if the .okectl directory exists and create it if it doesn't
	var renderedfile string = fmt.Sprintf("%s/%s", dir, "main.tf")

	templateLocation := fmt.Sprintf("%s/%s", currentDir, t.TemplateLocation)

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
	file, err := os.Create(renderedfile)
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

// RenderToEncodedZip renders the template to a zip file and returns the base64 encoded string
// This is a requirement for the OCI SDK Resouce Manager
func (t Template) RenderToEncodedZip() (string, error) {
	// Buffer to store the compressed data
	var zipBuffer bytes.Buffer

	// zip writer
	zipWriter := zip.NewWriter(&zipBuffer)

	// Create a new zip file header
	fileHeader := &zip.FileHeader{
		Name:   "main.tf",
		Method: zip.Deflate,
	}

	// Open a writer for the zip file
	zipFileWriter, err := zipWriter.CreateHeader(fileHeader)
	if err != nil {
		return "", err
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	templateLocation := fmt.Sprintf("%s/%s", currentDir, t.TemplateLocation)
	//Open the template file
	templateFile, err := os.Open(templateLocation)
	if err != nil {
		return "", err
	}

	defer templateFile.Close()

	// Parse the template file
	tmpl, err := template.ParseFiles(templateLocation)
	if err != nil {
		return "", err
	}

	err = tmpl.Execute(zipFileWriter, t)
	if err != nil {
		return "", err
	}

	// Close the zip writer
	err = zipWriter.Close()
	if err != nil {
		return "", err
	}

	// Encode the buffer's contents using base64 encoding
	base64Encoded := base64.StdEncoding.EncodeToString(zipBuffer.Bytes())

	return base64Encoded, nil

}
