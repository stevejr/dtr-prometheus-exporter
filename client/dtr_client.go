package client

import (
  "context"
  "time"
  "github.com/stevejr/dtr-prometheus-exporter/config"
  "github.com/stevejr/dtr-prometheus-exporter/dtrconnector"
  api "github.com/stevejr/dtr-prometheus-exporter/api" 
  "crypto/tls"
  "encoding/json"
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
  Replicas  *[]ReplicaStats
}

// ReplicaStats - main struct to hold the replica stats
type ReplicaStats struct {
  ID                  string
  HealthyCount        int
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
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

  var dtrClusterStatus api.ClusterStatus
	
  jsonData, err := dtrconnector.MakeClientRequest(c.connectionString, c.tlsConfig, c.username, c.password)
  if err != nil {
    return nil, err
  }

  err = json.Unmarshal(jsonData, &dtrClusterStatus)
  if err != nil {
    return nil, err
  }
  
	replicas, err := getReplicaStats(ctx, &dtrClusterStatus)
	if err != nil {
		return nil, err
	}

	return &Stats{
		Replicas: replicas,
	}, nil
}

func getReplicaStats(ctx context.Context, rh *api.ClusterStatus) (*[]ReplicaStats, error) {
  var result []ReplicaStats
  var healthCount int

  for replica, health := range rh.ReplicaHealth {
    // fmt.Println("replica:", replica, "=>", "health:", health)
    if health == "OK" {
      healthCount = 1
    }
    result = append(result, ReplicaStats{
      ID: replica,
      HealthyCount: healthCount,
    }) 
  }

  return &result, nil
}