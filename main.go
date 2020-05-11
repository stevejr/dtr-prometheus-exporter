package main

import (
	"fmt"
	"net/http"
	"time"
	"crypto/tls"
	// "io/ioutil"
	// "encoding/json"

	"github.com/namsral/flag"
	"github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	dtr "github.com/stevejr/dtr-prometheus-exporter/client"
	"github.com/stevejr/dtr-prometheus-exporter/collector"
  "github.com/stevejr/dtr-prometheus-exporter/dtrconnector"
  "github.com/stevejr/dtr-prometheus-exporter/config"
)

var (
  log = logrus.New()
  cfg config.Config
)

func readAndValidateConfig(result *config.Config) {

	flag.StringVar(&result.DTR.DTRAPIAddress, "connection-string", "", "DTR connection string")
	flag.UintVar(&result.Web.ListenPort, "port", 9580, "Port to expose scraping endpoint on")
	flag.DurationVar(&result.Scrape.Timeout, "timeout", time.Second*30, "Timeout for scrape")
	flag.BoolVar(&result.Log.Debug, "verbose", false, "Enable verbose logging")

	flag.StringVar(&result.DTR.Username, "dtr-username", "", "Username of DTR user")
	flag.StringVar(&result.DTR.Password, "dtr-password", "", "Password of DTR user")
	flag.BoolVar(&result.DTR.EnableTLS, "enable-tls", true, "Enable to use tls connection")
	flag.StringVar(&result.DTR.CAFile, "dtr-ca", "", "Path to CA certificate file for tls connection")
	flag.StringVar(&result.DTR.CertificateFile, "dtr-cert", "", "Path to certificate file for tls connection")
	flag.StringVar(&result.DTR.KeyFile, "dtr-key", "", "Path to key file for tls connection")
	// API related flags
	flag.UintVar(&result.API.JobCount, "job-count", 10, "Number of Job entries to retrieve from Jobs API")
	

	flag.Parse()

	if result.DTR.DTRAPIAddress == "" {
		log.Fatal("DTR connection string not provided")
	}
}

func configureRoutes() {
	var landingPage = []byte(`<html>
		<head><title>DTR exporter for Prometheus</title></head>
		<body>
		<h1>DTR exporter for Prometheus</h1>
		<p><a href='/metrics'>Metrics</a></p>
		</body>
		</html>
		`)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(landingPage) // nolint: errcheck
	})

	http.Handle("/metrics", promhttp.Handler())
}

func setupLogger(cfg config.Config) {
	if cfg.Log.Debug {
		log.Level = logrus.DebugLevel
	}
}

func startHTTPServer(cfg config.Config) {
  listenAddr := fmt.Sprintf(":%d", cfg.Web.ListenPort)
  // fmt.Printf("HTTP Server Listen Port s%\n", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}

func main() {
	var tlsConfig *tls.Config
  var err error
  var cfg config.Config
	
	readAndValidateConfig(&cfg)
	setupLogger(cfg)

	log.WithFields(logrus.Fields{
		"port":    cfg.Web.ListenPort,
		"timeout": cfg.Scrape.Timeout,
		"verbose": cfg.Log.Debug,
		"connection-string": cfg.DTR.DTRAPIAddress,
	}).Infof("DTR exporter configured")

	configureRoutes()

	if cfg.DTR.EnableTLS {
		tlsConfig, err = dtrconnector.PrepareTLSConfig(cfg.DTR.CAFile, cfg.DTR.CertificateFile, cfg.DTR.KeyFile)
		if err != nil {
			log.Fatal("failed to read tls credentials")
		}
	}

  // fmt.Println("Creating new DTR Client")
  client := dtr.New(cfg, tlsConfig)
  // fmt.Println("Creating new Prom Collector")
  coll := collector.New(client, log)
  // fmt.Println("Registring new Prom Collector")
	prometheus.MustRegister(coll)

  // fmt.Println("Starting new HTTP Server")
  startHTTPServer(cfg)
}