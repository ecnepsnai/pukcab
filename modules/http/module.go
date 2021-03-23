package http

import (
	"crypto/tls"
	"fmt"
	"io"
	nhttp "net/http"
	"os"

	"github.com/ecnepsnai/logtic"
	"github.com/ecnepsnai/pukcab"
)

var log = logtic.Connect("pukcab/http")

type HTTPConfig struct {
	URL                        string            `json:"url"`
	FileName                   string            `json:"file_name"`
	AllowUntrustedCertificates bool              `json:"allow_untrusted_certificates"`
	Headers                    map[string]string `json:"headers"`
}

// HTTPModule the HTTP pukcab module
type HTTPModule struct{}

func (m HTTPModule) Name() string {
	return "http"
}

func (m HTTPModule) Run(c interface{}) ([]pukcab.File, error) {
	config := HTTPConfig{}
	if err := pukcab.MarshallConfig(c, &config); err != nil {
		return nil, fmt.Errorf("invalid config for module")
	}

	filePath := pukcab.GetFilePath("http", config.FileName)
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Error("Error opening destination file for writing: file_path='%s' error='%s'", filePath, err.Error())
		return nil, err
	}
	defer f.Close()
	request, err := nhttp.NewRequest("GET", config.URL, nil)
	if err != nil {
		log.Error("Error forming HTTP request: url='%s' error='%s'", config.URL, err.Error())
		return nil, err
	}
	for k, v := range config.Headers {
		request.Header.Add(k, v)
	}

	client := &nhttp.Client{}
	tr := &nhttp.Transport{}
	if config.AllowUntrustedCertificates {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	client.Transport = tr

	resp, err := client.Do(request)
	if err != nil {
		log.Error("Error making HTTP request: url='%s' error='%s'", config.URL, err.Error())
		return nil, err
	}
	if resp.StatusCode != 200 {
		log.Error("Error making HTTP request: url='%s' error='http %d'", config.URL, resp.StatusCode)
		return nil, fmt.Errorf("http %d", resp.StatusCode)
	}

	if _, err := io.Copy(f, resp.Body); err != nil {
		log.Error("Error writing to destination file: file_path='%s' error='%s'", filePath, err.Error())
		return nil, err
	}

	return []pukcab.File{
		{
			Path: filePath,
		},
	}, nil
}
