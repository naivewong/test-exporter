package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	metrics := []prometheus.Gauge{}

	if len(os.Args) != 3 {
		panic("go run main.go #metrics port")
	}

	//Register metrics with prometheus
	numMetrics, _ := strconv.Atoi(os.Args[1])
	for i := 0; i < numMetrics; i++ {
		lbs := prometheus.Labels{}
		for j := 0; j < 20; j++ {
			lbs[fmt.Sprintf("label%d", j)] = fmt.Sprintf("label%d", i)
		}
		metrics = append(metrics, prometheus.NewGauge(prometheus.GaugeOpts{Name: fmt.Sprintf("test%d_metric", i), Help: fmt.Sprintf("Test %d", i), ConstLabels: lbs}))
		prometheus.MustRegister(metrics[len(metrics) - 1])
		metrics[len(metrics) - 1].Set(0)
	}

	ticker := time.NewTicker(2 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <- ticker.C:
				for i := 0; i < len(metrics); i++ {
					metrics[i].Add(float64(2))
				}
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
	level.Info(logger).Log("msg", "Beginning to serve on port :" + os.Args[2])
	http.ListenAndServe(":" + os.Args[2], nil)
}