package prometheus

import (
	"context"
	"cse/internal/pkg/exporter"
	"cse/internal/pkg/schema"
	"log"
	"net/http"

	"contrib.go.opencensus.io/exporter/prometheus"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
)

// Prometheus exporter implementation
type Prometheus struct {
	PrometheusExporter *prometheus.Exporter
}

// Record new stat
func (pm *Prometheus) Record(ctx context.Context, stat *stats.Float64Measure) {

}

// Register new stat
func (pm *Prometheus) Register(lineCountView *view.View) (exporter.Exporter, error) {
	pe, err := prometheus.NewExporter(prometheus.Options{})
	if err != nil {
		log.Fatalf("Failed to create the Prometheus stats exporter: %v", err)
	}
	pm.PrometheusExporter = pe
	return pm, err
}

// Start prometheus exporter
func (pm *Prometheus) Start() {
	go func() {
		mux := http.NewServeMux()

		mux.Handle("/metrics", pm.PrometheusExporter)
		if err := http.ListenAndServe(":8888", mux); err != nil {
			log.Fatalf("Failed to run Prometheus scrape endpoint: %v", err)
		}
	}()
}

func init() {
	schema.RegisterExporter(&Prometheus{}, "prometheus")
}
