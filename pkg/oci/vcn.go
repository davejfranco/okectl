package oci

import (
	"context"

	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/core"
)

const (
	//Resource default values
	defaultName = "okectl_quick"
	//Network default settings
	defaultNetworkCIDR = "10.0.0.0/16"
	defaultWorkerCIDR  = "10.0.10.0/24"
	defaultAPICIDR     = "10.0.0.0/28"
	defaultLBCIDR      = "10.0.20.0/24"
	//Network Ports
	sshport              = 22
	httpsPort            = 443
	K8sAPIPort           = 6443
	WorkerToControlPlane = 12250
)

//Global variables
var (
	apiPort core.PortRange = core.PortRange{
		Max: common.Int(K8sAPIPort),
		Min: common.Int(K8sAPIPort),
	}

	controlPlane core.PortRange = core.PortRange{
		Max: common.Int(WorkerToControlPlane),
		Min: common.Int(WorkerToControlPlane),
	}

	httpsConn core.PortRange = core.PortRange{
		Max: common.Int(httpsPort),
		Min: common.Int(httpsPort),
	}

	//Destination Unreachable
	icmp core.IcmpOptions = core.IcmpOptions{
		Type: common.Int(3),
		Code: common.Int(4),
	}
)

//VcnClient allows to create an abstraction of the core.NewVirtualNetworkClientWithConfigurationProvider
//to facilitate testing
type VcnClient interface {
	CreateVcn(context.Context, core.CreateVcnRequest) (core.CreateVcnResponse, error)
	CreateRouteTable(context.Context, core.CreateRouteTableRequest) (core.CreateRouteTableResponse, error)
	CreateInternetGateway(context.Context, core.CreateInternetGatewayRequest) (core.CreateInternetGatewayResponse, error)
	CreateSecurityList(context.Context, core.CreateSecurityListRequest) (core.CreateSecurityListResponse, error)
	CreateSubnet(context.Context, core.CreateSubnetRequest) (core.CreateSubnetResponse, error)
}

type Network struct {
	Name          string
	CIDR          string
	CompartmentID string
}

//DescribeVcn returns vcn details of a given vcnID
/*func DescribeVcn(client VcnClient, vcnID string) (core.Vcn, error) {

	req := core.GetVcnRequest{VcnId: common.String(vcnID)}
	resp, err := client.GetVcn(context.Background(), req)
	if err != nil {
		return core.Vcn{}, err
	}
	return resp.Vcn, nil
}
*/

//CreateVCN will deploy a vcn in a given compartmentID
func (n Network) CreateVCN(client VcnClient) (core.CreateVcnResponse, error) {

	req := core.CreateVcnRequest{CreateVcnDetails: core.CreateVcnDetails{
		CompartmentId: &n.CompartmentID,
		CidrBlock:     &n.CIDR,
		DisplayName:   &n.Name,
	},
	}

	resp, err := client.CreateVcn(context.Background(), req)
	if err != nil {
		return core.CreateVcnResponse{}, err
	}

	return resp, nil
}

func (n Network) AddRouteTable(client VcnClient, displayName, vcnID string, routes []core.RouteRule) (core.CreateRouteTableResponse, error) {
	req := core.CreateRouteTableRequest{
		CreateRouteTableDetails: core.CreateRouteTableDetails{
			CompartmentId: &n.CompartmentID,
			RouteRules:    routes,
			VcnId:         &vcnID,
			DisplayName:   &displayName,
		},
	}
	resp, err := client.CreateRouteTable(context.Background(), req)
	if err != nil {
		return core.CreateRouteTableResponse{}, err
	}
	return resp, nil
}

