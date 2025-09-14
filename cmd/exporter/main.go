package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "time"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/collectors"
    "github.com/prometheus/client_golang/prometheus/promhttp"

    "github.com/jbyun0101/yellowstone-metrics-exporter/internal/build"
    "github.com/jbyun0101/yellowstone-metrics-exporter/internal/metrics"
    "github.com/jbyun0101/yellowstone-metrics-exporter/internal/stream"
)

func main() {
    addr := getenv("METRICS_ADDR", ":9108")

    reg := prometheus.NewRegistry()
    reg.MustRegister(collectors.NewGoCollector())
    reg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

    em := metrics.NewExporterMetrics()
    em.MustRegister(reg, build.Version, build.Commit, build.BuildDate)

    ctx := context.Background()
    go func() {
        client, err := stream.Dial("localhost:10000")
        if err != nil {
            log.Fatalf("failed to connect to Dragon's Mouth: %v", err)
        }
        defer client.Close()

        err = client.StreamSlots(ctx, func(slot uint64) {
            em.LatestSlot.Set(float64(slot))
            log.Printf("slot=%d", slot)
        })
        if err != nil {
            log.Fatalf("stream error: %v", err)
        }
    }()

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
