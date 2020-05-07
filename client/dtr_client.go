package client

import (
  "time"
  "github.com/stevejr/dtr-prometheus-exporter/config"
  "github.com/stevejr/dtr-prometheus-exporter/dtrconnector"
  api "github.com/stevejr/dtr-prometheus-exporter/api" 
	"crypto/tls"
	"fmt"
)

// DTRClient - DTR Client struct
type DTRClient struct {
  connectionString    string
  tlsConfig           *tls.Config
  timeout             time.Duration
  username            string
	password            string
	jobCount						uint
}

// Stats - main struct to hold captured stats
type Stats struct {
	CSReplicasHealth  *[]api.CSReplicaHealth
	JobCounts					*[]api.JobCounts
}

// New - prepares a new http client
func New (cfg config.Config, tlsConfig *tls.Config) (*DTRClient) {
  return &DTRClient{
    connectionString: cfg.DTR.DTRAPIAddress,
    password: cfg.DTR.Password,
    timeout: cfg.Scrape.Timeout,
    tlsConfig: tlsConfig,
		username: cfg.DTR.Username,
		jobCount: cfg.API.JobCount,
  }
}

// GetDTRStats - Retrieve the DTR API related stats
func (c *DTRClient) GetDTRStats() (*Stats, error) {
  replicas, err := getClusterStatusStats(c)	
	if err != nil {
    return nil, err
	}
	
	jobCounts, err := getJobStats(c)
	if err != nil {
    return nil, err
  }	

	return &Stats{
		CSReplicasHealth: replicas,
		JobCounts: jobCounts,	
	}, nil
}

func getClusterStatusStats(c *DTRClient) (*[]api.CSReplicaHealth, error) {
	apiEndpoint := api.CSAPIEndpoint
	
  jsonData, err := dtrconnector.MakeClientRequest(c.connectionString, c.tlsConfig, c.username, c.password, apiEndpoint)
  if err != nil {
    return nil, err
  }

	replicas, err := api.GetCSReplicaHealthStats(jsonData)
  if err != nil {
    return nil, err
  }
	
	return replicas, nil
}

func getJobStats(c *DTRClient) (*[]api.JobCounts, error) {
	apiEndpoint := api.JobsAPIEndpoint
	
	apiEndpoint = fmt.Sprintf("%s&limit=%d", apiEndpoint, c.jobCount)
	
  jsonData, err := dtrconnector.MakeClientRequest(c.connectionString, c.tlsConfig, c.username, c.password, apiEndpoint)
  if err != nil {
    return nil, err
  }

	jobCounts, err := api.GetJobCountsStats(jsonData)
  if err != nil {
    return nil, err
  }

	return jobCounts, nil
}