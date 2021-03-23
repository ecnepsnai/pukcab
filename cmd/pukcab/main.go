package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ecnepsnai/logtic"
	"github.com/ecnepsnai/pukcab"
	"github.com/ecnepsnai/pukcab/modules/cloudflare"
	httpModule "github.com/ecnepsnai/pukcab/modules/http"
	"github.com/ecnepsnai/pukcab/modules/pfsense"
	"github.com/ecnepsnai/pukcab/modules/tar"
)

var pukcabModules = []pukcab.Module{
	pfsense.PFSenseModule{},
	cloudflare.CloudflareModule{},
	tar.TarModule{},
	httpModule.HTTPModule{},
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <Path to config JSON>\n", os.Args[0])
		os.Exit(1)
	}

	configFilePath := os.Args[1]
	f, err := os.OpenFile(configFilePath, os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	config := pukcab.Config{}
	if err := json.NewDecoder(f).Decode(&config); err != nil {
		panic(err)
	}

	if config.Verbose {
		logtic.Log.Level = logtic.LevelDebug
	}
	logtic.Open()
	pukcab.Configure(config)

	moduleMap := map[string]pukcab.Module{}
	for _, module := range pukcabModules {
		moduleMap[module.Name()] = module
	}

	for _, module := range config.Modules {
		if _, ok := moduleMap[module.Name]; !ok {
			fmt.Fprintf(os.Stderr, "No module named '%s'\n", module.Name)
			os.Exit(1)
		}
	}

	for _, module := range config.Modules {
		pukcab.RunModule(moduleMap[module.Name], module.Config)
		pukcab.CleanupModule(moduleMap[module.Name])
	}
}
