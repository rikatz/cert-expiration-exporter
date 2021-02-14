package prometheus

import (
	"context"
	"cse/internal/pkg/exporter"
	"cse/internal/pkg/schema"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
)

// Prometheus exporter implementation
type Prometheus struct {
}

// New returns a new StackDriver exporter implementation
func (sd *Prometheus) New() (exporter.Exporter, error) {

	return sd, nil
}

// Record new stat
func (sd *Prometheus) Record(ctx context.Context, stat *stats.Float64Measure) {

}

// Register new stat
func (sd *Prometheus) Register(view *view.View) {

}

func init() {
	schema.RegisterExporter(&Prometheus{}, "prometheus")
}
