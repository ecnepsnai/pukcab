// Package pukcab provides a common interface for a modular backup system.
package pukcab

import (
	"encoding/json"
	"os"
	"path"
	"regexp"
	"time"

	"github.com/ecnepsnai/logtic"
)

var log = logtic.Connect("pukcab")

var pukcabConfig *Config

// Configure will prepare pukcab for running with the given config instance
func Configure(config Config) {
	if pukcabConfig != nil {
		panic("pukcab already configred")
	}

	pukcabConfig = &config
	makeDirectoryIfNotExists(config.OutputDir)
}

// Module describes the interface for a pukcab module
type Module interface {
	Name() string
	Run(c interface{}) ([]File, error)
}

// File describes a backed-up file
type File struct {
	Path string
}

// RunModule will run the given backup module
func RunModule(module Module, config interface{}) error {
	name := module.Name()
	log.Info("Starting module: module_name='%s'", name)
	start := time.Now()
	makeDirectoryIfNotExists(path.Join(pukcabConfig.OutputDir, name, time.Now().Format("2006-01-02")))
	files, err := module.Run(config)
	if err != nil {
		log.Error("Error running module: module_name='%s' error='%s'", name, err.Error())
		return err
	}
	for _, file := range files {
		log.Info("Backup artifact saved: module_name='%s' file_path='%s'", name, file.Path)
	}
	log.Info("Module finished: module_name='%s' number_files=%d duration_s=%f", name, len(files), time.Since(start).Seconds())
	return nil
}

// CleanupModule remove expired artifacts
func CleanupModule(module Module) error {
	if pukcabConfig.ArtifactRetention <= 0 {
		return nil
	}

	name := module.Name()
	log.Info("Starting module cleanup: module_name='%s'", name)
	start := time.Now()

	moduleOutputPath := path.Join(pukcabConfig.OutputDir, name)
	items, err := os.ReadDir(moduleOutputPath)
	if err != nil {
		log.Error("Error reading directory: module_name='%s' directory='%s'", name, moduleOutputPath)
		return err
	}

	retentionHours := float64(pukcabConfig.ArtifactRetention) * 24.0
	datePattern := regexp.MustCompile("[0-9]{4}-[0-9]{2}-[0-9]{2}")

	for _, item := range items {
		if !item.IsDir() {
			continue
		}
		itemPath := path.Join(moduleOutputPath, item.Name())

		dateStr := datePattern.FindString(item.Name())
		if dateStr == "" {
			continue
		}
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			continue
		}
		if time.Since(date).Hours() <= retentionHours {
			log.Debug("Artifact not expired: module='%s' path='%s'", name, itemPath)
			continue
		}
		log.Warn("Artifact expired: module='%s' path='%s'", name, itemPath)
		if err := os.RemoveAll(itemPath); err != nil {
			log.Error("Error removing expired artifact: module='%s' path='%s' error='%s'", name, itemPath, err.Error())
		}
	}

	log.Info("Module cleanup finished: module_name='%s' duration_s=%f", name, time.Since(start).Seconds())
	return nil
}

// GetFilePath return an absolute path for a backup artifact.
// Must specify the module name, a description (which needs to be a filename-safe string), and an extension (without a dot)
func GetFilePath(moduleName, fileName string) string {
	return path.Join(pukcabConfig.OutputDir, moduleName, time.Now().Format("2006-01-02"), fileName)
}

func MarshallConfig(in interface{}, out interface{}) error {
	data, err := json.Marshal(in)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, out); err != nil {
		return err
	}

	return nil
}

func directoryExists(dirName string) bool {
	info, err := os.Stat(dirName)
	if os.IsNotExist(err) {
		return false
	}
	if info.IsDir() {
		return true
	}
	panic("File found at path when directory wanted")
}

func makeDirectoryIfNotExists(dirName string) error {
	if !directoryExists(dirName) {
		return os.MkdirAll(dirName, os.ModePerm)
	}

	return nil
}
