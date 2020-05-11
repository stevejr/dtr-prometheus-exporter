package collector

import (
	"fmt"
  dtr "github.com/stevejr/dtr-prometheus-exporter/client"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

type csReplicaHealthMetrics struct {
  health  					*prometheus.Desc
}

type csTableStatusMetrics struct {
	readiness 				*prometheus.Desc
}

type jobMetrics struct {
	jobCount				 	*prometheus.Desc		
}

type sumOfJobs struct {
  action 				    string
  status            string 
}

type csTableQueryEngineMetrics struct {
	readDocsSec									*prometheus.Desc
	readDocsTotal								*prometheus.Desc
	writtenDocsSec							*prometheus.Desc
	writtenDocsTotal						*prometheus.Desc
}

type csTableDiskMetrics struct {
	preallocatedBytes							*prometheus.Desc
	readBytesSec									*prometheus.Desc
	readBytesTotal								*prometheus.Desc
	usedDataBytes								  *prometheus.Desc
	usedGarbageBytes							*prometheus.Desc
	usedMetadataBytes							*prometheus.Desc	
	writtenBytesSec								*prometheus.Desc
	writtenBytesTotal							*prometheus.Desc
}

// Collector - Prometheus Collector struct
type Collector struct {
	client *dtr.DTRClient
	log    *logrus.Logger

  up *prometheus.Desc
  
	csReplicaHealthCounts       csReplicaHealthMetrics
	csTableStatusCounts					csTableStatusMetrics				

	csTableQueryEngineCounts		csTableQueryEngineMetrics
	csTableDiskCounts						csTableDiskMetrics	
	jobActionCounts							jobMetrics
}

func newCSReplicaHealthMetrics(labels ...string) csReplicaHealthMetrics {

	return csReplicaHealthMetrics{
		health: prometheus.NewDesc(fmt.Sprintf("dtr_replica_health_total"),
			fmt.Sprintf("DTR replica health total"), labels, nil),
	}
}

func newCSTableStatusMetrics(labels ...string) csTableStatusMetrics {

	return csTableStatusMetrics{
		readiness: prometheus.NewDesc(fmt.Sprintf("dtr_table_replica_status_total"),
			fmt.Sprintf("DTR table replica status total"), labels, nil),
	}
}

func newCSTableQueryEngineMetrics(labels ...string) csTableQueryEngineMetrics {

	return csTableQueryEngineMetrics{
		readDocsSec: prometheus.NewDesc(fmt.Sprintf("dtr_table_qe_read_docs_seconds"),
		fmt.Sprintf("DTR table query engine read docs per second per replica"), labels, nil),

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

		csReplicaHealthCounts: newCSReplicaHealthMetrics("health"),
		csTableStatusCounts: newCSTableStatusMetrics("db", "table", "status"),
		
		csTableQueryEngineCounts: newCSTableQueryEngineMetrics("db", "table", "replica"),

		jobActionCounts: newJobMetrics("action", "status"),
	}
}

func describeCSReplicaHealthMetrics(ch chan<- *prometheus.Desc, metrics *csReplicaHealthMetrics) {
	ch <- metrics.health
}

func describeCSTableStatusMetrics(ch chan<- *prometheus.Desc, metrics *csTableStatusMetrics) {
	ch <- metrics.readiness
}

func describeCSTableQueryEngineMetrics(ch chan<- *prometheus.Desc, metrics *csTableQueryEngineMetrics) {
	ch <- metrics.readDocsSec
}

func describeJobMetrics(ch chan<- *prometheus.Desc, metrics *jobMetrics) {
	ch <- metrics.jobCount
}

// Describe - called to get descriptors of the metrics provided by the collector
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.up

	describeCSReplicaHealthMetrics(ch, &c.csReplicaHealthCounts)
	describeCSTableStatusMetrics(ch, &c.csTableStatusCounts)
	describeCSTableQueryEngineMetrics(ch, &c.csTableQueryEngineCounts)
	describeJobMetrics(ch, &c.jobActionCounts)
}

func scrapeDTR(c *Collector, ch chan<- prometheus.Metric) {
	if stats, err := c.client.GetDTRStats(); err != nil {
		ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 0)

		c.log.WithError(err).Error("Error during scrape")
	} else {
		ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 1)

		collectCSReplicaHealthMetrics(c, ch, stats)
		collectCSTableStatusMetrics(c, ch, stats)
		collectCSTableQueryEngineMetrics(c, ch, stats)

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

  collectReplicaHealth(ch, &c.csReplicaHealthCounts, health, "healthy")
  collectReplicaHealth(ch, &c.csReplicaHealthCounts, unhealthy, "unhealthy")
}

func collectCSTableStatusMetrics(c *Collector, ch chan<- prometheus.Metric, stats *dtr.Stats) {

	for _, table := range *stats.CSTableStatus {
		db := table.Db
		name := table.Name
		ready := 0
		notReady := 0
		
		for _, replica := range table.Shards[0].Replicas {
			if replica.State == "ready" {
				ready++
			} else {
				notReady++
			}
		}

		collectTableStatus(ch, &c.csTableStatusCounts, ready, db, name, "ready")
		collectTableStatus(ch, &c.csTableStatusCounts, notReady, db, name, "notready")
	}
}

func collectCSTableQueryEngineMetrics(c *Collector, ch chan<- prometheus.Metric, stats *dtr.Stats) {

	for _, stat := range *stats.CSStats {
		fmt.Printf ("stat: %+v\n", stat)
		collectTableQueryEngineCounts(ch, &c.csTableQueryEngineCounts, stat.QueryEngine.ReadDocsPerSec, stat.Db, stat.Table, stat.Server)
	}
}

func collectJobMetrics(c *Collector, ch chan<- prometheus.Metric, stats *dtr.Stats) {

  jobStats := make(map[sumOfJobs]int) 
  
  for _, job := range *stats.JobCounts {
    jobStats[sumOfJobs{job.Action, job.Status}] ++
  }  

  for jobStat, jobCount := range jobStats {    
    collectJobCounts(ch, &c.jobActionCounts, jobCount, jobStat.action, jobStat.status)
  }
}

func collectReplicaHealth(ch chan<- prometheus.Metric, metrics *csReplicaHealthMetrics, health int, labelValues ...string) {
	ch <- prometheus.MustNewConstMetric(metrics.health, prometheus.GaugeValue, float64(health), labelValues...)
}

func collectTableStatus(ch chan<- prometheus.Metric, metrics *csTableStatusMetrics, count int, labelValues ...string) {
	ch <- prometheus.MustNewConstMetric(metrics.readiness, prometheus.GaugeValue, float64(count), labelValues...)	
}

func collectTableQueryEngineCounts(ch chan<- prometheus.Metric, metrics *csTableQueryEngineMetrics, counts float64, labelValues ...string) {
	ch <- prometheus.MustNewConstMetric(metrics.readDocsSec, prometheus.GaugeValue, float64(counts), labelValues...)	
}

func collectJobCounts(ch chan<- prometheus.Metric, metrics *jobMetrics, count int, labelValues ...string) {
	ch <- prometheus.MustNewConstMetric(metrics.jobCount, prometheus.GaugeValue, float64(count), labelValues...)	
}