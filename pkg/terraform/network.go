package terraform

import (
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/davejfranco/okectl/pkg/util"
)

const (
	subwrkName = "oke-wrk-subnet-quick-"
	sublbName  = "oke-lb-subnet-quick-"
	seclistWrk = "oke-wrk-seclist-"
	seclistlb  = "oke-lb-seclist-"
	igwname    = "oke-igw-quick-"
	ngwname    = "oke-ngw-quick-"
	vcnName    = "oke-vcn-quick-"
	rtwrk      = "oke-wrk-rt-"
	rtlb       = "oke-lb-rt-"
)

//Network components
type networkVCN struct {
	vcn, igw, ngw Resource
	publicOnly    bool
	compartmentID string
	subnets       []subnet
	rtables       []rtable
}

type vcn struct {
	name, cidr, compartmentID string
}

//Create VCN Resource
func (v vcn) create() (Resource, error) {

	name := vcnName + util.RandomKey(4)
	if !validVcnCIDR(v.cidr) {
		return Resource{}, errors.New("Invalid CIDR for the VCN")
	}

	vcn := Resource{Type: "oci_core_vcn",
		Name: "vcn", Rfield: Field{"cidr_block": v.cidr,
			"compartment_id": v.compartmentID,
			"display_name":   name}}

	return vcn, nil
}

//ValidVcnCIDR that CIDR provided is a valid one
func validVcnCIDR(cidr string) bool {

	_, ipv4net, err := net.ParseCIDR(cidr)
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

type gateway struct {
	vcnID, compartmentID string
}

func (gw gateway) createIGW() Resource {
	return Resource{Type: "oci_core_internet_gateway",
		Name: "igw",
		Rfield: Field{"compartment_id": gw.compartmentID,
			"vcn_id": gw.vcnID}}
}

func (gw gateway) createNGW() Resource {
	return Resource{Type: "oci_core_nat_gateway",
		Name: "ngw",
		Rfield: Field{"compartment_id": gw.compartmentID,
			"vcn_id": gw.vcnID}}
}

type route struct {
	netIdentityID, destinationCIDR, destinationType string
}

func (r route) createRoute() Field {

	return Field{"network_entity_id": r.netIdentityID,
		"destination":      r.destinationCIDR,
		"destination_type": "CIDR_BLOCK"}
}

type rtable struct {
	routeRule            Field
	vcnID, compartmentID string
}

func (rt *rtable) create() (Resource, error) {

	f := &Field{}
	if &rt.routeRule == f {
		return Resource{}, errors.New("A rule is required")
	}

	rtr := Resource{Type: "oci_core_route_table",
		Name: "rt",
		Rfield: Field{"compartment_id": rt.compartmentID,
			"vcn_id":      rt.vcnID,
			"route_rules": rt.routeRule}}

	return rtr, nil
}

//Subnet describe a subnet to be use in vcn
type subnet struct {
	vcnID, cidr, rtID, compartmentID string
}

//Create subnet fields
func (sub subnet) create() (Resource, error) {

	if sub.cidr == "" || sub.vcnID == "" || sub.rtID == "" || sub.compartmentID == "" {
		return Resource{}, errors.New("One or more subnet fields are empty")
	}

	s := Resource{Type: "oci_core_subnet",
		Name: "pub_elb_subnet",
		Rfield: Field{"vcn_id": sub.vcnID,
			"cidr_block":     sub.cidr,
			"compartment_id": sub.compartmentID,
			"route_table_id": sub.rtID}}

	return s, nil
}

//security list struct
type secList struct {
	name, vncID, compartmentID string
	egressRule, ingressRule    Field
}
