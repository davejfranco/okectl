package oke

import (
	"context"

	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/core"
)

//Vcn to be used con virtual network ops
type Vcn struct {
	Client *core.VirtualNetworkClient
}

//getAllVcns returns all vcn created on a given CompartmentID
func (v Vcn) getAllVcns(ctx context.Context, compartmentID string) (response []core.Vcn, err error) {

	lvcn := core.ListVcnsRequest{CompartmentId: common.String(compartmentID)}
	req, err := v.Client.ListVcns(ctx, lvcn)
	if err != nil {
		return req.Items, err
	}
	return req.Items, nil
}

//DescribeVcn returns vcn details of a given vcnID
func (v Vcn) DescribeVcn(ctx context.Context, vcnID string) (core.Vcn, error) {

	req := core.GetVcnRequest{VcnId: common.String(vcnID)}
	vcnDetail, err := v.Client.GetVcn(ctx, req)
	if err != nil {
		return core.Vcn{}, err
	}
	return vcnDetail.Vcn, nil

}
