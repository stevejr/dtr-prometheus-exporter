package dtrconnector

import (
	"crypto/tls"
	"fmt"
	"github.com/stevejr/dtr-prometheus-exporter/config"
	"io/ioutil"
	"net/http"
	"strings"
)

func setAPIEndpoint(ae string, cs *string) {

	if ae[0:0] == "/" {
		ae = strings.Replace(ae, "/", "", 1)
		fmt.Printf("apiEndpoint: %s\n", ae)
	}
	*cs = fmt.Sprintf("%s/%s", *cs, ae)
}

// MakeRequest prepares a new http client request
func MakeRequest(cfg config.Config, tlsConfig *tls.Config) (*http.Response, error) {

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				RootCAs:            tlsConfig.RootCAs,
				Certificates:       tlsConfig.Certificates,
			},
		},
	}

	req, err := http.NewRequest("GET", cfg.DTR.DTRAPIAddress, nil)
	if err != nil {
		return nil, fmt.Errorf("Could not create new http request")
	}

	req.SetBasicAuth(cfg.DTR.Username, cfg.DTR.Password)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// TODO - Users should be able to pass in their own client and we should default to a http client if it isn't
// MakeClientRequest prepares a new http client request
func MakeClientRequest(cs string, tlsConfig *tls.Config, username string, password string, apiEndpoint string) ([]byte, error) {

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				RootCAs:            tlsConfig.RootCAs,
				Certificates:       tlsConfig.Certificates,
			},
		},
	}

	setAPIEndpoint(apiEndpoint, &cs)

	req, err := http.NewRequest("GET", cs, nil)
	if err != nil {
		return nil, fmt.Errorf("Could not create new http request")
	}

	req.SetBasicAuth(username, password)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	jsonData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return jsonData, nil
}
