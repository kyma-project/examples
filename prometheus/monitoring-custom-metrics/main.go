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
	cpuEnergy = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "cpu_energy_watt",
			Help:    "Current power usage reported by the CPU.",
			Buckets: prometheus.LinearBuckets(0, 20, 5),
		},
		[]string{})
)

func init() {
	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(cpuTemp)
	prometheus.MustRegister(hdFailures)
}

// randomTemp generates the temperature ranging from 60 to 90
func randomTemp() float64 {
	return math.Round(rand.Float64()*300)/10 + 60
}

// randomTemp generates the energy ranging from 0 to 100
func randomEnergy() float64 {
	return math.Round(rand.Float64() * 100)
}

func main() {
	hdFailures.With(prometheus.Labels{"device": "/dev/sda"}).Inc()

	cpuEnergy.WithLabelValues("core", "0").Observe(randomEnergy())

	// The Handler function provides a default handler to expose metrics
	// via an HTTP server. "/metrics" is the usual endpoint for that.
	http.Handle("/metrics", promhttp.Handler())

	http.DefaultClient.Timeout = 30 * time.Second

	//Port 8080 to be redirected to Envoy proxy
	log.Fatal(http.ListenAndServe(":8080", nil))

}
