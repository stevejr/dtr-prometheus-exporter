# dtr-prometheus-exporter

Docker Trusted Registry (DTR) metrics Prometheus exporter. Issues HTTPs calls to the DTR REST APIs and scrapes the data.

## Installation

### From source

You need to have a Go 1.10+ environment configured. Clone the repo (outside your `GOPATH`) and then:

```bash
go build -o dtr-prometheus-exporter
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
stevejr/dtr-prometheus-exporter:alpine
```

>**NOTE:**  with the above you have to bind mount your DTR certificates into the container so that the can be used by the app.

## Configuration

The exporter can be configured with commandline arguments, environment variables and a configuration file. For the details on how to format the configuration file, visit [namsral/flag](https://github.com/namsral/flag) repo.

|Flag|ENV variable|Default|Meaning|
|---|---|---|---|
|--connection-string|CONNECTION_STRING|_no default_|Connection string for DTR.|
|--dtr-ca|DTR_CA|_no default_|The DTR CA certificate file|
|--dtr-cert|DTR_CERT|_no default_|The DTR Cert certificate file|
|--dtr-key|DTR_KEY|_no default_|The DTR Key certificate file|
|--dtr-username|DTR_USERNAME|_no default_|The DTR username|
|--dtr-password|DTR_PASSWORD|_no default_|The DTR password|
|--enable-tls|ENABLE_TLS|true|Enable TLS on HTTP Client connection to DTR|
|--port|PORT|9580|Port to expose scrape endpoint on|
|--timeout|TIMEOUT|30s|Timeout when scraping the Service Bus|
|--verbose|VERBOSE|false|Enable verbose logging|
|--job-count|JOB_COUNT|10|Number of results to retrieve from the Jobs API|

## Exported metrics

```bash
 HELP dtr_job_total DTR job total
# TYPE dtr_job_total gauge
dtr_job_total{action="license_update",status="done"} 1
dtr_job_total{action="poll_mirror",status="done"} 49
dtr_job_total{action="tag_prune",status="done"} 50
# HELP dtr_replica_health_total DTR Replica Health count
# TYPE dtr_replica_health_total gauge
dtr_replica_health_total{health="healthy"} 3
dtr_replica_health_total{health="unhealthy"} 0
# HELP dtr_up Whether the DTR scrape was successful
# TYPE dtr_up gauge
dtr_up 1
```
