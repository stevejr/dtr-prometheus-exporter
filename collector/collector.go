package collector

import (
	"fmt"
	"strings"
  dtr "github.com/stevejr/dtr-prometheus-exporter/client"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

type replicaHealthMetrics struct {
  health           *prometheus.Desc
}

// Collector - Prometheus Collector struct
type Collector struct {
	client *dtr.DTRClient
	log    *logrus.Logger

  up *prometheus.Desc
  
  replicaHealthyCounts      replicaHealthMetrics
  replicaUnHealthyCounts    replicaHealthMetrics
}

func newReplicaHealthMetrics(itemName string) replicaHealthMetrics {

	return replicaHealthMetrics{
		health: prometheus.NewDesc(fmt.Sprintf("dtr_replica_%s_count", strings.ToLower(itemName)),
			fmt.Sprintf("%s count of replicas", itemName), []string{"name"}, nil),
	}
}

// New - Creates new Prometheus collector
func New(client *dtr.DTRClient, log *logrus.Logger) *Collector {
	return &Collector{
		client: client,
		log:    log,

		up: prometheus.NewDesc("dtr_up", "Whether the DTR scrape was successful", nil, nil),

    replicaHealthyCounts: newReplicaHealthMetrics("Healthy"),
    replicaUnHealthyCounts: newReplicaHealthMetrics("UnHealthy"),
	}
}

func describeReplicaHealthMetrics(ch chan<- *prometheus.Desc, metrics *replicaHealthMetrics) {
	ch <- metrics.health
}

// Describe - called to get descriptors of the metrics provided by the collector
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.up

	describeReplicaHealthMetrics(ch, &c.replicaHealthyCounts)
	describeReplicaHealthMetrics(ch, &c.replicaUnHealthyCounts)
}

// Collect - called to get the metric values
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	c.log.Info("Running scrape")

	if stats, err := c.client.GetClusterStatusStats(); err != nil {
		ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 0)

		c.log.WithError(err).Error("Error during scrape")
	} else {
		ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 1)

		collectHealthMetrics(c, ch, stats)

		c.log.Info("Scrape completed")
	}
}

func collectHealthMetrics(c *Collector, ch chan<- prometheus.Metric, stats *dtr.Stats) {
  var health, unhealthy int
  for _, replica := range *stats.Replicas {
    if replica.HealthyCount == 1 {
      health++
    } else {
      unhealthy++
    }
  }

  collectSizes(ch, &c.replicaHealthyCounts, health, "Healthy")
  collectSizes(ch, &c.replicaUnHealthyCounts, unhealthy, "UnHealthy")
}

func collectSizes(ch chan<- prometheus.Metric, metrics *replicaHealthMetrics, health int, labelValues ...string) {
	ch <- prometheus.MustNewConstMetric(metrics.health, prometheus.GaugeValue, float64(health), labelValues...)
}