//AddInternetGateway will create a IG to allow access from the internet
func (n Network) AddInternetGateway(client VcnClient, vcnID, displayName string) (core.CreateInternetGatewayResponse, error) {

	req := core.CreateInternetGatewayRequest{
		CreateInternetGatewayDetails: core.CreateInternetGatewayDetails{
			CompartmentId: &n.CompartmentID,
			IsEnabled:     common.Bool(true),
			VcnId:         &vcnID,
			DisplayName:   &displayName,
		},
	}
	resp, err := client.CreateInternetGateway(context.Background(), req)
	if err != nil {
		return core.CreateInternetGatewayResponse{}, err
	}
	return resp, nil
}

type SecurityList struct {
	Name        string
	VcnID       string
	EgressRule  []core.EgressSecurityRule
	IngressRule []core.IngressSecurityRule
}

func ClusterAPIIngressRules(sourceCIDR string, workersCIDR []string) []core.IngressSecurityRule {
	var rules []core.IngressSecurityRule

	if sourceCIDR == "0.0.0.0/0" {
		publicAccess := core.IngressSecurityRule{
			Protocol: common.String("6"), //TCP
			Source:   common.String(sourceCIDR),
			TcpOptions: &core.TcpOptions{
				DestinationPortRange: &apiPort,
			},
			Description: common.String("External access to Kubernetes API endpoint"),
		}
		rules = append(rules, publicAccess)
	} else {
		sourceAccess := core.IngressSecurityRule{
			Protocol: common.String("6"),
			Source:   common.String(sourceCIDR),
			TcpOptions: &core.TcpOptions{
				DestinationPortRange: &apiPort,
			},
			Description: common.String("Access to Kubernetes API endpoint"),
		}
		rules = append(rules, sourceAccess)
	}

	//workers access to K8s API endpoint
	for _, cidr := range workersCIDR {
		//Grant workers nodes to the K8s API endpoint
		workerAPIaccess := core.IngressSecurityRule{
			Protocol: common.String("6"),
			Source:   common.String(cidr),
			TcpOptions: &core.TcpOptions{
				DestinationPortRange: &apiPort,
			},
			Description: common.String("Kubernetes worker to Kubernetes API endpoint communication"),
		}
		rules = append(rules, workerAPIaccess)

		//Grant workers to the control Plane
		workerToCPaccess := core.IngressSecurityRule{
			Protocol: common.String("6"),
			Source:   common.String(cidr),
			TcpOptions: &core.TcpOptions{
				DestinationPortRange: &controlPlane,
			},
			Description: common.String("Kubernetes worker to control plane communication"),
		}
		rules = append(rules, workerToCPaccess)

		//Path Discovery ICMP
		workerDiscoveryDU := core.IngressSecurityRule{
			Protocol:    common.String("1"), //ICMP
			Source:      common.String(cidr),
			IcmpOptions: &icmp,
			Description: common.String("Path discovery"),
		}
		rules = append(rules, workerDiscoveryDU)
	}

	return rules
}

func ClusterAPIEgressRules(workersCIDR []string) []core.EgressSecurityRule {
	var rules []core.EgressSecurityRule

	for _, cidr := range workersCIDR {
		//Path Discovery ICMP
		workerDiscoveryDU := core.EgressSecurityRule{
			Protocol:    common.String("1"), //ICMP
			Destination: common.String(cidr),
			IcmpOptions: &icmp,
			Description: common.String("Path discovery"),
		}
		rules = append(rules, workerDiscoveryDU)

		//Grant workers all tcp traffic
		workersAll := core.EgressSecurityRule{
			Protocol:    common.String("6"), //TCP
			Destination: common.String(cidr),
			Description: common.String("All traffic to worker nodes"),
		}
		rules = append(rules, workersAll)
	}
	//Grant egress access to OCI services
	httpsOCI := core.EgressSecurityRule{
		Protocol:        common.String("6"), //TCP
		Destination:     common.String("all-iad-services-in-oracle-services-network"),
		DestinationType: "SERVICE_CIDR_BLOCK",
		TcpOptions: &core.TcpOptions{
			DestinationPortRange: &httpsConn,
		},
		Description: common.String("All traffic to worker nodes"),
	}
	rules = append(rules, httpsOCI)
	return rules
}
func (n Network) AddSecurityList(client VcnClient, sl SecurityList) (core.CreateSecurityListResponse, error) {
	req := core.CreateSecurityListRequest{
		CreateSecurityListDetails: core.CreateSecurityListDetails{
			CompartmentId:        &n.CompartmentID,
			VcnId:                &sl.VcnID,
			DisplayName:          &sl.Name,
			EgressSecurityRules:  sl.EgressRule,
			IngressSecurityRules: sl.IngressRule,
		},
	}

	resp, err := client.CreateSecurityList(context.Background(), req)
	if err != nil {
		return core.CreateSecurityListResponse{}, err
	}
	return resp, nil
}

