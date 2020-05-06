package dtrconnector

import (
  "fmt"
  "net/http"
  "github.com/stevejr/dtr-prometheus-exporter/config"
  "io/ioutil"
  "crypto/tls"
)

// MakeRequest prepares a new http client request
func MakeRequest (cfg config.Config, tlsConfig *tls.Config) (*http.Response, error) {

  client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				RootCAs:      tlsConfig.RootCAs,
				Certificates: tlsConfig.Certificates,
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

// MakeClientRequest prepares a new http client request
func MakeClientRequest (cs string, tlsConfig *tls.Config, username string, password string) ([]byte, error) {

	client := &http.Client{
    Transport: &http.Transport{
      TLSClientConfig: &tls.Config{
        InsecureSkipVerify: true,
        RootCAs:      tlsConfig.RootCAs,
        Certificates: tlsConfig.Certificates,
      },
    },
	}
	
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
