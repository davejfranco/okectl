package oci

import (
	"reflect"
	"testing"

	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/containerengine"
	"github.com/stretchr/testify/assert"
)

// test configProvider
func TestConfigProvider(t *testing.T) {

	asserts := assert.New(t)

	var cases = []Config{
		{Profile: "DEFAULT", ConfigLocation: ""},
		{Profile: "", ConfigLocation: ""},
		{Profile: "DEFAULT", ConfigLocation: "test"},
	}

	for _, c := range cases {
		config := c.configProvider()
		asserts.Equal(reflect.TypeOf(config), reflect.TypeOf(common.DefaultConfigProvider()))
	}
}

// test OKE
func TestOKE(t *testing.T) {
	asserts := assert.New(t)

	cfg := Config{
		Profile:        "DEFAULT",
		ConfigLocation: "",
	}

	client, err := cfg.Oke()

	//Should be nil
	asserts.Nil(err)

	//Check if client is of type ContainerEngineClient
	asserts.Equal(reflect.TypeOf(client), reflect.TypeOf(containerengine.ContainerEngineClient{}))

}
