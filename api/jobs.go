package api

import (
	"encoding/json"
	"time"
)

const (
	// JobsAPIEndpoint - the DTR API endpoint for Jobs
	JobsAPIEndpoint = "api/v0/jobs?action=any&worker=any&running=any&start=0"
)

// Jobs - DTR Jobs list
type Jobs struct {
	Jobs []Job `json:"jobs"`
}

// Job - DTR Job
type Job struct {
	ID           string      `json:"id"`
	RetryFromID  string      `json:"retryFromID"`
	WorkerID     string      `json:"workerID"`
	Status       string      `json:"status"`
	ScheduledAt  time.Time   `json:"scheduledAt"`
	LastUpdated  time.Time   `json:"lastUpdated"`
	Action       string      `json:"action"`
	RetriesLeft  int         `json:"retriesLeft"`
	RetriesTotal int         `json:"retriesTotal"`
	CapacityMap  interface{} `json:"capacityMap"`
	Parameters   interface{} `json:"parameters"`
	Deadline     string      `json:"deadline"`
	StopTimeout  string      `json:"stopTimeout"`
}

// JobCounts - Count of jobs execute by Action and Status
type JobCounts struct {
	Action string
	Status string
}

// GetJobCountsStats - Get the JobCounts stats
func GetJobCountsStats(jsonData []byte) (*[]JobCounts, error) {
	var result []JobCounts
	var jobs Jobs

	err := json.Unmarshal(jsonData, &jobs)
	if err != nil {
		return nil, err
	}

	for _, job := range jobs.Jobs {
		result = append(result, JobCounts{
			Action: job.Action,
			Status: job.Status})
	}

	return &result, nil
}
