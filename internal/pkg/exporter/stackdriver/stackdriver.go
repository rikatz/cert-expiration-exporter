package stackdriver

import (
	"context"
	"cse/internal/pkg/exporter"
	"cse/internal/pkg/schema"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
)

// StackDriver exporter implementation
type StackDriver struct {
}

// New returns a new StackDriver exporter implementation
func (sd *StackDriver) New() (exporter.Exporter, error) {

	return sd, nil
}

// Record new stat
func (sd *StackDriver) Record(ctx context.Context, stat *stats.Float64Measure) {

}

// Register new stat
func (sd *StackDriver) Register(view *view.View) {

}

func init() {
	schema.RegisterExporter(&StackDriver{}, "stackdriver")
}
