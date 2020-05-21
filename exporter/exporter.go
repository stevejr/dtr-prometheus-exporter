package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	dtr "github.com/stevejr/dtr-prometheus-exporter/client"
)

// DTRRethinkDBExporter - is a prometheus exporter of the rethinkdb statistics
type DTRRethinkDBExporter struct {
	client *dtr.DTRClient
	log    *logrus.Logger

	metrics struct {
		// General Metrics
		up            *prometheus.Desc
		scrapeLatency *prometheus.Desc

		// DTR Details from rethink_system_tables:stats json struct in cluster_status response
		clusterClientConnections *prometheus.Desc
		clusterClientsActive     *prometheus.Desc
		clusterDocsPerSecond     *prometheus.Desc
		clusterQueriesPerSecond  *prometheus.Desc

		serverClientConnections *prometheus.Desc
		serverClientsActive     *prometheus.Desc
		serverQueriesPerSecond  *prometheus.Desc
		serverQueriesTotal      *prometheus.Desc
		serverDocsPerSecond     *prometheus.Desc
		serverDocsTotal         *prometheus.Desc

		tableDocsPerSecond *prometheus.Desc
		tableRowsCount     *prometheus.Desc

		tableReplicaDocsPerSecond *prometheus.Desc
		tableReplicaCacheBytes    *prometheus.Desc
		tableReplicaIO            *prometheus.Desc
		tableReplicaDataBytes     *prometheus.Desc
		tableReplicaGarbageBytes  *prometheus.Desc
		tableReplicaMetaDataBytes *prometheus.Desc

		// DTR Details from replica_health json struct in cluster_status API response
		serverHealth *prometheus.Desc
		serverTotal  *prometheus.Desc

		// DTR Details from table_status json struct in cluster_status API response
		tableStatus *prometheus.Desc
		tableTotal  *prometheus.Desc

		// DTR Details from jobs API response
		jobTotals *prometheus.Desc
	}
}

// New - Creates new Prometheus collector
func New(client *dtr.DTRClient, log *logrus.Logger) *DTRRethinkDBExporter {
	exporter := &DTRRethinkDBExporter{
		client: client,
		log:    log,
	}

	exporter.initMetrics()

	return exporter
}
