package vcn

import (
	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/core"
)

//Conn struct to specify config and context for oracle's services
type Conn struct {
	Config common.ConfigurationProvider
}

//Client creates a connection to Oracle virtual network
func (c Conn) Client() (core.VirtualNetworkClient, error) {

	client, err := core.NewVirtualNetworkClientWithConfigurationProvider(c.Config)
	if err != nil {
		return core.VirtualNetworkClient{}, err
	}
	return client, nil
}
