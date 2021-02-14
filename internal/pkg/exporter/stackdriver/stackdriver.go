package stackdriver

import (
	"context"
	"cse/internal/pkg/exporter"
	"cse/internal/pkg/schema"
	"log"

	"contrib.go.opencensus.io/exporter/stackdriver"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
)

// StackDriver exporter implementation
type StackDriver struct {
}

// Record new stat
func (sd *StackDriver) Record(ctx context.Context, stat *stats.Float64Measure) {

}

// Register new stat
func (sd *StackDriver) Register(lineCountView *view.View) (exporter.Exporter, error) {
	if err := view.Register(lineCountView); err != nil {
		log.Fatalf("Failed to register the views: %v", err)
	}

	// TODO: Define each exporter in a separate file/package
	exporter, err := stackdriver.NewExporter(stackdriver.Options{})
	if err != nil {
		log.Fatal(err)
	}

	defer exporter.Flush()
	if err := exporter.StartMetricsExporter(); err != nil {
		log.Fatalf("Error starting metric exporter: %v", err)
	}
	defer exporter.StopMetricsExporter()
	return sd, err
}

// Start StackDriver exporter
func (sd *StackDriver) Start() {

}

func init() {
	schema.RegisterExporter(&StackDriver{}, "stackdriver")
}
