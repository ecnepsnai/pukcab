package pfsense

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"

	"github.com/ecnepsnai/pukcab"
)

func runBackup(config PFSenseConfig) (*pukcab.File, error) {
	backupURL := "https://" + config.HostAddress + "/diag_backup.php"

	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Error("Error making new cookiejar: error='%s'", err.Error())
		return nil, err
	}
	client := &http.Client{
		Jar: jar,
	}
	tr := &http.Transport{}
	if config.AllowUntrustedCertificates {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	client.Transport = tr

	csrfRequest, err := http.NewRequest("GET", backupURL, nil)
	if err != nil {
		log.Error("Error forming CSRF request: error='%s'", err.Error())
		return nil, err
	}

	csrfResponse, err := client.Do(csrfRequest)
	if err != nil {
		log.Error("Error making CSRF request: error='%s'", err.Error())
		return nil, err
	}
	if csrfResponse.StatusCode != 200 {
		log.Error("HTTP error making CSRF request: status_code=%d", csrfResponse.StatusCode)
		return nil, fmt.Errorf("http %d", csrfResponse.StatusCode)
	}

	csrfToken, err := getCSRFTokenFromResponse(csrfResponse)
	if err != nil {
		return nil, err
	}

	loginParams := url.Values{}
	loginParams.Add("login", "Login")
	loginParams.Add("usernamefld", config.Username)
	loginParams.Add("passwordfld", config.Password)
	loginParams.Add("__csrf_magic", csrfToken)
	loginRequest, err := http.NewRequest("POST", backupURL, bytes.NewReader([]byte(loginParams.Encode())))
	if err != nil {
		log.Error("Error forming login request: error='%s'", err.Error())
		return nil, err
	}
	loginRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	loginResponse, err := client.Do(loginRequest)
	if err != nil {
		log.Error("Error making login request: error='%s'", err.Error())
		return nil, err
	}
	if loginResponse.StatusCode != 200 {
		log.Error("HTTP error making login request: status_code=%d", loginResponse.StatusCode)
		return nil, fmt.Errorf("http %d", loginResponse.StatusCode)
	}
	csrfToken, err = getCSRFTokenFromResponse(loginResponse)
	if err != nil {
		return nil, err
	}

	backupParams := url.Values{}
	backupParams.Add("download", "download")
	backupParams.Add("donotbackuprrd", "yes")
	backupParams.Add("__csrf_magic", csrfToken)
	if config.EncryptPassword != "" {
		backupParams.Add("encrypt", "yes")
		backupParams.Add("encrypt_password", config.EncryptPassword)
	}
	backupRequest, err := http.NewRequest("POST", backupURL, bytes.NewReader([]byte(backupParams.Encode())))
	if err != nil {
		log.Error("Error forming backup request: error='%s'", err.Error())
		return nil, err
	}
	backupRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	backupRequest.Header.Add("Accept", "*/*")

	backupResponse, err := client.Do(backupRequest)
	if err != nil {
		log.Error("Error making backup request: error='%s'", err.Error())
		return nil, err
	}
	if backupResponse.StatusCode != 200 {
		log.Error("HTTP error making backup request: status_code=%d", backupResponse.StatusCode)
		return nil, fmt.Errorf("http %d", backupResponse.StatusCode)
	}

	filePath := pukcab.GetFilePath("pfsense", config.HostAddress+".xml")
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Error("Error opening destination file: file_path='%s' error='%s'", filePath, err.Error())
		return nil, err
	}
	defer f.Close()

	if _, err := io.Copy(f, backupResponse.Body); err != nil {
		log.Error("Error writing backup file: file_path='%s' error='%s'", filePath, err.Error())
		return nil, err
	}

	return &pukcab.File{
		Path: filePath,
	}, nil
}

func getCSRFTokenFromResponse(resp *http.Response) (string, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("Error reading HTTP body: error='%s'", err.Error())
		return "", err
	}
	resp.Body.Close()
	csrfTokenPattern := regexp.MustCompile("var csrfMagicToken = \"[a-zA-Z0-9:;,]+\";")
	csrfTokenRaw := csrfTokenPattern.Find(body)
	if len(csrfTokenRaw) == 0 {
		log.Error("No CSRF token found")
		return "", fmt.Errorf("no csrf token found")
	}
	token := string(regexp.MustCompile("\".*\"").Find(csrfTokenRaw))
	token = token[1 : len(token)-1]
	log.Debug("CSRF token found: %s", token)
	return token, nil
}
