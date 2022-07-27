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

var log = logtic.Log.Connect("pukcab")

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
	log.PInfo("Starting module", map[string]interface{}{
		"module_name": name,
	})
	start := time.Now()
	makeDirectoryIfNotExists(path.Join(pukcabConfig.OutputDir, name, time.Now().Format("2006-01-02")))
	files, err := module.Run(config)
	if err != nil {
		log.PError("Error running module", map[string]interface{}{
			"module_name": name,
			"error":       err.Error(),
		})
	}
	nFiles := 0
	for _, file := range files {
		info, err := os.Stat(file.Path)
		if err != nil {
			log.PError("Unable to stat module artifact", map[string]interface{}{
				"module_name": name,
				"file_path":   file.Path,
				"error":       err.Error(),
			})
			os.Remove(file.Path)
			continue
		}
		if info.Size() == 0 {
			log.PError("Module produced empty artifact", map[string]interface{}{
				"module_name": name,
				"file_path":   file.Path,
			})
			os.Remove(file.Path)
			continue
		}

		log.PInfo("Backup artifact saved", map[string]interface{}{
			"module_name": name,
			"file_path":   file.Path,
			"size":        logtic.FormatBytesB(uint64(info.Size())),
		})
		nFiles++
	}
	log.PInfo("Module finished", map[string]interface{}{
		"module_name": name,
		"n_files":     nFiles,
		"duration":    time.Since(start).String(),
	})
	return nil
}

// CleanupModule remove expired artifacts
func CleanupModule(module Module) error {
	if pukcabConfig.ArtifactRetention <= 0 {
		return nil
	}

	name := module.Name()
	log.PInfo("Starting module cleanup", map[string]interface{}{
		"module_name": name,
	})
	start := time.Now()

	moduleOutputPath := path.Join(pukcabConfig.OutputDir, name)
	items, err := os.ReadDir(moduleOutputPath)
	if err != nil {
		log.PError("Error reading directory", map[string]interface{}{
			"module_name": name,
			"directory":   moduleOutputPath,
		})
		return err
	}

	retentionHours := float64(pukcabConfig.ArtifactRetention) * 24.0
	datePattern := regexp.MustCompile("[0-9]{4}-[0-9]{2}-[0-9]{2}")

	for _, item := range items {
		if !item.IsDir() {
			continue
		}
		itemPath := path.Join(moduleOutputPath, item.Name())

		subItems, _ := os.ReadDir(itemPath)
		for _, subItem := range subItems {
			info, _ := subItem.Info()
			subItemPath := path.Join(itemPath, subItem.Name())
			if info.Size() == 0 {
				log.PWarn("Removing empty artifact", map[string]interface{}{
					"module": name,
					"path":   subItemPath,
				})
				os.Remove(subItemPath)
			}
		}
		subItems, _ = os.ReadDir(itemPath)
		if len(subItems) == 0 {
			log.PWarn("Empty artifact directory", map[string]interface{}{
				"module": name,
				"path":   itemPath,
			})
			if err := os.RemoveAll(itemPath); err != nil {
				log.Error("Error removing expired artifact: module='%s' path='%s' error='%s'", name, itemPath, err.Error())
			}
			continue
		}

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
