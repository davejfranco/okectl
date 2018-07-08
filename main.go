package main

import (
	"fmt"
	"os"

	"github.com/davejfranco/okectl/oke"
)

func main() {

	c, err := oke.NewDefaultClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	client := oke.NewClient{Client: c}
	compartmentid := "ocid1.compartment.oc1..aaaaaaaan3ck32sx7mhiioxoxnpwtzcziyqzqi5pftawcnql47zvabbglgaa"
	r, err := client.GetAllClusters(compartmentid)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(r)

}
