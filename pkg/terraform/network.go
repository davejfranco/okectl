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

//Subnet describe a subnet to be use in vcn
type subnet struct {
	cidr          string
	vcnID         string
	rtID          string
	compartmentID string
}

//ValidCIDR will return error if a subnet is not part of the VCN CIDR
// func (sub subnet) validCIDR(vcnCIDR string) error {
// 	return nil
// }

//Create subnet fields
func (sub subnet) create() (Field, error) {

	if sub.cidr == "" || sub.vcnID == "" || sub.rtID == "" || sub.compartmentID == "" {
		return Field{}, errors.New("One or more subnet fields are empty")
	}

	return Field{"vcn_id": sub.vcnID,
		"cidr_block":     "10.0.1.0/24",
		"compartment_id": sub.compartmentID}, nil
}

//Vcn network
type vcn struct {
	name          string
	cidr          string
	compartmentID string
}

//Create VCN Resource
func (v vcn) create() (Field, error) {

	if v.name == "" {
		v.name = vcnName + util.RandomKey()
	}

	if !v.validVcnCIDR() {
		return Field{}, errors.New("Invalid CIDR for the VCN")
	}

	return Field{"cidr_block": v.cidr,
		"compartment_id": v.compartmentID,
		"display_name":   v.name}, nil
}

//ValidVcnCIDR that CIDR provided is a valid one
func (v vcn) validVcnCIDR() bool {

	_, ipv4net, err := net.ParseCIDR(v.cidr)
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

//Rtable represents a Route Table on a VCN
type rtable struct {
	vcnID         string
	compartmentID string
	name          string
	routes        []route
}

type route struct {
	netIdentityID   string
	destinationCIDR string
	destinationType string
}

//add routes to a given route table
func (rt *rtable) addRoutes(rts []route) {

	for _, r := range rts {
		rt.routes = append(rt.routes, r)
	}
}

func (rt rtable) create() (Field, error) {

	if len(rt.routes) == 0 {
		return Field{}, errors.New("At least one route should be added to create this field")
	}

	//Add all routes available
	var routeRules []Field
	for _, ro := range rt.routes {
		ro := Field{"network_entity_id": ro.netIdentityID,
			"destination":      ro.destinationCIDR,
			"destination_type": "CIDR_BLOCK"}

		routeRules = append(routeRules, ro)
	}
	return Field{"compartment_id": rt.compartmentID,
		"vcn_id":      rt.vcnID,
		"route_rules": routeRules}, nil
}

//Network components
type network struct {
	vcnID         vcn
	publicOnly    bool
	compartmentID string
	rtables       []rtable
	subnets       []subnet
}
