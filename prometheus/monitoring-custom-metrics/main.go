package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
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
		[]string{"core"})
	hwHumidity = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "hw_humidity",
			Help:       "The summary of humidity rate reported by a humidity sensor, as a ratio.",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"sensor"})
	tracer = otel.Tracer("github.com/kyma-project/examples/prometheus/monitoring-custom-metrics")
)

func init() {
	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(cpuTemp)
	prometheus.MustRegister(hdFailures)
	prometheus.MustRegister(cpuEnergy)
	prometheus.MustRegister(hwHumidity)
}

// randomTemp generates the temperature ranging from 60 to 90
func randomTemp() float64 {
	return math.Round(rand.Float64()*300)/10 + 60
}

// randomEnergy generates the energy ranging from 0 to 100
func randomEnergy() float64 {
	return math.Round(rand.Float64() * 100)
}

// randomHumidity generates the humidity from 0 to 1
func randomHumidity() float64 {
	return rand.Float64()
}

func newTraceProvider(exp sdktrace.SpanExporter) *sdktrace.TracerProvider {
	// Ensure default SDK resources and the required service name are set.
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("monitoring-custom-metrics"),
		),
	)

	if err != nil {
		panic(err)
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(r),
	)
}

func main() {
	client := otlptracehttp.NewClient()
	ctx := context.Background()
	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		panic(fmt.Errorf("creating OTLP trace exporter: %w", err))
	}
	// Create a new tracer provider with a batch span processor and the given exporter.
	tp := newTraceProvider(exporter)
	otel.SetTracerProvider(tp)

	// Handle shutdown properly so nothing leaks.
	defer func() { _ = tp.Shutdown(ctx) }()

	hdFailures.With(prometheus.Labels{"device": "/dev/sda"}).Inc()
	cpuEnergy.With(prometheus.Labels{"core": "0"}).Observe(randomEnergy())
	hwHumidity.With(prometheus.Labels{"sensor": "0"}).Observe(randomHumidity())

	// The Handler function provides a default handler to expose metrics
	handler := promhttp.Handler()

	// Wrap the handler with OpenTelemetry instrumentation
	wrappedHandler := otelhttp.NewHandler(handler, "/metrics")

	// via an HTTP server. "/metrics" is the usual endpoint for that.
	http.Handle("/metrics", wrappedHandler)

	http.DefaultClient.Timeout = 30 * time.Second

	//Port 8080 to be redirected to Envoy proxy
	log.Fatal(http.ListenAndServe(":8080", nil))

}