//AddSubnet will create a subnet for a given VCN
func (n Network) AddSubnet(client VcnClient, vcnID string, subnet Network, seclistIDs []string) (core.CreateSubnetResponse, error) {

	req := core.CreateSubnetRequest{
		CreateSubnetDetails: core.CreateSubnetDetails{
			CompartmentId:   &n.CompartmentID,
			CidrBlock:       &subnet.CIDR,
			VcnId:           &vcnID,
			DisplayName:     &subnet.Name,
			SecurityListIds: seclistIDs,
		},
	}

	resp, err := client.CreateSubnet(context.Background(), req)
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
/*
func QuickNetworking(client VcnClient, compartmentID string) error {

	//Random Naming is necessary
	random := util.RandomInt(6)
	// VCN First
	net := Network{
		Name:          defaultName + "_vcn_" + random,
		CIDR:          defaultNetworkCIDR,
		CompartmentID: compartmentID,
	}

	fmt.Println("Creating VCN...")
	resp, err := net.CreateVCN() //vcn.Create(net)
	if err != nil {
		return err
	}

	//Security list
	apiIngress := ClusterAPIIngressRules("0.0.0.0/0", []string{defaultWorkerCIDR})
	fmt.Println(apiIngress)
	apiEgress := ClusterAPIEgressRules([]string{defaultWorkerCIDR})
	fmt.Println(apiEgress)

	k8sapisl := SecurityList{
		Name:        defaultName + "_sl_k8sapi_" + random,
		VcnID:       *vcnresp.Id,
		EgressRule:  apiEgress,
		IngressRule: apiIngress,
	}

	fmt.Println("Creating Security Lists for API Endpoint...")
	slResp, err := vcn.AddSecurityList(k8sapisl)
	if err != nil {
		return err
	}
	fmt.Println(slResp)

	var quickSubnets []Network

	nodeSubnet := Network{
		Name:          defaultName + "_subnet_workers_" + random,
		CIDR:          defaultWorkerCIDR,
		CompartmentID: vcn.CompartmentID,
	}

	quickSubnets = append(quickSubnets, nodeSubnet)

	svclbSubnet := Network{
		Name:          defaultName + "_subnet_svclb_" + random,
		CIDR:          defaultLBCIDR,
		CompartmentID: vcn.CompartmentID,
	}

	quickSubnets = append(quickSubnets, svclbSubnet)

	k8sApiSubnet := Network{
		Name:          defaultName + "_subnet_k8sApiEndpoint_" + random,
		CIDR:          defaultAPICIDR,
		CompartmentID: vcn.CompartmentID,
	}

	quickSubnets = append(quickSubnets, k8sApiSubnet)

	//Create three subnets for worker nodes, public services, and k8s api endpoint
	fmt.Println("Creating Subnets...")
	for _, subnet := range quickSubnets {
		_, err = vcn.AddSubnet(*vcnresp.Id, subnet, []string{})
		if err != nil {
			return err
		}
	}

	//Create Internet Gateway
	fmt.Println("Creating Internet Gateway...")
	igwName := defaultName + "_igw_" + random
	_, err = vcn.AddInternetGateway(*vcnresp.Id, igwName)
	if err != nil {
		return err
	}

	//Lets create a route table
	return nil
}
*/
