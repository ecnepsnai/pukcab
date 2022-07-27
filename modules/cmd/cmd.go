package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/ecnepsnai/logtic"
	"github.com/ecnepsnai/pukcab"
)

var log = logtic.Log.Connect("pukcab/cmd")

const Name = "cmd"

type CmdConfig struct {
	ExecPath      string   `json:"exec_path"`
	Args          []string `json:"args"`
	Env           []string `json:"env"`
	Wd            string   `json:"wd"`
	OutputName    string   `json:"output_name"`
	IncludeStderr bool     `json:"include_stderr"`
}

// CmdModule the Cmd pukcab module
type CmdModule struct{}

func (m CmdModule) Name() string {
	return Name
}

func (m CmdModule) Run(c interface{}) ([]pukcab.File, error) {
	config := CmdConfig{}
	if err := pukcab.MarshallConfig(c, &config); err != nil {
		return nil, fmt.Errorf("invalid config for module")
	}

	outputFile := pukcab.GetFilePath(Name, config.OutputName)
	f, err := os.OpenFile(outputFile, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		log.PPanic("Error opening output file", map[string]interface{}{
			"error": err.Error(),
		})
	}
	defer f.Close()

	command := exec.Command(config.ExecPath, config.Args...)
	command.Env = os.Environ()
	if config.Env != nil {
		command.Env = append(command.Env, config.Env...)
	}
	if config.Wd != "" {
		command.Dir = config.Wd
	}
	command.Stdout = f
	if config.IncludeStderr {
		command.Stderr = f
	} else {
		command.Stderr = os.Stderr
	}
	if err := command.Run(); err != nil {
		log.PError("Error running command", map[string]interface{}{
			"exec":  config.ExecPath,
			"args":  config.Args,
			"error": err.Error(),
		})
		return nil, err
	}

	return []pukcab.File{{Path: outputFile}}, nil
}
