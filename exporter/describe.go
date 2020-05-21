package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	readOperation    = "read"
	writtenOperation = "written"
)

type sumOfJobs struct {
	action string
	status string
}

// Describe sends metrics descriptions to the prometheus chan
func (e *DTRRethinkDBExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.metrics.clusterClientConnections
	ch <- e.metrics.clusterClientsActive
	ch <- e.metrics.clusterDocsPerSecond
	ch <- e.metrics.clusterQueriesPerSecond

	ch <- e.metrics.serverClientConnections
	ch <- e.metrics.serverClientsActive
	ch <- e.metrics.serverQueriesPerSecond
	ch <- e.metrics.serverQueriesTotal
	ch <- e.metrics.serverDocsPerSecond
	ch <- e.metrics.serverDocsTotal

	ch <- e.metrics.tableDocsPerSecond

	ch <- e.metrics.tableReplicaDocsPerSecond
	ch <- e.metrics.tableReplicaCacheBytes
	ch <- e.metrics.tableReplicaIO
	ch <- e.metrics.tableReplicaDataBytes
	ch <- e.metrics.tableReplicaGarbageBytes
	ch <- e.metrics.tableReplicaMetaDataBytes

	ch <- e.metrics.scrapeLatency

	ch <- e.metrics.serverHealth
	ch <- e.metrics.serverTotal

	ch <- e.metrics.tableStatus

	ch <- e.metrics.jobTotals
}

func (e *DTRRethinkDBExporter) initMetrics() {

	e.metrics.up = prometheus.NewDesc(
		"dtr_up",
		"Whether the DTR scrape was successful",
		nil, nil)

	e.metrics.clusterClientConnections = prometheus.NewDesc(
		"dtr_cluster_client_connections",
		"Total number of connections from the cluster",
		nil, nil)
	e.metrics.clusterClientsActive = prometheus.NewDesc(
		"dtr_cluster_clients_active",
		"Total number of active clients in the cluster",
		nil, nil)
	e.metrics.clusterDocsPerSecond = prometheus.NewDesc(
		"dtr_cluster_docs_per_second",
		"Total number of reads and writes of documents per second from the cluster",
		[]string{"operation"}, nil)
	e.metrics.clusterQueriesPerSecond = prometheus.NewDesc(
		"dtr_cluster_queries_per_second",
		"Total number of queries per second from the cluster",
		nil, nil)

	e.metrics.serverClientConnections = prometheus.NewDesc(
		"dtr_server_client_connections",
		"Number of client connections to the server(replica)",
		[]string{"replica"}, nil)
	e.metrics.serverClientsActive = prometheus.NewDesc(
		"dtr_server_clients_active",
		"Total number of active clients in the server(replica)",
		[]string{"replica"}, nil)
	e.metrics.serverQueriesPerSecond = prometheus.NewDesc(
		"dtr_server_queries_per_second",
		"Number of queries per second from the server(replica)",
		[]string{"replica"}, nil)
	e.metrics.serverQueriesTotal = prometheus.NewDesc(
		"dtr_server_queries_total",
		"Number of total queries from the server(replica)",
		[]string{"replica"}, nil)
	e.metrics.serverDocsPerSecond = prometheus.NewDesc(
		"dtr_server_docs_per_second",
		"Total number of reads and writes of documents per second from the server(replica)",
		[]string{"replica", "operation"}, nil)
	e.metrics.serverDocsTotal = prometheus.NewDesc(
		"dtr_server_docs_total",
		"Total number of reads and writes of documents from the server(replica)",
		[]string{"replica", "operation"}, nil)

	e.metrics.tableDocsPerSecond = prometheus.NewDesc(
		"dtr_table_docs_per_second",
		"Number of reads and writes of documents per second from the table",
		[]string{"db", "table", "replica", "operation"}, nil)

	e.metrics.tableReplicaDocsPerSecond = prometheus.NewDesc(
		"dtr_tablereplica_docs_per_second",
		"Number of reads and writes of documents per second from the table replica",
		[]string{"db", "table", "replica", "operation"}, nil)
	e.metrics.tableReplicaCacheBytes = prometheus.NewDesc(
		"dtr_tablereplica_cache_bytes",
		"Table replica cache size in bytes",
		[]string{"db", "table", "replica"}, nil)
	e.metrics.tableReplicaIO = prometheus.NewDesc(
		"dtr_tablereplica_io",
		"Table replica reads and writes of bytes per second",
		[]string{"db", "table", "replica", "operation"}, nil)
	e.metrics.tableReplicaDataBytes = prometheus.NewDesc(
		"dtr_tablereplica_data_bytes",
		"Table replica size in stored bytes",
		[]string{"db", "table", "replica"}, nil)
	e.metrics.tableReplicaGarbageBytes = prometheus.NewDesc(
		"dtr_tablereplica_garbage_bytes",
		"Table replica garbage size in stored bytes",
		[]string{"db", "table", "replica"}, nil)
	e.metrics.tableReplicaMetaDataBytes = prometheus.NewDesc(
		"dtr_tablereplica_metadata_bytes",
		"Table replica metadata size in stored bytes",
		[]string{"db", "table", "replica"}, nil)

	e.metrics.scrapeLatency = prometheus.NewDesc(
		"dtr_scrape_latency",
		"Latency of collecting scrape",
		nil, nil)

	e.metrics.serverHealth = prometheus.NewDesc(
		"dtr_server_health_count",
		"Count of healthy/unhealthy DTR servers(replicas)",
		[]string{"health"}, nil)
	e.metrics.serverTotal = prometheus.NewDesc(
		"dtr_server_count",
		"Count of DTR servers(replicas)",
		nil, nil)

	e.metrics.tableStatus = prometheus.NewDesc(
		"dtr_table_server_state_count",
		"Count of healthy/unhealthy DTR servers for each table",
		[]string{"db", "table", "state"}, nil)
	e.metrics.tableTotal = prometheus.NewDesc(
		"dtr_table_count",
		"Count of DTR tables",
		nil, nil)

	e.metrics.jobTotals = prometheus.NewDesc(
		"dtr_job_total",
		"Count of Job Status",
		[]string{"status", "action"}, nil)
}
