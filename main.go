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
	"context"
	"fmt"

	//"github.com/davejfranco/okectl/cmd"
	"github.com/davejfranco/okectl/pkg/oci"
)

func main() {
	//cmd.Execute()
	config := oci.Config{}
	cfg, err := config.Load()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(cfg)

	client, err := oci.VcnClient(cfg)
	if err != nil {
		fmt.Println(err)
	}

	ctx := context.Background()
	compartmentid := "ocid1.compartment.oc1..aaaaaaaahm7337hqcgxgdvzckf3kvwvsltwtftuscqajeydkcm2m5gcdy6ka"

	netClient := oci.Vcn{&client, ctx}
	if err = netClient.CreateVCN(compartmentid, ctx); err != nil {
		fmt.Println(err)
	}
}
