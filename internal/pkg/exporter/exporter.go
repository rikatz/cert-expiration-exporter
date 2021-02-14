package exporter

import (
	"context"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
)

// Exporter is a common interface for interacting with monitoring services
type Exporter interface {
	New() (Exporter, error)
	Register(view *view.View)
	Record(ctx context.Context, stat *stats.Float64Measure)
}
