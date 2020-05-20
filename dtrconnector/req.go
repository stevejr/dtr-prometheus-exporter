package dtrconnector

import (
	"crypto/tls"
	"fmt"
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
