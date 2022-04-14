package oci

import (
	"context"
	"errors"

	"github.com/davejfranco/okectl/pkg/util"
	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/core"
)

var (
	defaultName        = "okectl_quick"
	defaultCIDR string = "10.0.0.0/16"
)

type Vcn struct {
	Client        *core.VirtualNetworkClient
	CompartmentID string
	Ctx           context.Context
}

//Network to be used con virtual network ops
type Network struct {
	Name          string
	CIDR          string
	CompartmentID string
}

//DescribeVcn returns vcn details of a given vcnID
func (v Vcn) DescribeVcn(vcnID string) (core.Vcn, error) {

	req := core.GetVcnRequest{VcnId: common.String(vcnID)}
	vcnDetail, err := v.Client.GetVcn(v.Ctx, req)
	if err != nil {
		return core.Vcn{}, err
	}
	return vcnDetail.Vcn, nil

}

//CreateVCN will deploy a vcn in a given compartmentID
func (v Vcn) Create(net Network) (core.CreateVcnResponse, error) {

	req := core.CreateVcnRequest{CreateVcnDetails: core.CreateVcnDetails{
		CompartmentId: &v.CompartmentID,
		CidrBlock:     &net.CIDR,
		DisplayName:   &net.Name,
	},
	}

	resp, err := v.Client.CreateVcn(v.Ctx, req)
	if err != nil {
		return core.CreateVcnResponse{}, err
	}

	return resp, nil
}

func (v Vcn) addRouteTable(displayName, vcnID string, routes []core.RouteRule) (core.CreateRouteTableResponse, error) {
	req := core.CreateRouteTableRequest{
		CreateRouteTableDetails: core.CreateRouteTableDetails{
			CompartmentId: &v.CompartmentID,
			RouteRules:    routes,
			VcnId:         &vcnID,
			DisplayName:   &displayName,
		},
	}
	resp, err := v.Client.CreateRouteTable(v.Ctx, req)
	if err != nil {
		return core.CreateRouteTableResponse{}, err
	}
	return resp, nil
}

//AddInternetGateway will create a IG to allow access from the internet
func (v Vcn) AddInternetGateway(vcnID, displayName string) (core.CreateInternetGatewayResponse, error) {

	isenabled := true
	req := core.CreateInternetGatewayRequest{
		CreateInternetGatewayDetails: core.CreateInternetGatewayDetails{
			CompartmentId: &v.CompartmentID,
			IsEnabled:     &isenabled,
			VcnId:         &vcnID,
			DisplayName:   &displayName,
		},
	}
	resp, err := v.Client.CreateInternetGateway(v.Ctx, req)
	if err != nil {
		return core.CreateInternetGatewayResponse{}, err
	}
	return resp, nil
}

type SecurityList struct {
	VcnID       string
	EgressRule  []core.EgressSecurityRule
	IngressRule []core.IngressSecurityRule
}

func (v Vcn) addSecurityList(sl SecurityList) (core.CreateSecurityListResponse, error) {
	req := core.CreateSecurityListRequest{
		CreateSecurityListDetails: core.CreateSecurityListDetails{
			CompartmentId:        &v.CompartmentID,
			VcnId:                &sl.VcnID,
			EgressSecurityRules:  sl.EgressRule,
			IngressSecurityRules: sl.IngressRule,
		},
	}
	resp, err := v.Client.CreateSecurityList(v.Ctx, req)
	if err != nil {
		return core.CreateSecurityListResponse{}, nil
	}
	return resp, nil
}

//AddSubnet will create a subnet for a given VCN
func (v Vcn) AddSubnet(vcnID string, subnet Network) (core.CreateSubnetResponse, error) {

	req := core.CreateSubnetRequest{
		CreateSubnetDetails: core.CreateSubnetDetails{
			CompartmentId: &v.CompartmentID,
			CidrBlock:     &subnet.CIDR,
			VcnId:         &vcnID,
			DisplayName:   &subnet.Name,
		},
	}

	resp, err := v.Client.CreateSubnet(v.Ctx, req)
	if err != nil {
		return core.CreateSubnetResponse{}, err
	}
	return resp, nil
}

//QuickNetworking will deploy a default vcn to be used the an OKE cluster
//The QuickVCN method will deploy
//* Virtual Cloud Network (VCN)
//* Internet Gateway (IG)
//* NAT Gateway (NAT)
//* Service Gateway (SGW)
func QuickNetworking(vcn Vcn) error {

	//VCN First
	random := util.RandomInt(6)
	net := Network{
		Name:          defaultName + "_vcn_" + random,
		CIDR:          defaultCIDR,
		CompartmentID: vcn.CompartmentID,
	}
	vcnresp, err := vcn.Create(net)
	if err != nil {
		return err
	}

	var quickSubnets []Network

	nodeSubnet := Network{
		Name:          defaultName + "_subnet_workers_" + random,
		CIDR:          "10.0.10.0/24",
		CompartmentID: vcn.CompartmentID,
	}

	quickSubnets = append(quickSubnets, nodeSubnet)

	svclbSubnet := Network{
		Name:          defaultName + "_subnet_svclb_" + random,
		CIDR:          "10.0.20.0/24",
		CompartmentID: vcn.CompartmentID,
	}

	quickSubnets = append(quickSubnets, svclbSubnet)

	k8sApiSubnet := Network{
		Name:          defaultName + "_subnet_k8sApiEndpoint_" + random,
		CIDR:          "10.0.0.0/28",
		CompartmentID: vcn.CompartmentID,
	}

	quickSubnets = append(quickSubnets, k8sApiSubnet)

	//Create three subnets for worker nodes, public services, and k8s api endpoint
	for _, subnet := range quickSubnets {
		_, err = vcn.AddSubnet(*vcnresp.Id, subnet)
		if err != nil {
			return err
		}
	}

	//Create Internet Gateway
	igwName := defaultName + "_igw_" + random
	_, err = vcn.AddInternetGateway(*vcnresp.Id, igwName)
	if err != nil {
		return err
	}

	return errors.New("This is an error")
}
