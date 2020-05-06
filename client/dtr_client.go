package client

import (
  "time"
  "github.com/stevejr/dtr-prometheus-exporter/config"
  "github.com/stevejr/dtr-prometheus-exporter/dtrconnector"
  api "github.com/stevejr/dtr-prometheus-exporter/api" 
  "crypto/tls"
)

// DTRClient - DTR Client struct
type DTRClient struct {
  connectionString    string
  tlsConfig           *tls.Config
  timeout             time.Duration
  username            string
  password            string
}

// Stats - main struct to hold captured stats
type Stats struct {
  CSReplicasHealth  *[]api.CSReplicaHealth
}

// New - prepares a new http client
func New (cfg config.Config, tlsConfig *tls.Config) (*DTRClient) {

  return &DTRClient{
    connectionString: cfg.DTR.DTRAPIAddress,
    password: cfg.DTR.Password,
    timeout: cfg.Scrape.Timeout,
    tlsConfig: tlsConfig,
    username: cfg.DTR.Username,
  }
}

// GetClusterStatusStats - Retrieve the ClusterStatus API related stats
func (c *DTRClient) GetClusterStatusStats() (*Stats, error) {
	
	apiEndpoint := api.CSAPIEndpoint
	
  jsonData, err := dtrconnector.MakeClientRequest(c.connectionString, c.tlsConfig, c.username, c.password, apiEndpoint)
  if err != nil {
    return nil, err
  }

	replicas, err := api.GetCSReplicaHealthStats(jsonData)

	return &Stats{
		CSReplicasHealth: replicas,
	}, nil
}