# dtr-prometheus-exporter

Docker Trusted Registry (DTR) metrics Prometheus exporter. Issues HTTPs calls to the DTR REST APIs and scrapes the data.

## Installation

### From source

You need to have a Go 1.10+ environment configured. Clone the repo (outside your `GOPATH`) and then:

```bash
go build -o dtr-prometheus-exporter 
```

To run the dtr-prometheus-exporter tool:

```bash
./dtr-prometheus-exporter \
--connection-string=[YOUR CONNECTION STRING] \
--dtr-ca=[YOUR DTR CA.PEM] \
--dtr-username=[YOUR DTR USERNAME] \
--dtr-password=[YOUR DTR PASSWORD] \
--enable-tls=[TRUE||FALSE]
```

>**NOTE:**  with the above only the DTR Root CA is required as DTR does not enforce Client authentication.


### Using Docker

To build the docker image from source:

```bash
docker image build \
--no-cache=true \
--build-arg BUILD_DATE=$(date -u +'%Y-%m-%dT%H:%M:%SZ') \
--build-arg VCS_REF=$(git rev-parse --short HEAD) \
--build-arg BUILD_VERSION=alpine \
-t dockerps/dtr-prometheus-exporter:alpine .
```

To run the docker image:

```bash
docker run \
-d \
-p 9580:9580 \
--mount type=bind,source=[YOUR DTR CERTS DIR],target=/dtrcerts,readonly \
-e CONNECTION_STRING=[YOUR CONNECTION STRING] \
-e DTR_CA=/dtrcerts/[YOUR DTR CA.PEM FILENAME] \
-e DTR_USERNAME=[YOUR DTR USERNAME] \
-e DTR_PASSWORD=[YOUR DTR PASSWORD] \
dockerps/dtr-prometheus-exporter:alpine
```

>**NOTE:**  with the above you have to bind mount your DTR certificates into the container so that the can be used by the app.

### Using Helm

To deploy the dtr-prometheus-exporter into a Kubernetes cluster a helm chart is available in the helm directory. You should install Helm as per the Helm documentation and then amend the values.yaml as required before deploying the chart.

You will need to place the contents of you DTR CA.pem into the file ./helm/files/dtr_ca.pem as this file will be created in a configmap which will then be mounted into the pod.

## Configuration

