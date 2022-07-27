package scp

import (
	"fmt"

	"github.com/ecnepsnai/logtic"
	"github.com/ecnepsnai/pukcab"
)

var log = logtic.Log.Connect("pukcab/scp")

const Name = "scp"

type SCPConfig struct {
	HostAddress   string `json:"host_address"`
	Port          uint16 `json:"port"`
	Username      string `json:"username"`
	PrivateKey    string `json:"private_key"`
	HostPublicKey string `json:"host_public_key"`
	FilePath      string `json:"file_path"`
	ScpPath       string `json:"scp_path"`
}

// SCPModule the SCP pukcab module
type SCPModule struct{}

func (m SCPModule) Name() string {
	return Name
}

func (m SCPModule) Run(c interface{}) ([]pukcab.File, error) {
	config := SCPConfig{}
	if err := pukcab.MarshallConfig(c, &config); err != nil {
		return nil, fmt.Errorf("invalid config for module")
	}

	file, err := runBackup(config)
	if err != nil {
		return nil, err
	}

	return []pukcab.File{*file}, nil
}
