package main

import (
	"flag"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"
	"strconv"
	"time"
)

// boolToFloat64 converts a bool to a float64
func boolToFloat64(b bool) float64 {
	if b {
		return float64(1)
	}
	return float64(0)
}

// startMonitor refreshes the Sia metrics periodically as defined by refreshRate
func startMonitor(refreshRate time.Duration, passwd string, address string) {
	for range time.Tick(time.Minute * refreshRate) {
		updateMetrics(passwd, address)
	}
}

// updateMetrics calls the various metric collection functions
func updateMetrics(passwd string, address string) {
	//do something every timeRefresh

	//call collector's function for curl values
	callClient(passwd, address)
}

func main() {
	// TEST VARIABLES
	port := flag.Int("port", 8101, "Port to serve Prometheus Metrics on")
	refresh := flag.Int("refresh", 1, "Frequency to get Metrics from Hostd (minutes)")
	passwd := flag.String("passwd", "Sia is Awesome", "Hostd API password")
	address := flag.String("address", "127.0.0.1:9980", "Hostd API address")

	flag.Parse()

	passwdEnv, isSet := os.LookupEnv("HOSTD_PASSWD")
	if isSet {
		*passwd = passwdEnv
	}

	// Set the metrics initially before starting the monitor and HTTP server
	// If you don't do this all the metrics start with a "0" until they are set
	updateMetrics(*passwd, *address)

	// start the metrics collector
	go startMonitor(time.Duration(*refresh), *passwd, *address)

	// This section will start the HTTP server and expose
	// any metrics on the /metrics endpoint.
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":"+strconv.Itoa(*port), nil)
}
