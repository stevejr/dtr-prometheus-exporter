package collector

import (
	"fmt"
  dtr "github.com/stevejr/dtr-prometheus-exporter/client"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"strings"
	api "github.com/stevejr/dtr-prometheus-exporter/api"
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

type csTableStorageEngineMetrics struct {
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
	csTableStorageEngineCounts	csTableStorageEngineMetrics	
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
		readDocsSec: prometheus.NewDesc(fmt.Sprintf("dtr_table_queryengine_read_docs_seconds"),
		fmt.Sprintf("DTR table query engine read docs per second per replica"), labels, nil),
		readDocsTotal: prometheus.NewDesc(fmt.Sprintf("dtr_table_queryengine_read_docs_total"),
		fmt.Sprintf("DTR table query engine read docs total per replica"), labels, nil),
		writtenDocsSec: prometheus.NewDesc(fmt.Sprintf("dtr_table_queryengine_written_docs_seconds"),
		fmt.Sprintf("DTR table query engine written docs per second per replica"), labels, nil),
		writtenDocsTotal: prometheus.NewDesc(fmt.Sprintf("dtr_table_queryengine_written_docs_total"),
		fmt.Sprintf("DTR table query engine written docs total per replica"), labels, nil),
	}
}

func newCSTableStorageEngineMetrics(labels ...string) csTableStorageEngineMetrics {

	return csTableStorageEngineMetrics{
		preallocatedBytes: prometheus.NewDesc(fmt.Sprintf("dtr_table_disk_preallocated_bytes"),
		fmt.Sprintf("DTR table disk space preallocated bytes per replica"), labels, nil),
		readBytesSec: prometheus.NewDesc(fmt.Sprintf("dtr_table_disk_read_bytes_seconds"),
		fmt.Sprintf("DTR table disk read bytes per second per replica"), labels, nil),
		readBytesTotal: prometheus.NewDesc(fmt.Sprintf("dtr_table_disk_read_total_bytes"),
		fmt.Sprintf("DTR table disk read bytes total per replica"), labels, nil),
		usedDataBytes: prometheus.NewDesc(fmt.Sprintf("dtr_table_disk_used_data_bytes"),
		fmt.Sprintf("DTR table disk space used data bytes per replica"), labels, nil),
		usedGarbageBytes: prometheus.NewDesc(fmt.Sprintf("dtr_table_disk_used_garbage_bytes"),
		fmt.Sprintf("DTR table disk space used garbage bytes per replica"), labels, nil),
		usedMetadataBytes: prometheus.NewDesc(fmt.Sprintf("dtr_table_disk_used_metadata_bytes"),
		fmt.Sprintf("DTR table disk space use metadata bytes per replica"), labels, nil),
		writtenBytesSec: prometheus.NewDesc(fmt.Sprintf("dtr_table_disk_written_bytes_second"),
		fmt.Sprintf("DTR table disk written bytes per second replica"), labels, nil),
		writtenBytesTotal: prometheus.NewDesc(fmt.Sprintf("dtr_table_disk_written_total_bytes"),
		fmt.Sprintf("DTR table disk written bytes total per replica"), labels, nil),
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
		
		csTableQueryEngineCounts: newCSTableQueryEngineMetrics("type", "db", "table", "replica"),
		csTableStorageEngineCounts: newCSTableStorageEngineMetrics("type", "db", "table", "replica"),

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
	ch <- metrics.readDocsTotal
	ch <- metrics.writtenDocsSec
	ch <- metrics.writtenDocsTotal
}

func describeCSTableStorageEngineMetrics(ch chan<- *prometheus.Desc, metrics *csTableStorageEngineMetrics) {
	ch <- metrics.preallocatedBytes
	ch <- metrics.readBytesSec
	ch <- metrics.readBytesTotal
	ch <- metrics.usedDataBytes
	ch <- metrics.usedGarbageBytes
	ch <- metrics.usedMetadataBytes
	ch <- metrics.writtenBytesSec
	ch <- metrics.writtenBytesTotal
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
	describeCSTableStorageEngineMetrics(ch, &c.csTableStorageEngineCounts)
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
		collectCSTableStorageEngineMetrics(c, ch, stats)

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

	replica := ""

	for _, stat := range *stats.CSStats {
		if stat.Server != "" {
			replica = strings.Split(stat.Server, "_")[2]
		}
		idType := stat.ID[0]		
		collectTableQueryEngineCounts(ch, &c.csTableQueryEngineCounts, &stat.QueryEngine, idType, stat.Db, stat.Table, replica)
	}
}

func collectCSTableStorageEngineMetrics(c *Collector, ch chan<- prometheus.Metric, stats *dtr.Stats) {

	replica := ""

	for _, stat := range *stats.CSStats {
		if stat.Server != "" {
			replica = strings.Split(stat.Server, "_")[2]
		}
		idType := stat.ID[0]	
		collectTableStorageEngineCounts(ch, &c.csTableStorageEngineCounts, &stat.StorageEngine, idType, stat.Db, stat.Table, replica)
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

func collectTableQueryEngineCounts(ch chan<- prometheus.Metric, metrics *csTableQueryEngineMetrics, qe *api.QueryEngine, labelValues ...string) {
	ch <- prometheus.MustNewConstMetric(metrics.readDocsSec, prometheus.GaugeValue, float64(qe.ReadDocsPerSec), labelValues...)	
	ch <- prometheus.MustNewConstMetric(metrics.readDocsTotal, prometheus.GaugeValue, float64(qe.ReadDocsTotal), labelValues...)
	ch <- prometheus.MustNewConstMetric(metrics.writtenDocsSec, prometheus.GaugeValue, float64(qe.WrittenDocsPerSec), labelValues...)	
	ch <- prometheus.MustNewConstMetric(metrics.writtenDocsTotal, prometheus.GaugeValue, float64(qe.WrittenDocsTotal), labelValues...)
}

func collectTableStorageEngineCounts(ch chan<- prometheus.Metric, metrics *csTableStorageEngineMetrics, se *api.StorageEngine, labelValues ...string) {
	ch <- prometheus.MustNewConstMetric(metrics.preallocatedBytes, prometheus.GaugeValue, float64(se.Disk.SpaceUsage.PreallocatedBytes), labelValues...)
	ch <- prometheus.MustNewConstMetric(metrics.readBytesSec, prometheus.GaugeValue, float64(se.Disk.ReadBytesPerSec), labelValues...)	
	ch <- prometheus.MustNewConstMetric(metrics.readBytesTotal, prometheus.GaugeValue, float64(se.Disk.ReadBytesTotal), labelValues...)	
	ch <- prometheus.MustNewConstMetric(metrics.usedDataBytes, prometheus.GaugeValue, float64(se.Disk.SpaceUsage.DataBytes), labelValues...)	
	ch <- prometheus.MustNewConstMetric(metrics.usedGarbageBytes, prometheus.GaugeValue, float64(se.Disk.SpaceUsage.GarbageBytes), labelValues...)	
	ch <- prometheus.MustNewConstMetric(metrics.usedMetadataBytes, prometheus.GaugeValue, float64(se.Disk.SpaceUsage.MetadataBytes), labelValues...)	
	ch <- prometheus.MustNewConstMetric(metrics.writtenBytesSec, prometheus.GaugeValue, float64(se.Disk.WrittenBytesPerSec), labelValues...)	
	ch <- prometheus.MustNewConstMetric(metrics.writtenBytesTotal, prometheus.GaugeValue, float64(se.Disk.WrittenBytesTotal), labelValues...)		
}
func collectJobCounts(ch chan<- prometheus.Metric, metrics *jobMetrics, count int, labelValues ...string) {
	ch <- prometheus.MustNewConstMetric(metrics.jobCount, prometheus.GaugeValue, float64(count), labelValues...)	
}