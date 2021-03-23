package pukcab

// Config describes a configuration object for pukcab
type Config struct {
	Modules           []ModuleType `json:"modules"`
	OutputDir         string       `json:"output_dir"`
	Verbose           bool         `json:"verbose"`
	ArtifactRetention int          `json:"artifact_retention"`
}

// ModuleType describes a module configuration for pukcab
type ModuleType struct {
	Name   string      `json:"name"`
	Config interface{} `json:"config"`
}
