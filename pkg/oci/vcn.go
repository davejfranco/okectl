package oci

import (
	"context"
	"errors"

	"github.com/davejfranco/okectl/pkg/util"
	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/core"
)

var (
	defaultVcnName        = "okectl_quick_vcn_" + util.RandomInt(6)
	defaultCIDR    string = "10.0.0.0/16"
)

//Vcn to be used con virtual network ops
type Vcn struct {
	Client *core.VirtualNetworkClient
	Ctx    context.Context
}

//getAllVcns returns all vcn created on a given CompartmentID
/* func (v *Vcn) getAllVcns(compartmentID string) (response []core.Vcn, err error) {

	lvcn := core.ListVcnsRequest{CompartmentId: common.String(compartmentID)}
	req, err := v.Client.ListVcns(v.Ctx, lvcn)
	if err != nil {
		return req.Items, err
	}
	return req.Items, nil
} */

//DescribeVcn returns vcn details of a given vcnID
func (v *Vcn) DescribeVcn(vcnID string) (core.Vcn, error) {

	req := core.GetVcnRequest{VcnId: common.String(vcnID)}
	vcnDetail, err := v.Client.GetVcn(v.Ctx, req)
	if err != nil {
		return core.Vcn{}, err
	}
	return vcnDetail.Vcn, nil

}

//QuickNetworking will deploy a default vcn to be used the an OKE cluster
//The QuickVCN method will deploy
//* Virtual Cloud Network (VCN)
//* Internet Gateway (IG)
//* NAT Gateway (NAT)
//* Service Gateway (SGW)
func (v *Vcn) QuickNetworking(compartmendID string, ctx context.Context) error {

	return errors.New("This is an error")
}

//CreateVCN will deploy a vcn in a given compartmentID
func (v *Vcn) CreateVCN(compartmendID string, ctx context.Context) (core.CreateVcnResponse, error) {

	req := core.CreateVcnRequest{CreateVcnDetails: core.CreateVcnDetails{
		CompartmentId: &compartmendID,
		CidrBlock:     &defaultCIDR,
		DisplayName:   &defaultVcnName,
	}}

	resp, err := v.Client.CreateVcn(ctx, req)
	if err != nil {
		return core.CreateVcnResponse{}, err
	}

	return resp, nil
}

func (v *Vcn) addRouteTable(compartmentID, displayName, vcnID string, routes []core.RouteRule, ctx context.Context) (core.CreateRouteTableResponse, error) {
	req := core.CreateRouteTableRequest{
		CreateRouteTableDetails: core.CreateRouteTableDetails{
			CompartmentId: &compartmentID,
			RouteRules:    routes,
			VcnId:         &vcnID,
			DisplayName:   &displayName,
		},
	}
	resp, err := v.Client.CreateRouteTable(ctx, req)
	if err != nil {
		return core.CreateRouteTableResponse{}, err
	}
	return resp, nil
}

//AddSubnet will create a subnet for a given VCN
func (v *Vcn) AddSubnet(compartmentID, cidr, displayName, vcnID string, ctx context.Context) (core.CreateSubnetResponse, error) {

	req := core.CreateSubnetRequest{
		CreateSubnetDetails: core.CreateSubnetDetails{
			CompartmentId: &compartmentID,
			CidrBlock:     &cidr,
			VcnId:         &vcnID,
			DisplayName:   &displayName,
		},
	}

	resp, err := v.Client.CreateSubnet(ctx, req)
	if err != nil {
		return core.CreateSubnetResponse{}, err
	}
	return resp, nil
}

//AddInternetGateway will create a IG to allow access from the internet
func (v *Vcn) AddInternetGateway(compartmentID, vcnID, displayName string, ctx context.Context) (core.CreateInternetGatewayResponse, error) {

	isenabled := true
	req := core.CreateInternetGatewayRequest{
		CreateInternetGatewayDetails: core.CreateInternetGatewayDetails{
			CompartmentId: &compartmentID,
			IsEnabled:     &isenabled,
			VcnId:         &vcnID,
			DisplayName:   &displayName,
		},
	}
	resp, err := v.Client.CreateInternetGateway(ctx, req)
	if err != nil {
		return core.CreateInternetGatewayResponse{}, err
	}
	return resp, nil
}

func (v *Vcn) addSecurityList(compartmentID, vcnID string, egressRules []core.EgressSecurityRule, ingressRules []core.IngressSecurityRule, ctx context.Context) (core.CreateSecurityListResponse, error) {
	req := core.CreateSecurityListRequest{
		CreateSecurityListDetails: core.CreateSecurityListDetails{
			CompartmentId:        &compartmentID,
			VcnId:                &vcnID,
			EgressSecurityRules:  egressRules,
			IngressSecurityRules: ingressRules,
		},
	}
	resp, err := v.Client.CreateSecurityList(ctx, req)
	if err != nil {
		return core.CreateSecurityListResponse{}, nil
	}
	return resp, nil
}
