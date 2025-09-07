package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/JiminByun0101/yellowstone-metrics-exporter/internal/build"
	"github.com/JiminByun0101/yellowstone-metrics-exporter/internal/metrics"
)

func main() {
	addr := getenv("METRICS_ADDR", ":9108") // listen addr (can override with env var)

	// Create a new registry
	reg := prometheus.NewRegistry()
	// Add Go runtime and process metrics
	reg.MustRegister(collectors.NewGoCollector())
	reg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	// Add our exporter-specific metrics
	em := metrics.NewExporterMetrics()
	em.MustRegister(reg, build.Version, build.Commit, build.BuildDate)

	// Expose /metrics
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	srv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("exporter listening on %s (GET /metrics)\n", addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("http server error: %v", err)
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
