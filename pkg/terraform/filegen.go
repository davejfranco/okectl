package terraform

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

const (
	tfFileName = "okectl.tf.json"
)

//Tfile generator struct
type Tfile struct {
	Provider []Field `json:"provider,omitempty"`
	Variable []Field `json:"variable,omitempty"`
	Data     []Field `json:"data,omitempty"`
	Resource []Field `json:"resource,omitempty"`
	Output   []Field `json:"output,omitempty"`
}

//Field to represent fields in terraform resources
type Field map[string]interface{}

//Genfile will generate a tf.json file to be use by terraform
func (tf Tfile) Genfile() error {

	tojson, err := json.MarshalIndent(tf, "", " ")
	if err != nil {
		return err
	}

	jsonFile, err := os.Create(tfFileName)
	if err != nil {
		return err
	}

	jsonFile.Write(tojson)

	return nil
}

func validResource(resource string) bool {

	var validr = []string{"oci_identity_compartment",
		"oci_core_vcn",
		"oci_core_internet_gateway",
		"oci_core_nat_gateway",
		"oci_core_route_table",
		"oci_core_subnet"}

	for _, v := range validr {
		if v == resource {
			return true
		}
	}
	return false
}

//Resource to identifu different resources in the terraform
type Resource struct {
	Type   string
	Name   string
	Rfield Field
}

//Create resource
func (r Resource) Create() Field {

	if r.Type == "" {
		log.Fatal("Resource type cannot be empty")
	}

	if !validResource(r.Type) {
		log.Fatalf("Invalid resource type %v", r.Type)
	}

	return Field{r.Type: Field{r.Name: r.Rfield}}
}

//Export will generate the resource export like "${oci_core_vcn.vcn.id}"
func (r *Resource) Export() string {
	return fmt.Sprintf("${%v.%v.id}", r.Type, r.Name)
}
