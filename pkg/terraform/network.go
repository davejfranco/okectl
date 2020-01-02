package terraform

import (
	"fmt"
	"log"
	"net"

	"github.com/davejfranco/okectl/pkg/util"
)

//Subnet describe a subnet to be use in vcn
type Subnet struct {
	Cidr   string
	Public bool //
}

//Rtable represents a Route Table on a VCN
type Rtable struct {
	Subnets []Subnet
	Routes  []string
}

type Vcn struct {
	Name         string
	Cidr         string
	CompartmenID string
}

//Network components
type Network struct {
	VcnID   string
	IgwID   string
	NgwID   string
	Cidr    string
	Rtables []Rtable
	Subnets []Subnet
}

//validate that CIDR provided is a valid one
func (vcn Vcn) validCIDR() bool {

	_, ipv4net, err := net.ParseCIDR(vcn.Cidr)
	if err != nil {
		log.Fatal(err)
		return false
	}

	//Check Network Mask lenght should be between either /16 or /24
	mask := util.HexaMask(ipv4net.Mask.String())
	fmt.Println(mask)
	if mask != "16" && mask != "24" {
		return false
	}
	return true
}

//ResourceGen will generate tfile resource
func (vcn Vcn) ResourceGen() TField {

	if vcn.validCIDR() {
		return TField{"vcn": TField{"cidr_block": vcn.Cidr, "compartment_id": vcn.CompartmenID, "display_name": vcn.Name}}
	}
	return TField{}
}
