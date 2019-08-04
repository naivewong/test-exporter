package main

import (
	"net/http"
	"os"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

//Define the metrics we wish to expose
var test1Metric = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "test1_metric", Help: "Test 1"})

var test2Metric = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "test2_metric", Help: "Test 2"})


func main() {
	//Register metrics with prometheus
	prometheus.MustRegister(test1Metric)
	prometheus.MustRegister(test2Metric)

	test1Metric.Set(0)
	test2Metric.Set(0)

	ticker := time.NewTicker(2 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <- ticker.C:
				test1Metric.Add(float64(2))
				test2Metric.Add(float64(2))
				case <- quit:
				ticker.Stop()
				return
			}
		}
	}()

	//This section will start the HTTP server and expose
	//any metrics on the /metrics endpoint.
	http.Handle("/metrics", promhttp.Handler())
	logger := level.NewFilter(log.NewLogfmtLogger(os.Stdout), level.AllowInfo())
	level.Info(logger).Log("msg", "Beginning to serve on port :8080")
	http.ListenAndServe(":8080", nil)
}