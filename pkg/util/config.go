package util

import (
	"os"

	"github.com/oracle/oci-go-sdk/common"
)

//OciUser required connection arguments
type OciUser struct {
	Tenancy        string
	User           string
	Region         string
	Fingerprint    string
	privateKey     string
	KeyPassphrasse string
}

func (ouser OciUser) isEmpty() bool {
	if (OciUser{}) == ouser {
		return true
	}
	return false
}

//Config client to connect to tenant
type Config struct {
	Path    string
	Profile string
	OciUser
}

func (cfg Config) validPath() bool {

	if _, err := os.Stat(cfg.Path); err == nil {
		return true
	}
	return false
}

//Load creates a condig object to connect to OCI
func (cfg Config) Load() (common.ConfigurationProvider, error) {

	switch {
	case cfg.Path != "" && cfg.Profile != "":
		c, err := common.ConfigurationProviderFromFileWithProfile(
			cfg.Path,
			cfg.Profile,
			cfg.KeyPassphrasse,
		)
		if err != nil {
			return nil, err
		}
		return c, nil
	case cfg.Path != "" && cfg.Profile == "":
		c, err := common.ConfigurationProviderFromFile(cfg.Path, cfg.KeyPassphrasse)
		if err != nil {
			return nil, err
		}
		return c, nil
	default:
		return common.DefaultConfigProvider(), nil

	}
}
