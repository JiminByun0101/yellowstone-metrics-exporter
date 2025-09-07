package metrics

import "github.com/prometheus/client_golang/prometheus"

// ExporterMetrics holds the few metrics we expose at Step 1.
// We'll add Solana-specific ones in later steps.
type ExporterMetrics struct {
	Up        prometheus.Gauge
	BuildInfo *prometheus.GaugeVec
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
	}
	em.Up.Set(1) // exporter started successfully
	return em
}

// Register all metrics with a registry and stamp build info.
func (em *ExporterMetrics) MustRegister(reg *prometheus.Registry, version, commit, date string) {
	reg.MustRegister(em.Up)
	reg.MustRegister(em.BuildInfo)
	em.BuildInfo.WithLabelValues(version, commit, date).Set(1)
}
