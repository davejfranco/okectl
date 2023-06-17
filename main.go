/*
Copyright © 2022 Dave Franco davefranco1987@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"okectl/cmd"
	"okectl/pkg/template"
	"okectl/pkg/util"
)

var compartmentID string = "ocid1.compartment.oc1..aaaaaaaan5u7sgad2xlxamnxi4nmb5zvtdrww62ec4i22qvgs2mqtmwczkea"

/*
func rmclient() resourcemanager.ResourceManagerClient {
	config := oci.Config{}
	rmclient, err := config.ResourceManager()
	if err != nil {
		panic(err)
	}
	return rmclient
}*/

func testTemplate() string {
	t := template.Template{
		CidrBlock:     template.CidrBlock,
		Random:        util.RandomString(4),
		Region:        "us-ashburn-1",
		CompartmentID: compartmentID,
		Cluster: template.Cluster{
			Name: "devCluster",
			//Version: "1.24",
		},
		NodePool: template.NodePool{
			Name:    "devNP1",
			Shape:   "VM.Standard.E3.Flex",
			ImageID: "ocid1.image.oc1.iad.aaaaaaaas2zhgcuhfarrwqxow4ffdrllxbfqkm32b4y3bovmatntgjvv56va",
			Size:    "1",
			ShapeConfig: template.ShapeConfig{
				RAM: "8",
				CPU: "2",
			},
		},
	}

	if err := template.RenderFile(t); err != nil {
		panic(err)
	}
	return "ok"

	//return file
	// file, err := util.ZipAndEncodeFile("renderedmain.tf")
	// if err != nil {
	// 	panic(err)
	// }
	//zip, err := template.ZipAndEncodeTemplate(t)
	//if err != nil {
	//	panic(err)
	//}
	//return zip
}

/*
func testStack() oci.Stack {

	cfg := oci.Config{}
	client, err := cfg.ResourceManager()
	if err != nil {
		panic(err)
	}

	stack := oci.NewStack(client)
	stack.Name = "devStack"
	stack.CompartmentID = compartmentID
	return *stack
}*/

/*
// Create a new client
func testStackCreation() {
	// client := rmclient()
	zip := testTemplate()
	stack := testStack()

	resp, err := stack.Create(zip)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp)
}*/

/*
func testGetStack(name, compartmentID string) {
	cfg := oci.Config{}
	client, err := cfg.ResourceManager()
	if err != nil {
		panic(err)
	}
	stack, err := oci.GetStack(name, compartmentID, client)
	if err != nil {
		panic(err)
	}
	fmt.Println(stack)
}*/

func main() {
	cmd.Execute()
	//testTemplate()

}
