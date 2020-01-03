package terraform

import (
	"encoding/json"
	"log"
	"os"
)

const (
	tfFileName = "maingo.tf.json"
	k8sVersion = "v1.14.8"
	sshKeyfile = "~/.ssh/id_rsa.pub"
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
		"oci_core_route_table",
		"oci_core_subnet"}

	for _, v := range validr {
		if v == resource {
			return true
		}
	}
	return false
}

func validVarTypes(vtype string) bool {
	var validtypes = []string{"string", "list", ""}
}

//Var respresents Variable in tf files
type Var struct {
	Vname    string
	Vtype    string
	Vdefault string
}

//Resource to identifu different resources in the terraform
type Resource struct {
	Rtype  string
	Rname  string
	Rfield Field
}

//Create resource
func (r *Resource) Create() Field {

	if r.Rtype == "" {
		log.Fatal("Resource type cannot be empty")
	}

	if !validResource(r.Rtype) {
		log.Fatalf("Invalid resource type %v", r.Rtype)
	}

	return Field{r.Rtype: Field{r.Rname: r.Rfield}}
}
