package cloudflare

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/ecnepsnai/pukcab"
)

type cloudflareZone struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func getZones(creds CloudflareConfig) ([]cloudflareZone, error) {
	request, err := http.NewRequest("GET", "https://api.cloudflare.com/client/v4/zones", nil)
	if err != nil {
		log.Error("Error forming zones request: error='%s'", err.Error())
		return nil, err
	}
	request.Header.Add("X-Auth-Key", creds.APIKey)
	request.Header.Add("X-Auth-Email", creds.Email)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Error("Error making zone list request: error='%s'", err.Error())
		return nil, err
	}
	if response.StatusCode != 200 {
		log.Error("Error making zone list request: error='http %d'", response.StatusCode)
		return nil, fmt.Errorf("http %d", response.StatusCode)
	}

	type zoneResponse struct {
		Result []cloudflareZone `json:"result"`
	}
	result := zoneResponse{}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		log.Error("Error decoding zone list request: error='%s'", err.Error())
		return nil, err
	}

	return result.Result, nil
}

func downloadZoneFile(creds CloudflareConfig, zone cloudflareZone) (*pukcab.File, error) {
	request, err := http.NewRequest("GET", "https://api.cloudflare.com/client/v4/zones/"+zone.ID+"/dns_records/export", nil)
	if err != nil {
		log.Error("Error forming zones request: error='%s'", err.Error())
		return nil, err
	}
	request.Header.Add("X-Auth-Key", creds.APIKey)
	request.Header.Add("X-Auth-Email", creds.Email)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Error("Error making zone list request: error='%s'", err.Error())
		return nil, err
	}
	if response.StatusCode != 200 {
		log.Error("Error making zone list request: error='http %d'", response.StatusCode)
		return nil, fmt.Errorf("http %d", response.StatusCode)
	}

	filePath := pukcab.GetFilePath(Name, zone.Name+".txt")
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Error("Error opening destination file: file_path='%s' error='%s'", filePath, err.Error())
		return nil, err
	}
	defer f.Close()

	if _, err := io.Copy(f, response.Body); err != nil {
		log.Error("Error writing backup file: file_path='%s' error='%s'", filePath, err.Error())
		return nil, err
	}

	return &pukcab.File{
		Path: filePath,
	}, nil
}
