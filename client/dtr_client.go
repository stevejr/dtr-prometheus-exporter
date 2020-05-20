package client

import (
	"crypto/tls"
	"fmt"
	api "github.com/stevejr/dtr-prometheus-exporter/api"
	"github.com/stevejr/dtr-prometheus-exporter/config"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// DTRClient - DTR Client struct
type DTRClient struct {
	httpClient *http.Client
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

func newHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{Timeout: timeout}
}

func newHTTPSClient(timeout time.Duration, tlsConfig *tls.Config) *http.Client {
	var hc *http.Client 
	
	if tlsConfig.Certificates == nil {
		fmt.Println("HTTPS Client will be configured with RootCA only")
		hc = &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs:            tlsConfig.RootCAs,
				},
			},
		}
	} else {
		fmt.Println("HTTPS Client will be configured with RootCA and Client Certs")
		hc = &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs:            tlsConfig.RootCAs,
					Certificates:       tlsConfig.Certificates,
				},
			},
		}
	}

	return hc
}


// New - prepares a new http client
func New(httpClient *http.Client, cfg config.Config, tlsConfig *tls.Config) *DTRClient {
	hc := httpClient

	if hc == nil {
		fmt.Println("Using default Golang HTTP Client")
		if tlsConfig == nil {
			fmt.Println("Creating standard HTTP Client")
			hc = newHTTPClient(cfg.Scrape.Timeout) 
		} else {
			fmt.Println("Creating standard HTTPS Client")
			hc = newHTTPSClient(cfg.Scrape.Timeout, tlsConfig) 
		}
	}

	return &DTRClient{
		httpClient: hc,
		connectionString: cfg.DTR.DTRAPIAddress,
		password:         cfg.DTR.Password,
		timeout:          cfg.Scrape.Timeout,
		tlsConfig:        tlsConfig,
		username:         cfg.DTR.Username,
		jobCount:         cfg.API.JobCount,
	}
}

func setAPIEndpoint(ae string, cs string) (string, error) {

	if ae[0:0] == "/" {
		ae = strings.Replace(ae, "/", "", 1)
		fmt.Printf("apiEndpoint: %s\n", ae)
	}
	return fmt.Sprintf("%s/%s", cs, ae), nil
}

// MakeRequest prepares a new http client request
func (c *DTRClient) MakeRequest(apiEndpoint string) (*http.Response, error) {

	apiEP, err := setAPIEndpoint(apiEndpoint, c.connectionString)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", apiEP, nil)
	if err != nil {
		return nil, fmt.Errorf("Could not create new http request")
	}

	req.SetBasicAuth(c.username, c.password)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
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

	resp, err := c.MakeRequest(apiEndpoint)
	if err != nil {
		return nil, err
	}

	jsonData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

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

	resp, err := c.MakeRequest(apiEndpoint)
	if err != nil {
		return nil, err
	}

	jsonData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	jobCounts, err := api.GetJobCountsStats(jsonData)
	if err != nil {
		return nil, err
	}

	return jobCounts, nil
}
