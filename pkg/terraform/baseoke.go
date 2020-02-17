package terraform

import (
	"github.com/davejfranco/okectl/pkg/util"
)

//BaseOKE will generate tf.json file of default OKE cluster including vcn
func vcnResources(cid string) []Resource {

	networkResources := []Resource{}

	//VCN
	vname := "vcn-" + util.RandomKey(4)
	v := vcn{
		name:          vname,
		cidr:          "172.16.0.0/16",
		compartmentID: cid,
	}
	net, _ := v.create()
	networkResources = append(networkResources, net)

	//Internet Gateway
	gw := gateway{
		vcnID:         net.Export(),
		compartmentID: cid,
	}
	igw := gw.createIGW()
	networkResources = append(networkResources, igw)

	//Default Route to Internet gateway
	r := route{
		netIdentityID:   igw.Export(),
		destinationCIDR: "0.0.0.0/0",
	}

	rt := rtable{
		vcnID:         net.Export(),
		compartmentID: cid,
		name:          "Public Route Table",
	}

	//Public Route Table
	prt, _ := rt.create(r.createRoute())
	networkResources = append(networkResources, prt)

	//Public Subnet
	s := subnet{
		vcnID:         net.Export(),
		cidr:          "172.16.1.0/24",
		rtID:          prt.Export(),
		compartmentID: cid,
		name:          "public_subnet_1",
	}

	psub, _ := s.create()
	networkResources = append(networkResources, psub)
	return networkResources

}

//OKE Terraform Resources
func okeResources() {}

//BaseOKE will generate basic cluster
func BaseOKE(cid string) error {

	base := Tfile{}

	oci := Field{"oci": Field{"version": "v3.62.0",
		"region":           "${var.region}",
		"tenancy_ocid":     "${var.tenancy_ocid}",
		"user_ocid":        "${var.user_ocid}",
		"fingerprint":      "${var.fingerprint}",
		"private_key_path": "${var.private_key_path}"}}

	//set default provider
	base.Provider = []Field{oci}

	v := vcnResources(cid)

	for _, v := range v {
		base.Resource = append(base.Resource, v.Create())
	}

	err := base.Genfile()
	if err != nil {
		return err
	}

	return nil
}
