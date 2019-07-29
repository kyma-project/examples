package main

import (
	"log"
	"math"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	cpuTemp = prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Name: "cpu_temperature_celsius",
		Help: "Current temperature of the CPU.",
	}, randomTemp)
	hdFailures = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "hd_errors_total",
			Help: "Number of hard-disk errors.",
		},
		[]string{"device"},
	)
)

func init() {
	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(cpuTemp)
	prometheus.MustRegister(hdFailures)
}

// randomTemp generates temperature in the range from 60 to 90
func randomTemp() float64 {
	return math.Round(rand.Float64()*300)/10 + 60
}

func main() {
	hdFailures.With(prometheus.Labels{"device": "/dev/sda"}).Inc()

	// The Handler function provides a default handler to expose metrics
	// via an HTTP server. "/metrics" is the usual endpoint for that.
	http.Handle("/metrics", promhttp.Handler())

	http.DefaultClient.Timeout = 30 * time.Second
	//Excluded port 8081, to be reached by Prometheus
	go func() {
		log.Fatal(http.ListenAndServe(":8081", nil))
	}()

	//Port 8080 to be redirected to Envoy proxy not reacheble by Prometheus
	log.Fatal(http.ListenAndServe(":8080", nil))

}
