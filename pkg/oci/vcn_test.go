package oci

import (
	"context"
	"net/http"
	"testing"

	"github.com/oracle/oci-go-sdk/core"
)

var (
	mockVcnID             string = "ocid1.vcn.oc1.ash.abuw4ljrlsfiqw6vzzxb43vyypt4pkodawglp3wqxjqofakrwvou52gb6s5a"
	mockCompartmentID     string = "ocid1.compartment.oc1..ash.aaaaaajrlsfiqw6vzzxb43vyypt4pkodawglp3wqxjqofakrwvou52gb6s5a"
	mockRouteTableID      string = "ocid1.routetable.oc1.iad.aaaaabbb4csskd4reftq2womsopvwq2t5e2bl4fkskkl5mhyqzseywi62faa"
	mockInternetGatewayID string = "ocid1.vcn.oc1.iad.azaaaaaazluk3aaa43uodzm7daxtxkq3vaq6ml5aj6mq2vagyucavl4konra"
	mockSecurityListID    string = "ocid1.securitylist.oc1.iad.bbbaaaaac3yhu7fadx6rt5eijbrg3fdjttzcouwuvdeuwlsg6kl2w3dafafa"
	mockSubnetID          string = "ocid1.subnet.oc1.iad.bbbaaaaa3lqwtfx4nvtn4hecezvshvqfvp7bnk5xmooqze5pzsqqjjyx2izq"
)

type mockVcnClient struct {
}

var httpResponseOK http.Response = http.Response{
	Status:     "200 OK",
	StatusCode: 200,
	Proto:      "HTTP/1.1",
	ProtoMajor: 1,
	ProtoMinor: 1,
}

func (m mockVcnClient) CreateVcn(ctx context.Context, req core.CreateVcnRequest) (core.CreateVcnResponse, error) {

	return core.CreateVcnResponse{
		RawResponse: &httpResponseOK,
		Vcn: core.Vcn{
			CidrBlock:      req.CidrBlock,
			CompartmentId:  req.CompartmentId,
			Id:             &mockVcnID,
			LifecycleState: core.VcnLifecycleStateAvailable,
		},
	}, nil
}

func (m mockVcnClient) CreateRouteTable(ctx context.Context, req core.CreateRouteTableRequest) (core.CreateRouteTableResponse, error) {
	return core.CreateRouteTableResponse{
		RawResponse: &httpResponseOK,
		RouteTable: core.RouteTable{
			CompartmentId:  req.CompartmentId,
			Id:             &mockRouteTableID,
			LifecycleState: core.RouteTableLifecycleStateAvailable,
			RouteRules:     req.RouteRules,
		},
	}, nil
}

func (m mockVcnClient) CreateInternetGateway(ctx context.Context, req core.CreateInternetGatewayRequest) (core.CreateInternetGatewayResponse, error) {
	return core.CreateInternetGatewayResponse{
		RawResponse: &httpResponseOK,
		InternetGateway: core.InternetGateway{
			CompartmentId:  req.CompartmentId,
			Id:             &mockInternetGatewayID,
			LifecycleState: core.InternetGatewayLifecycleStateAvailable,
			VcnId:          &mockVcnID,
		},
	}, nil
}

func (m mockVcnClient) CreateSecurityList(ctx context.Context, req core.CreateSecurityListRequest) (core.CreateSecurityListResponse, error) {
	return core.CreateSecurityListResponse{
		RawResponse: &httpResponseOK,
		SecurityList: core.SecurityList{
			CompartmentId:        req.CompartmentId,
			DisplayName:          req.DisplayName,
			EgressSecurityRules:  req.EgressSecurityRules,
			IngressSecurityRules: req.IngressSecurityRules,
			Id:                   &mockSecurityListID,
			LifecycleState:       core.SecurityListLifecycleStateAvailable,
		},
	}, nil
}

func (m mockVcnClient) CreateSubnet(ctx context.Context, req core.CreateSubnetRequest) (core.CreateSubnetResponse, error) {
	return core.CreateSubnetResponse{
		RawResponse: &httpResponseOK,
		Subnet: core.Subnet{
			CidrBlock:      req.CidrBlock,
			CompartmentId:  req.CompartmentId,
			Id:             &mockSubnetID,
			RouteTableId:   req.RouteTableId,
			VcnId:          req.VcnId,
			LifecycleState: core.SubnetLifecycleStateAvailable,
		},
	}, nil
}
func TestCreateVCN(t *testing.T) {

	var m mockVcnClient
	net := Network{
		Name:          "MockVCN",
		CIDR:          "172.16.0.0/16",
		CompartmentID: mockCompartmentID,
	}

	_, err := net.CreateVCN(m)
	if err != nil {
		t.Error(err)
	}

}
