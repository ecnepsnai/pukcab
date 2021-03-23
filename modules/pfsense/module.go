package pfsense

import (
	"fmt"

	"github.com/ecnepsnai/logtic"
	"github.com/ecnepsnai/pukcab"
)

var log = logtic.Connect("pukcab/pfsense")

type PFSenseConfig struct {
	HostAddress                string `json:"host_address"`
	Username                   string `json:"username"`
	Password                   string `json:"password"`
	AllowUntrustedCertificates bool   `json:"allow_untrusted_certificates"`
	EncryptPassword            string `json:"encrypt_password"`
}

// PFSenseModule the PFSense pukcab module
type PFSenseModule struct{}

func (m PFSenseModule) Name() string {
	return "pfsense"
}

func (m PFSenseModule) Run(c interface{}) ([]pukcab.File, error) {
	config := PFSenseConfig{}
	if err := pukcab.MarshallConfig(c, &config); err != nil {
		return nil, fmt.Errorf("invalid config for module")
	}

	file, err := runBackup(config)
	if err != nil {
		return nil, err
	}

	return []pukcab.File{*file}, nil
}
