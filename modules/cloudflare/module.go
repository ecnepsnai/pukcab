package cloudflare

import (
	"fmt"

	"github.com/ecnepsnai/logtic"
	"github.com/ecnepsnai/pukcab"
)

var log = logtic.Log.Connect("pukcab/cloudflare")

const Name = "cloudflare"

type CloudflareConfig struct {
	Email  string `json:"cloudflare_email"`
	APIKey string `json:"cloudflare_api_key"`
}

// CloudflareModule the Cloudflare pukcab module
type CloudflareModule struct{}

func (m CloudflareModule) Name() string {
	return Name
}

func (m CloudflareModule) Run(c interface{}) ([]pukcab.File, error) {
	config := CloudflareConfig{}
	if err := pukcab.MarshallConfig(c, &config); err != nil {
		return nil, fmt.Errorf("invalid config for module")
	}

	zones, err := getZones(config)
	if err != nil {
		return nil, err
	}

	files := []pukcab.File{}
	for _, zone := range zones {
		file, err := downloadZoneFile(config, zone)
		if err != nil {
			return nil, err
		}
		files = append(files, *file)
	}

	return files, nil
}
