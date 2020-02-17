package terraform

import (
	"errors"
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

/*
resource "oci_core_vcn" "test_vcn" {
    #Required
    cidr_block = "${var.vcn_cidr_block}"
    compartment_id = "${var.compartment_id}"

    #Optional
    defined_tags = {"Operations.CostCenter"= "42"}
    display_name = "${var.vcn_display_name}"
    dns_label = "${var.vcn_dns_label}"
    freeform_tags = {"Department"= "Finance"}
    ipv6cidr_block = "${var.vcn_ipv6cidr_block}"
    is_ipv6enabled = "${var.vcn_is_ipv6enabled}"
}
*/
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

/*
resource "oci_core_route_table" "test_route_table" {
    #Required
    compartment_id = "${var.compartment_id}"
    vcn_id = "${oci_core_vcn.test_vcn.id}"

    #Optional
    defined_tags = {"Operations.CostCenter"= "42"}
    display_name = "${var.route_table_display_name}"
    freeform_tags = {"Department"= "Finance"}
    route_rules {
        #Required
        network_entity_id = "${oci_core_internet_gateway.test_internet_gateway.id}"

        #Optional
        cidr_block = "${var.route_table_route_rules_cidr_block}"
        description = "${var.route_table_route_rules_description}"
        destination = "${var.route_table_route_rules_destination}"
        destination_type = "${var.route_table_route_rules_destination_type}"
    }
}
*/
type route struct {
	netIdentityID, destinationCIDR string
}

func (r route) createRoute() Field {

	return Field{"network_entity_id": r.netIdentityID,
		"destination":      r.destinationCIDR,
		"destination_type": "CIDR_BLOCK"}
}

type rtable struct {
	vcnID, compartmentID, name string
}

//Create route table
func (rt rtable) create(r Field) (Resource, error) {

	rtr := Resource{Type: "oci_core_route_table",
		Name: "rt",
		Rfield: Field{"compartment_id": rt.compartmentID,
			"vcn_id":       rt.vcnID,
			"display_name": rt.name,
			"route_rules":  r}}

	return rtr, nil
}

/*
resource "oci_core_subnet" "test_subnet" {
    #Required
    cidr_block = "${var.subnet_cidr_block}"
    compartment_id = "${var.compartment_id}"
    vcn_id = "${oci_core_vcn.test_vcn.id}"

    #Optional
    availability_domain = "${var.subnet_availability_domain}"
    defined_tags = {"Operations.CostCenter"= "42"}
    dhcp_options_id = "${oci_core_dhcp_options.test_dhcp_options.id}"
    display_name = "${var.subnet_display_name}"
    dns_label = "${var.subnet_dns_label}"
    freeform_tags = {"Department"= "Finance"}
    ipv6cidr_block = "${var.subnet_ipv6cidr_block}"
    prohibit_public_ip_on_vnic = "${var.subnet_prohibit_public_ip_on_vnic}"
    route_table_id = "${oci_core_route_table.test_route_table.id}"
    security_list_ids = "${var.subnet_security_list_ids}"
}
*/
//Subnet describe a subnet to be use in vcn
type subnet struct {
	vcnID, cidr, rtID, compartmentID, name string
}

//Create subnet fields
func (sub subnet) create() (Resource, error) {

	if sub.cidr == "" || sub.vcnID == "" || sub.rtID == "" || sub.compartmentID == "" {
		return Resource{}, errors.New("One or more subnet fields are empty")
	}

	s := Resource{Type: "oci_core_subnet",
		Name: "sub",
		Rfield: Field{"vcn_id": sub.vcnID,
			"cidr_block":     sub.cidr,
			"display_name":   sub.name,
			"compartment_id": sub.compartmentID,
			"route_table_id": sub.rtID}}

	return s, nil
}

//security list struct
type secList struct {
	name, vncID, compartmentID string
	egressRule, ingressRule    Field
}
