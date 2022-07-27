package tar

import (
	"fmt"
	"os/exec"

	"github.com/ecnepsnai/logtic"
	"github.com/ecnepsnai/pukcab"
)

type TarConfig struct {
	TarPath     string   `json:"tar_path"`
	TarballName string   `json:"tarball_name"`
	Sources     []string `json:"sources"`
}

var log = logtic.Log.Connect("pukcab/tar")

const Name = "tar"

// TarModule the Tar pukcab module
type TarModule struct{}

func (m TarModule) Name() string {
	return Name
}

func (m TarModule) Run(c interface{}) ([]pukcab.File, error) {
	config := TarConfig{
		TarPath: "tar",
	}
	if err := pukcab.MarshallConfig(c, &config); err != nil {
		return nil, fmt.Errorf("invalid config for module")
	}

	args := []string{
		"-czf",
		pukcab.GetFilePath(Name, config.TarballName),
	}
	args = append(args, config.Sources...)
	cmd := exec.Command(config.TarPath, args...)
	out, err := cmd.CombinedOutput()
	log.Debug("tar output: %s", out)
	if err != nil {
		log.Error("Error running tar command: error='%s'", err.Error())
		return nil, err
	}

	return []pukcab.File{
		{
			Path: pukcab.GetFilePath(Name, config.TarballName),
		},
	}, nil
}
