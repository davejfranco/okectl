package terraform

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	tfFileName = "maingo.tf.json"
)

//Tfile generator struct
type Tfile struct {
	Provider []TField `json:"provider,omitempty"`
	Variable []TField `json:"variable,omitempty"`
	Data     []TField `json:"data,omitempty"`
	Resource []TField `json:"resource,omitempty"`
	Output   []TField `json:"output,omitempty"`
}

//TField to represent fields in terraform resources
type TField map[string]interface{}

//Genfile will generate a tf.json file to be use by terraform
func (tf *Tfile) Genfile() error {

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

func main() {

	var res []TField

	oci := TField{"oci": TField{"version": "3.0.0", "region": "us-ashburn-1"}}

	res = append(res, oci)

	tf := Tfile{Provider: res}
	err := tf.Genfile()
	if err != nil {
		fmt.Println(err.Error())
	}

}
