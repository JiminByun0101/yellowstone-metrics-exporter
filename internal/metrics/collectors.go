package metrics

import "github.com/prometheus/client_golang/prometheus"

type ExporterMetrics struct {
    Up         prometheus.Gauge
    BuildInfo  *prometheus.GaugeVec
    LatestSlot prometheus.Gauge
}

func NewExporterMetrics() *ExporterMetrics {
    em := &ExporterMetrics{
        Up: prometheus.NewGauge(prometheus.GaugeOpts{
            Name: "solana_exporter_up",
            Help: "1 if exporter is running",
        }),
        BuildInfo: prometheus.NewGaugeVec(prometheus.GaugeOpts{
            Name: "solana_exporter_build_info",
            Help: "Build information for the exporter",
        }, []string{"version", "commit", "date"}),
        LatestSlot: prometheus.NewGauge(prometheus.GaugeOpts{
            Name: "solana_latest_slot",
            Help: "Most recent slot observed from Dragon's Mouth",
        }),
    }
    em.Up.Set(1)
    em.LatestSlot.Set(0) // initialize
    return em
}

func (em *ExporterMetrics) MustRegister(reg *prometheus.Registry, version, commit, date string) {
    reg.MustRegister(em.Up)
    reg.MustRegister(em.BuildInfo)
    reg.MustRegister(em.LatestSlot)
    em.BuildInfo.WithLabelValues(version, commit, date).Set(1)
}
