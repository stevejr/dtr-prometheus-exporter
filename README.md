# dtr-prometheus-exporter

Docker Trusted Registry (DTR) metrics Prometheus exporter. Issues HTTPs calls to the DTR REST APIs and scrapes the data.

## Installation

### From source

You need to have a Go 1.10+ environment configured. Clone the repo (outside your `GOPATH`) and then:

```bash
go build -o dtr-prometheus-exporter && 
./dtr-prometheus-exporter \
--connection-string=[YOUR CONNECTION STRING] \
--dtr-ca=[YOUR DTR CA.PEM] \
--dtr-cert=[YOUR DTR CERT.PEM] \
--dtr-key=[YOUR DTR KEY.PEM] \
--dtr-username=[YOUR DTR USERNAME] \
--dtr-password=[YOUR DTR PASSWORD] \
--enable-tls=[TRUE||FALSE]
```

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
-e DTR_CA=/dtrcerts/[YOUR CA.PEM FILENAME] \
-e DTR_CERT=/dtrcerts/[YOUR CERT.PEM FILENAME] \
-e DTR_KEY=/dtrcerts/[YOUR KEY.PEM FILENAME] \
-e DTR_USERNAME=[YOUR DTR USERNAME] \
-e DTR_PASSWORD=[YOUR DTR PASSWORD] \
dockerps/dtr-prometheus-exporter:alpine
```

>**NOTE:**  with the above you have to bind mount your DTR certificates into the container so that the can be used by the app.

### Using Helm

To deploy the dtr-prometheus-exporter into a Kubernetes cluster a helm chart is available in the helm directory. You should install Helm as per the Helm documentation and then amend the values.yaml as required before deploying the chart.

## Configuration

The exporter can be configured with commandline arguments, environment variables and a configuration file. For the details on how to format the configuration file, visit [namsral/flag](https://github.com/namsral/flag) repo.

|Flag|ENV variable|Default|Meaning|
|---|---|---|---|
|--connection-string|CONNECTION_STRING|_no default_|Connection string for DTR including protocol and port|
|--dtr-ca|DTR_CA|_no default_|The DTR CA certificate file|
|--dtr-cert|DTR_CERT|_no default_|The DTR Cert certificate file|
|--dtr-key|DTR_KEY|_no default_|The DTR Key certificate file|
|--dtr-username|DTR_USERNAME|_no default_|The DTR username|
|--dtr-password|DTR_PASSWORD|_no default_|The DTR password|
|--enable-tls|ENABLE_TLS|true|Enable TLS on HTTP Client connection to DTR|
|--port|PORT|9580|Port to expose scrape endpoint on|
|--timeout|TIMEOUT|30s|Timeout when scraping the Service Bus|
|--verbose|VERBOSE|false|Enable verbose logging|
|--job-count|JOB_COUNT|100|Number of results to retrieve from the Jobs API|

## Exported metrics

```bash
# HELP dtr_job_total DTR job total
# TYPE dtr_job_total gauge
dtr_job_total{action="license_update",status="done"} 1
dtr_job_total{action="poll_mirror",status="done"} 49
dtr_job_total{action="tag_prune",status="done"} 50
# HELP dtr_replica_health_total DTR replica health total
# TYPE dtr_replica_health_total gauge
dtr_replica_health_total{health="healthy"} 3
dtr_replica_health_total{health="unhealthy"} 0
# HELP dtr_table_disk_preallocated_bytes DTR table disk space preallocated bytes per replica
# TYPE dtr_table_disk_preallocated_bytes gauge
dtr_table_disk_preallocated_bytes{db="",replica="",table="",type="cluster"} 0
dtr_table_disk_preallocated_bytes{db="",replica="24133d3caf7c",table="",type="server"} 0
dtr_table_disk_preallocated_bytes{db="",replica="86ac5ff6e381",table="",type="server"} 0
dtr_table_disk_preallocated_bytes{db="",replica="f009cd362354",table="",type="server"} 0
dtr_table_disk_preallocated_bytes{db="dtr2",replica="24133d3caf7c",table="blob_links",type="table_server"} 2.097152e+06# HELP dtr_table_disk_read_bytes_seconds DTR table disk read bytes per second per replica
.....
# HELP dtr_table_disk_read_bytes_seconds DTR table disk read bytes per second per replica
# TYPE dtr_table_disk_read_bytes_seconds gauge
dtr_table_disk_read_bytes_seconds{db="",replica="",table="",type="cluster"} 0
dtr_table_disk_read_bytes_seconds{db="",replica="24133d3caf7c",table="",type="server"} 0
dtr_table_disk_read_bytes_seconds{db="",replica="86ac5ff6e381",table="",type="server"} 0
dtr_table_disk_read_bytes_seconds{db="",replica="f009cd362354",table="",type="server"} 0
dtr_table_disk_read_bytes_seconds{db="dtr2",replica="24133d3caf7c",table="blob_links",type="table_server"} 0
.....
# HELP dtr_table_disk_read_total_bytes DTR table disk read bytes total per replica
# TYPE dtr_table_disk_read_total_bytes gauge
dtr_table_disk_read_total_bytes{db="",replica="",table="",type="cluster"} 0
dtr_table_disk_read_total_bytes{db="",replica="24133d3caf7c",table="",type="server"} 0
dtr_table_disk_read_total_bytes{db="",replica="86ac5ff6e381",table="",type="server"} 0
dtr_table_disk_read_total_bytes{db="",replica="f009cd362354",table="",type="server"} 0
dtr_table_disk_read_total_bytes{db="dtr2",replica="24133d3caf7c",table="blob_links",type="table_server"} 6.664704e+06
.....
# HELP dtr_table_disk_used_data_bytes DTR table disk space used data bytes per replica
# TYPE dtr_table_disk_used_data_bytes gauge
dtr_table_disk_used_data_bytes{db="",replica="",table="",type="cluster"} 0
dtr_table_disk_used_data_bytes{db="",replica="24133d3caf7c",table="",type="server"} 0
dtr_table_disk_used_data_bytes{db="",replica="86ac5ff6e381",table="",type="server"} 0
dtr_table_disk_used_data_bytes{db="",replica="f009cd362354",table="",type="server"} 0
dtr_table_disk_used_data_bytes{db="dtr2",replica="24133d3caf7c",table="blob_links",type="table_server"} 2.097152e+06
.....
# HELP dtr_table_disk_used_garbage_bytes DTR table disk space used garbage bytes per replica
# TYPE dtr_table_disk_used_garbage_bytes gauge
dtr_table_disk_used_garbage_bytes{db="",replica="",table="",type="cluster"} 0
dtr_table_disk_used_garbage_bytes{db="",replica="24133d3caf7c",table="",type="server"} 0
dtr_table_disk_used_garbage_bytes{db="",replica="86ac5ff6e381",table="",type="server"} 0
dtr_table_disk_used_garbage_bytes{db="",replica="f009cd362354",table="",type="server"} 0
dtr_table_disk_used_garbage_bytes{db="dtr2",replica="24133d3caf7c",table="blob_links",type="table_server"} 0
.....
# HELP dtr_table_disk_used_metadata_bytes DTR table disk space use metadata bytes per replica
# TYPE dtr_table_disk_used_metadata_bytes gauge
dtr_table_disk_used_metadata_bytes{db="",replica="",table="",type="cluster"} 0
dtr_table_disk_used_metadata_bytes{db="",replica="24133d3caf7c",table="",type="server"} 0
dtr_table_disk_used_metadata_bytes{db="",replica="86ac5ff6e381",table="",type="server"} 0
dtr_table_disk_used_metadata_bytes{db="",replica="f009cd362354",table="",type="server"} 0
dtr_table_disk_used_metadata_bytes{db="dtr2",replica="24133d3caf7c",table="blob_links",type="table_server"} 8.388608e+06
.....
# HELP dtr_table_disk_written_bytes_second DTR table disk written bytes per second replica
# TYPE dtr_table_disk_written_bytes_second gauge
dtr_table_disk_written_bytes_second{db="",replica="",table="",type="cluster"} 0
dtr_table_disk_written_bytes_second{db="",replica="24133d3caf7c",table="",type="server"} 0
dtr_table_disk_written_bytes_second{db="",replica="86ac5ff6e381",table="",type="server"} 0
dtr_table_disk_written_bytes_second{db="",replica="f009cd362354",table="",type="server"} 0
dtr_table_disk_written_bytes_second{db="dtr2",replica="24133d3caf7c",table="blob_links",type="table_server"} 0
.....
# HELP dtr_table_disk_written_total_bytes DTR table disk written bytes total per replica
# TYPE dtr_table_disk_written_total_bytes gauge
dtr_table_disk_written_total_bytes{db="",replica="",table="",type="cluster"} 0
dtr_table_disk_written_total_bytes{db="",replica="24133d3caf7c",table="",type="server"} 0
dtr_table_disk_written_total_bytes{db="",replica="86ac5ff6e381",table="",type="server"} 0
dtr_table_disk_written_total_bytes{db="",replica="f009cd362354",table="",type="server"} 0
dtr_table_disk_written_total_bytes{db="dtr2",replica="24133d3caf7c",table="blob_links",type="table_server"} 73728
.....
# HELP dtr_table_queryengine_read_docs_seconds DTR table query engine read docs per second per replica
# TYPE dtr_table_queryengine_read_docs_seconds gauge
dtr_table_queryengine_read_docs_seconds{db="",replica="",table="",type="cluster"} 4
dtr_table_queryengine_read_docs_seconds{db="",replica="24133d3caf7c",table="",type="server"} 0
dtr_table_queryengine_read_docs_seconds{db="",replica="86ac5ff6e381",table="",type="server"} 0
dtr_table_queryengine_read_docs_seconds{db="",replica="f009cd362354",table="",type="server"} 4
dtr_table_queryengine_read_docs_seconds{db="dtr2",replica="24133d3caf7c",table="blob_links",type="table_server"} 0
.....
# HELP dtr_table_queryengine_read_docs_total DTR table query engine read docs total per replica
# TYPE dtr_table_queryengine_read_docs_total gauge
dtr_table_queryengine_read_docs_total{db="",replica="",table="",type="cluster"} 0
dtr_table_queryengine_read_docs_total{db="",replica="24133d3caf7c",table="",type="server"} 7
dtr_table_queryengine_read_docs_total{db="",replica="86ac5ff6e381",table="",type="server"} 12
dtr_table_queryengine_read_docs_total{db="",replica="f009cd362354",table="",type="server"} 679252
dtr_table_queryengine_read_docs_total{db="dtr2",replica="24133d3caf7c",table="blob_links",type="table_server"} 0
.....
# HELP dtr_table_queryengine_written_docs_seconds DTR table query engine written docs per second per replica
# TYPE dtr_table_queryengine_written_docs_seconds gauge
dtr_table_queryengine_written_docs_seconds{db="",replica="",table="",type="cluster"} 0
dtr_table_queryengine_written_docs_seconds{db="",replica="24133d3caf7c",table="",type="server"} 0
dtr_table_queryengine_written_docs_seconds{db="",replica="86ac5ff6e381",table="",type="server"} 0
dtr_table_queryengine_written_docs_seconds{db="",replica="f009cd362354",table="",type="server"} 0
dtr_table_queryengine_written_docs_seconds{db="dtr2",replica="24133d3caf7c",table="blob_links",type="table_server"} 0
.....
# HELP dtr_table_queryengine_written_docs_total DTR table query engine written docs total per replica
# TYPE dtr_table_queryengine_written_docs_total gauge
dtr_table_queryengine_written_docs_total{db="",replica="",table="",type="cluster"} 0
dtr_table_queryengine_written_docs_total{db="",replica="24133d3caf7c",table="",type="server"} 222565
dtr_table_queryengine_written_docs_total{db="",replica="86ac5ff6e381",table="",type="server"} 222603
dtr_table_queryengine_written_docs_total{db="",replica="f009cd362354",table="",type="server"} 222536
dtr_table_queryengine_written_docs_total{db="dtr2",replica="24133d3caf7c",table="blob_links",type="table_server"} 0
.....
# HELP dtr_table_replica_status_total DTR table replica status total
# TYPE dtr_table_replica_status_total gauge
dtr_table_replica_status_total{db="dtr2",status="notready",table="blob_links"} 0
# HELP dtr_up Whether the DTR scrape was successful
# TYPE dtr_up gauge
dtr_up 1
```
