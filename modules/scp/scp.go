package scp

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/ecnepsnai/pukcab"
)

func runBackup(config SCPConfig) (*pukcab.File, error) {
	priv, err := os.CreateTemp("", "scp_prk")
	if err != nil {
		log.PPanic("Error making temp file", map[string]interface{}{
			"error": err.Error(),
		})
	}
	if err := priv.Chmod(0600); err != nil {
		log.PPanic("Error chmod temp file", map[string]interface{}{
			"error": err.Error(),
		})
	}
	privPath := priv.Name()
	pub, err := os.CreateTemp("", "scp_pbk")
	if err != nil {
		log.PPanic("Error making temp file", map[string]interface{}{
			"error": err.Error(),
		})
	}
	if err := pub.Chmod(0600); err != nil {
		log.PPanic("Error chmod temp file", map[string]interface{}{
			"error": err.Error(),
		})
	}
	pubPath := pub.Name()

	defer func() {
		priv.Close()
		pub.Close()
		os.Remove(privPath)
		os.Remove(pubPath)
	}()

	if _, err := priv.WriteString(config.PrivateKey); err != nil {
		log.PError("Error writing private key to temporary file", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	if _, err := pub.WriteString(fmt.Sprintf("%s %s\n", config.HostAddress, config.HostPublicKey)); err != nil {
		log.PError("Error writing public key to temporary file", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	scpPath := config.ScpPath
	if scpPath == "" {
		p, err := exec.LookPath("scp")
		if err != nil {
			return nil, fmt.Errorf("no scp bin found")
		}
		scpPath = p
	}

	port := config.Port
	if port == 0 {
		port = 22
	}

	outputFilePath := pukcab.GetFilePath(Name, fmt.Sprintf("%s_%s", sanitizePath(config.HostAddress), sanitizePath(config.FilePath)))

	args := []string{
		"-P", fmt.Sprintf("%d", port),
		"-o", fmt.Sprintf("UserKnownHostsFile=%s", pubPath),
		"-i", privPath,
		fmt.Sprintf("%s@%s:%s", config.Username, config.HostAddress, config.FilePath),
		outputFilePath,
	}

	log.Debug("exec %s %s", scpPath, args)
	cmd := exec.Command(scpPath, args...)
	output, err := cmd.CombinedOutput()
	log.Debug("scp output: %s", output)
	if err != nil {
		return nil, err
	}

	log.PInfo("SCP success", map[string]interface{}{
		"host_address": config.HostAddress,
		"username":     config.Username,
		"file_path":    config.FilePath,
	})

	return &pukcab.File{
		Path: outputFilePath,
	}, nil
}

func sanitizePath(fileName string) string {
	if fileName == "" {
		return fileName
	}

	p := fileName

	// Need to remove all NTFS characters,
	// as well as a couple more annoynances
	naughty := map[string]string{
		"<":    "",
		">":    "",
		":":    "",
		"\"":   "",
		"/":    "",
		"\\":   "",
		"|":    "",
		"?":    "",
		"*":    "",
		" ":    "_",
		",":    "",
		"#":    "",
		"\000": "",
	}

	for bad, good := range naughty {
		p = strings.Replace(p, bad, good, -1)
	}

	// Don't allow UNIX "hidden" files
	if p[0] == '.' {
		p = "_" + p
	}

	return p
}
