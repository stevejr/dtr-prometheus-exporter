package client

import (
	"crypto/tls"
	"fmt"
	api "github.com/stevejr/dtr-prometheus-exporter/api"
	"github.com/stevejr/dtr-prometheus-exporter/config"
	"github.com/stevejr/dtr-prometheus-exporter/dtrconnector"
	"time"
)

// DTRClient - DTR Client struct
type DTRClient struct {
	connectionString string
	tlsConfig        *tls.Config
	timeout          time.Duration
	username         string
	password         string
	jobCount         uint
}

type clusterStatusStats struct {
	replicaStats *[]api.CSReplicaHealth
	tableStats   *[]api.TableStatus
	stats        *[]api.Stats
}

// Stats - main struct to hold captured stats
type Stats struct {
	CSReplicasHealth *[]api.CSReplicaHealth
	CSStats          *[]api.Stats
	CSTableStatus    *[]api.TableStatus
	JobCounts        *[]api.JobCounts
}

// TODO - Look at being able to allow a client object stored as part of this struct as it allows users to pass in their own client
// New - prepares a new http client
func New(cfg config.Config, tlsConfig *tls.Config) *DTRClient {
	return &DTRClient{
		connectionString: cfg.DTR.DTRAPIAddress,
		password:         cfg.DTR.Password,
		timeout:          cfg.Scrape.Timeout,
		tlsConfig:        tlsConfig,
		username:         cfg.DTR.Username,
		jobCount:         cfg.API.JobCount,
	}
}

// GetDTRStats - Retrieve the DTR API related stats
func (c *DTRClient) GetDTRStats() (*Stats, error) {
	clusterStatusStats, err := getClusterStatusStats(c)
	if err != nil {
		return nil, err
	}

	jobCounts, err := getJobStats(c)
	if err != nil {
		return nil, err
	}

	return &Stats{
		CSReplicasHealth: clusterStatusStats.replicaStats,
		CSStats:          clusterStatusStats.stats,
		CSTableStatus:    clusterStatusStats.tableStats,
		JobCounts:        jobCounts,
	}, nil
}

func getClusterStatusStats(c *DTRClient) (*clusterStatusStats, error) {
	apiEndpoint := api.CSAPIEndpoint

	jsonData, err := dtrconnector.MakeClientRequest(c.connectionString, c.tlsConfig, c.username, c.password, apiEndpoint)
	if err != nil {
		return nil, err
	}

	replicas, err := api.GetCSReplicaHealthStats(jsonData)
	if err != nil {
		return nil, err
	}

	tablestatus, err := api.GetCSTableStatus(jsonData)
	if err != nil {
		return nil, err
	}

	stats, err := api.GetCSStats(jsonData)
	if err != nil {
		return nil, err
	}

	return &clusterStatusStats{
		replicaStats: replicas,
		tableStats:   tablestatus,
		stats:        stats}, nil
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