The exporter can be configured with commandline arguments, environment variables and a configuration file. For the details on how to format the configuration file, visit [namsral/flag](https://github.com/namsral/flag) repo.

|Flag|ENV variable|Default|Meaning|
|---|---|---|---|
|--connection-string|CONNECTION_STRING|_no default_|Connection string for DTR including protocol and port|
|--dtr-ca|DTR_CA|_no default_|The DTR Root CA certificate file|
|--dtr-cert|DTR_CERT|_no default_|The DTR Cert certificate file|
|--dtr-key|DTR_KEY|_no default_|The DTR Key certificate file|
|--dtr-username|DTR_USERNAME|_no default_|The DTR username|
|--dtr-password|DTR_PASSWORD|_no default_|The DTR password|
|--enable-tls|ENABLE_TLS|true|Enable TLS on HTTP Client connection to DTR|
|--port|PORT|9580|Port to expose scrape endpoint on|
|--timeout|TIMEOUT|10s|Timeout when scraping the Service Bus|
|--verbose|VERBOSE|false|Enable verbose logging|
|--job-count|JOB_COUNT|100|Number of results to retrieve from the Jobs API|

## Exported metrics

Below are the metrics that are exported. Where you see `...` this means the metric is repeated for each replica in the DTR cluster.

```bash
# HELP dtr_cluster_client_connections Total number of connections from the cluster
# TYPE dtr_cluster_client_connections gauge
dtr_cluster_client_connections 231
# HELP dtr_cluster_clients_active Total number of active clients in the cluster
# TYPE dtr_cluster_clients_active gauge
dtr_cluster_clients_active 181
# HELP dtr_cluster_docs_per_second Total number of reads and writes of documents per second from the cluster
# TYPE dtr_cluster_docs_per_second gauge
dtr_cluster_docs_per_second{operation="read"} 4
dtr_cluster_docs_per_second{operation="written"} 0
# HELP dtr_cluster_queries_per_second Total number of queries per second from the cluster
# TYPE dtr_cluster_queries_per_second gauge
dtr_cluster_queries_per_second 5
# HELP dtr_job_total Count of Job Status
# TYPE dtr_job_total gauge
dtr_job_total{action="poll_mirror",status="done"} 20
dtr_job_total{action="tag_prune",status="done"} 20
# HELP dtr_scrape_latency Latency of collecting scrape
# TYPE dtr_scrape_latency gauge
dtr_scrape_latency 0.5014301
# HELP dtr_server_client_connections Number of client connections to the server(replica)
# TYPE dtr_server_client_connections gauge
dtr_server_client_connections{replica="378171528545"} 77
dtr_server_client_connections{replica="595fe3ff470e"} 77
dtr_server_client_connections{replica="724936901731"} 77
# HELP dtr_server_clients_active Total number of active clients in the server(replica)
# TYPE dtr_server_clients_active gauge
dtr_server_clients_active{replica="378171528545"} 77
dtr_server_clients_active{replica="595fe3ff470e"} 77
dtr_server_clients_active{replica="724936901731"} 77
# HELP dtr_server_count Count of DTR servers(replicas)
# TYPE dtr_server_count gauge
dtr_server_count 3
# HELP dtr_server_docs_per_second Total number of reads and writes of documents per second from the server(replica)
# TYPE dtr_server_docs_per_second gauge
dtr_server_docs_per_second{operation="read",replica="378171528545"} 4
dtr_server_docs_per_second{operation="read",replica="595fe3ff470e"} 0
dtr_server_docs_per_second{operation="read",replica="724936901731"} 0
dtr_server_docs_per_second{operation="written",replica="378171528545"} 0
dtr_server_docs_per_second{operation="written",replica="595fe3ff470e"} 0
dtr_server_docs_per_second{operation="written",replica="724936901731"} 0
# HELP dtr_server_docs_total Total number of reads and writes of documents from the server(replica)
# TYPE dtr_server_docs_total gauge
dtr_server_docs_total{operation="read",replica="378171528545"} 39686
dtr_server_docs_total{operation="read",replica="595fe3ff470e"} 2
dtr_server_docs_total{operation="read",replica="724936901731"} 10
dtr_server_docs_total{operation="written",replica="378171528545"} 5184
dtr_server_docs_total{operation="written",replica="595fe3ff470e"} 5221
dtr_server_docs_total{operation="written",replica="724936901731"} 5155
# HELP dtr_server_health_count Count of healthy/unhealthy DTR servers(replicas)
# TYPE dtr_server_health_count gauge
dtr_server_health_count{health="healthy"} 3
dtr_server_health_count{health="unhealthy"} 0
# HELP dtr_server_queries_per_second Number of queries per second from the server(replica)
# TYPE dtr_server_queries_per_second gauge
dtr_server_queries_per_second{replica="378171528545"} 4
dtr_server_queries_per_second{replica="595fe3ff470e"} 0
dtr_server_queries_per_second{replica="724936901731"} 0
# HELP dtr_server_queries_total Number of total queries from the server(replica)
# TYPE dtr_server_queries_total gauge
dtr_server_queries_total{replica="378171528545"} 12730
dtr_server_queries_total{replica="595fe3ff470e"} 12750
dtr_server_queries_total{replica="724936901731"} 13230
# HELP dtr_table_count Count of DTR tables
# TYPE dtr_table_count gauge
dtr_table_count 29
# HELP dtr_table_docs_per_second Number of reads and writes of documents per second from the table
# TYPE dtr_table_docs_per_second gauge
dtr_table_docs_per_second{db="dtr2",operation="read",table="blob_links"} 0
...
# HELP dtr_table_server_state_count Count of healthy/unhealthy DTR servers for each table
# TYPE dtr_table_server_state_count gauge
dtr_table_server_state_count{db="dtr2",state="notready",table="blob_links"} 0
dtr_table_server_state_count{db="dtr2",state="ready",table="blob_links"} 3
...
# HELP dtr_tablereplica_cache_bytes Table replica cache size in bytes
# TYPE dtr_tablereplica_cache_bytes gauge
dtr_tablereplica_cache_bytes{db="dtr2",replica="378171528545",table="blob_links"} 45376
...
# HELP dtr_tablereplica_data_bytes Table replica size in stored bytes
# TYPE dtr_tablereplica_data_bytes gauge
dtr_tablereplica_data_bytes{db="dtr2",replica="378171528545",table="blob_links"} 2.097152e+06
...
# HELP dtr_tablereplica_docs_per_second Number of reads and writes of documents per second from the table replica
# TYPE dtr_tablereplica_docs_per_second gauge
dtr_tablereplica_docs_per_second{db="dtr2",operation="read",replica="378171528545",table="blob_links"} 0
dtr_tablereplica_docs_per_second{db="dtr2",operation="written",replica="378171528545",table="blob_links"} 0
...
# HELP dtr_tablereplica_garbage_bytes Table replica garbage size in stored bytes
# TYPE dtr_tablereplica_garbage_bytes gauge
dtr_tablereplica_garbage_bytes{db="dtr2",replica="378171528545",table="blob_links"} 0
...
# HELP dtr_tablereplica_io Table replica reads and writes of bytes per second
# TYPE dtr_tablereplica_io gauge
dtr_tablereplica_io{db="dtr2",operation="read",replica="378171528545",table="blob_links"} 0
dtr_tablereplica_io{db="dtr2",operation="written",replica="378171528545",table="blob_links"} 0
...
# HELP dtr_tablereplica_metadata_bytes Table replica metadata size in stored bytes
# TYPE dtr_tablereplica_metadata_bytes gauge
dtr_tablereplica_metadata_bytes{db="dtr2",replica="378171528545",table="blob_links"} 8.388608e+06
...
# HELP dtr_tablereplica_preallocated_bytes Table replica preallocated size in stored bytes
# TYPE dtr_tablereplica_preallocated_bytes gauge
dtr_tablereplica_preallocated_bytes{db="dtr2",replica="378171528545",table="blob_links"} 2.097152e+06
...
# HELP dtr_up Whether the DTR scrape was successful
# TYPE dtr_up gauge
dtr_up 1
```
