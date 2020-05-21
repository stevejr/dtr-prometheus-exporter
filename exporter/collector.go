package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
	api "github.com/stevejr/dtr-prometheus-exporter/api"
	dtr "github.com/stevejr/dtr-prometheus-exporter/client"
	"strings"
	"time"
)

const (
	statIDTypeIDX    = 0
	replicaIDStatIDX = 2
)

func scrapeDTR(e *DTRRethinkDBExporter, ch chan<- prometheus.Metric) {
	start := time.Now()
	if dtrStats, err := e.client.GetDTRStats(); err != nil {
		ch <- prometheus.MustNewConstMetric(e.metrics.up, prometheus.GaugeValue, 0)

		e.log.WithError(err).Error("Error during scrape")
	} else {
		ch <- prometheus.MustNewConstMetric(e.metrics.up, prometheus.GaugeValue, 1)

		processStats(e, ch, dtrStats)
	}

	elapsed := time.Since(start)
	ch <- prometheus.MustNewConstMetric(e.metrics.scrapeLatency, prometheus.GaugeValue, elapsed.Seconds())
}

// Collect - called to get the metric values
func (e *DTRRethinkDBExporter) Collect(ch chan<- prometheus.Metric) {
	e.log.Info("Running scrape")
	scrapeDTR(e, ch)
	e.log.Info("Scrape completed")
}

func processStats(e *DTRRethinkDBExporter, ch chan<- prometheus.Metric, dtrStats *dtr.Stats) {

	var healthy, unhealthy, tableCount int

	e.log.Info("Collecting DTR ClusterStatus API data")
	for _, stat := range *dtrStats.CSStats {

		replica := ""

		if len(stat.ID) == 0 {
			e.log.Error("unexpected empty stat id")
			return
		}

		if stat.Server != "" {
			replica = strings.Split(stat.Server, "_")[replicaIDStatIDX]
		}

		switch stat.ID[statIDTypeIDX] {
		case "cluster":
			e.processClusterStat(stat.QueryEngine, ch)
		case "server":
			e.processServerStat(stat.QueryEngine, ch, replica)
		case "table":
			e.processTableStat(stat.QueryEngine, ch, stat.Db, stat.Table, replica)
		case "table_server":
			e.processTableServerStat(stat.QueryEngine, stat.StorageEngine, ch, stat.Db, stat.Table, replica)
		default:
			e.log.Errorf("unexpected stat id: '%v'", stat.ID[statIDTypeIDX])
			return
		}
	}

	for _, table := range *dtrStats.CSTableStatus {
		db := table.Db
		name := table.Name
		ready := 0
		notReady := 0
		tableCount++

		for _, replica := range table.Shards[0].Replicas {
			if replica.State == "ready" {
				ready++
			} else {
				notReady++
			}
		}

		ch <- prometheus.MustNewConstMetric(e.metrics.tableStatus, prometheus.GaugeValue, float64(ready), db, name, "ready")
		ch <- prometheus.MustNewConstMetric(e.metrics.tableStatus, prometheus.GaugeValue, float64(notReady), db, name, "notready")
	}

	// Total tables
	ch <- prometheus.MustNewConstMetric(e.metrics.tableTotal, prometheus.GaugeValue, float64(tableCount))

	for _, replica := range *dtrStats.CSReplicasHealth {
		if replica.HealthyCount == 1 {
			healthy++
		} else {
			unhealthy++
		}
	}

	ch <- prometheus.MustNewConstMetric(e.metrics.serverHealth, prometheus.GaugeValue, float64(healthy), "healthy")
	ch <- prometheus.MustNewConstMetric(e.metrics.serverHealth, prometheus.GaugeValue, float64(unhealthy), "unhealthy")
	ch <- prometheus.MustNewConstMetric(e.metrics.serverTotal, prometheus.GaugeValue, float64(unhealthy+healthy))

	e.log.Info("Collected DTR Jobs API data")

	jobStats := make(map[sumOfJobs]int)

	for _, job := range *dtrStats.JobCounts {
		jobStats[sumOfJobs{job.Action, job.Status}]++
	}

	for jobStat, jobCount := range jobStats {
		if e.client.Debug {
			e.log.Debugf("func: %s :- jobStat received: '%v'", "processStats", jobStat)
		}

		ch <- prometheus.MustNewConstMetric(e.metrics.jobTotals, prometheus.GaugeValue, float64(jobCount), jobStat.status, jobStat.action)
	}

}

