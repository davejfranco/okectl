package main

import (
	"context"
	"fmt"

	"github.com/davejfranco/okectl/pkg/oke"
	"github.com/davejfranco/okectl/pkg/util"
)

func main() {
	OkeTest()
}

func OkeTest() {

	//user := util.OciUser{KeyPassphrasse: ""}
	c := util.Config{
		Path:    "~/.oci/configProfile",
		Profile: "DEV",
	}
	//c := util.Config{}
	config, err := c.Load()
	if err != nil {
		panic(err)
	}
	ceClient, err := util.OkeClient(config)
	if err != nil {
		fmt.Errorf("Error: %s", err)
	}

	test := oke.Oke{
		Client:        &ceClient,
		CompartmentID: "ocid1.compartment.oc1..aaaaaaaan3ck32sx7mhiioxoxnpwtzcziyqzqi5pftawcnql47zvabbglgaa",
		Ctx:           context.Background(),
	}

	oke.PrintAllClusters(test)

}
