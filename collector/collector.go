package collector

import (
	"fmt"
	"strings"
  dtr "github.com/stevejr/dtr-prometheus-exporter/client"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

type csReplicaHealthMetrics struct {
  health           *prometheus.Desc
}

type jobMetrics struct {
	jobCount				 	*prometheus.Desc		
}

type sumOfJobs struct {
  Action 				    string
  Status            string 
}

// Collector - Prometheus Collector struct
type Collector struct {
	client *dtr.DTRClient
	log    *logrus.Logger

  up *prometheus.Desc
  
  csReplicaHealthyCounts      csReplicaHealthMetrics
	csReplicaUnHealthyCounts    csReplicaHealthMetrics
	
	jobActionCounts							jobMetrics
}

func newCSReplicaHealthMetrics(itemName string) csReplicaHealthMetrics {

	return csReplicaHealthMetrics{
		health: prometheus.NewDesc(fmt.Sprintf("dtr_replica_%s_total", strings.ToLower(itemName)),
			fmt.Sprintf("%s count of replicas", itemName), []string{"name"}, nil),
	}
}

func newJobMetrics(labels ...string) jobMetrics {

	return jobMetrics{
		jobCount: prometheus.NewDesc(fmt.Sprintf("dtr_job_total"),
			fmt.Sprintf("DTR job total"), labels, nil),
	}
}

// New - Creates new Prometheus collector
func New(client *dtr.DTRClient, log *logrus.Logger) *Collector {
	return &Collector{
		client: client,
		log:    log,

		up: prometheus.NewDesc("dtr_up", "Whether the DTR scrape was successful", nil, nil),

    csReplicaHealthyCounts: newCSReplicaHealthMetrics("Healthy"),
		csReplicaUnHealthyCounts: newCSReplicaHealthMetrics("UnHealthy"),
		
		jobActionCounts: newJobMetrics("action", "status"),
	}
}

func describeCSReplicaHealthMetrics(ch chan<- *prometheus.Desc, metrics *csReplicaHealthMetrics) {
	ch <- metrics.health
}

func describeJobMetrics(ch chan<- *prometheus.Desc, metrics *jobMetrics) {
	ch <- metrics.jobCount
}

// Describe - called to get descriptors of the metrics provided by the collector
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.up

	describeCSReplicaHealthMetrics(ch, &c.csReplicaHealthyCounts)
	describeCSReplicaHealthMetrics(ch, &c.csReplicaUnHealthyCounts)

	describeJobMetrics(ch, &c.jobActionCounts)
}

func scrapeDTR(c *Collector, ch chan<- prometheus.Metric) {
	if stats, err := c.client.GetDTRStats(); err != nil {
		ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 0)

		c.log.WithError(err).Error("Error during scrape")
	} else {
		ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 1)

    collectCSReplicaHealthMetrics(c, ch, stats)
    c.log.Info("Collected DTR ClusterStatus API data")
    collectJobMetrics(c, ch, stats)
    c.log.Info("Collected DTR Jobs API data")
	}
}

// Collect - called to get the metric values
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	c.log.Info("Running scrape")

  scrapeDTR(c, ch)
  c.log.Info("Scrape completed")
}

func collectCSReplicaHealthMetrics(c *Collector, ch chan<- prometheus.Metric, stats *dtr.Stats) {
  var health, unhealthy int
  for _, replica := range *stats.CSReplicasHealth {
    if replica.HealthyCount == 1 {
      health++
    } else {
      unhealthy++
    }
  }

  collectReplicaHealth(ch, &c.csReplicaHealthyCounts, health, "healthy")
  collectReplicaHealth(ch, &c.csReplicaUnHealthyCounts, unhealthy, "unhealthy")
}

func collectJobMetrics(c *Collector, ch chan<- prometheus.Metric, stats *dtr.Stats) {

  jobStats := make(map[sumOfJobs]int) 
  
  for _, job := range *stats.JobCounts {
    jobStats[sumOfJobs{job.Action, job.Status}] ++
  }  

  for jobStat, jobCount := range jobStats {    
    collectJobCounts(ch, &c.jobActionCounts, jobCount, jobStat.Action, jobStat.Status)
  }
}

func collectReplicaHealth(ch chan<- prometheus.Metric, metrics *csReplicaHealthMetrics, health int, labelValues ...string) {
	ch <- prometheus.MustNewConstMetric(metrics.health, prometheus.GaugeValue, float64(health), labelValues...)
}

func collectJobCounts(ch chan<- prometheus.Metric, metrics *jobMetrics, count int, labelValues ...string) {
	ch <- prometheus.MustNewConstMetric(metrics.jobCount, prometheus.GaugeValue, float64(count), labelValues...)	
}