func (e *DTRRethinkDBExporter) processClusterStat(stat api.QueryEngine, ch chan<- prometheus.Metric) {
	if e.client.Debug {
		e.log.Debugf("func: %s :- stat received: '%v'", "processClusterStat", stat)
	}

	ch <- prometheus.MustNewConstMetric(e.metrics.clusterClientConnections, prometheus.GaugeValue, float64(stat.ClientConnections))
	ch <- prometheus.MustNewConstMetric(e.metrics.clusterClientsActive, prometheus.GaugeValue, float64(stat.ClientsActive))
	ch <- prometheus.MustNewConstMetric(e.metrics.clusterDocsPerSecond, prometheus.GaugeValue, float64(stat.ReadDocsPerSec), readOperation)
	ch <- prometheus.MustNewConstMetric(e.metrics.clusterDocsPerSecond, prometheus.GaugeValue, float64(stat.WrittenDocsPerSec), writtenOperation)
	ch <- prometheus.MustNewConstMetric(e.metrics.clusterQueriesPerSecond, prometheus.GaugeValue, float64(stat.QueriesPerSec))
}

func (e *DTRRethinkDBExporter) processServerStat(stat api.QueryEngine, ch chan<- prometheus.Metric, replica string) {
	if e.client.Debug {
		e.log.Debugf("func: %s :- stat received: '%v'", "processServerStat", stat)
	}

	ch <- prometheus.MustNewConstMetric(e.metrics.serverClientConnections, prometheus.GaugeValue, float64(stat.ClientConnections), replica)
	ch <- prometheus.MustNewConstMetric(e.metrics.serverClientsActive, prometheus.GaugeValue, float64(stat.ClientConnections), replica)

	ch <- prometheus.MustNewConstMetric(e.metrics.serverDocsPerSecond, prometheus.GaugeValue, stat.ReadDocsPerSec, replica, readOperation)
	ch <- prometheus.MustNewConstMetric(e.metrics.serverDocsTotal, prometheus.GaugeValue, float64(stat.ReadDocsTotal), replica, readOperation)
	ch <- prometheus.MustNewConstMetric(e.metrics.serverDocsPerSecond, prometheus.GaugeValue, stat.WrittenDocsPerSec, replica, writtenOperation)
	ch <- prometheus.MustNewConstMetric(e.metrics.serverDocsTotal, prometheus.GaugeValue, float64(stat.WrittenDocsTotal), replica, writtenOperation)

	ch <- prometheus.MustNewConstMetric(e.metrics.serverQueriesPerSecond, prometheus.GaugeValue, stat.ReadDocsPerSec, replica)
	ch <- prometheus.MustNewConstMetric(e.metrics.serverQueriesTotal, prometheus.GaugeValue, float64(stat.QueriesTotal), replica)
}

func (e *DTRRethinkDBExporter) processTableStat(stat api.QueryEngine, ch chan<- prometheus.Metric, db string, table string, replica string) {
	if e.client.Debug {
		e.log.Debugf("func: %s :- stat received: '%v'", "processTableStat", stat)
	}

	ch <- prometheus.MustNewConstMetric(e.metrics.tableDocsPerSecond, prometheus.GaugeValue, stat.ReadDocsPerSec, db, table, replica, readOperation)
	ch <- prometheus.MustNewConstMetric(e.metrics.tableDocsPerSecond, prometheus.GaugeValue, stat.WrittenDocsPerSec, db, table, replica, writtenOperation)
}

func (e *DTRRethinkDBExporter) processTableServerStat(statQE api.QueryEngine, statSE api.StorageEngine, ch chan<- prometheus.Metric, db string, table string, replica string) {
	if e.client.Debug {
		e.log.Debugf("func: %s :- statQE received: '%v', statSE received: '%v'", "processTableServerStat", statQE, statSE)
	}

	ch <- prometheus.MustNewConstMetric(e.metrics.tableReplicaDocsPerSecond, prometheus.GaugeValue, statQE.ReadDocsPerSec, db, table, replica, readOperation)
	ch <- prometheus.MustNewConstMetric(e.metrics.tableReplicaDocsPerSecond, prometheus.GaugeValue, statQE.WrittenDocsPerSec, db, table, replica, writtenOperation)

	ch <- prometheus.MustNewConstMetric(e.metrics.tableReplicaCacheBytes, prometheus.GaugeValue, float64(statSE.Cache.InUseBytes), db, table, replica)

	ch <- prometheus.MustNewConstMetric(e.metrics.tableReplicaIO, prometheus.GaugeValue, statSE.Disk.ReadBytesPerSec, db, table, replica, readOperation)
	ch <- prometheus.MustNewConstMetric(e.metrics.tableReplicaIO, prometheus.GaugeValue, statSE.Disk.WrittenBytesPerSec, db, table, replica, writtenOperation)

	ch <- prometheus.MustNewConstMetric(e.metrics.tableReplicaDataBytes, prometheus.GaugeValue, float64(statSE.Disk.SpaceUsage.DataBytes), db, table, replica)
	ch <- prometheus.MustNewConstMetric(e.metrics.tableReplicaGarbageBytes, prometheus.GaugeValue, float64(statSE.Disk.SpaceUsage.GarbageBytes), db, table, replica)
	ch <- prometheus.MustNewConstMetric(e.metrics.tableReplicaMetaDataBytes, prometheus.GaugeValue, float64(statSE.Disk.SpaceUsage.MetadataBytes), db, table, replica)
}
