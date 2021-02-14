package main

import (
	conf "cse/internal/pkg/config"
	"cse/internal/pkg/exporter"
	"cse/internal/pkg/metrics"
	interfaceSchema "cse/internal/pkg/schema"
	"fmt"
	"log"
	"time"

	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

var (
	err error

	// LineCountView a view to pass to the exporter Register
	LineCountView *view.View
)

func main() {

	LineCountView = &view.View{
		Name:        "certificate_expiration_remaining_seconds",
		Measure:     conf.CertRemainingSeconds,
		Description: "The number of lines from standard input",
		Aggregation: view.LastValue(),
		TagKeys:     []tag.Key{conf.KeyNamespaceKey, conf.KeyCertNameKey, conf.KeyOwnerKey},
	}

	exporterChoice := conf.Exporter
	log.Println("Chosen Exporter:", exporterChoice)
	chosenExporter, err := interfaceSchema.GetExporter(exporterChoice)
	if err != nil {
		log.Fatalf("Failed to select exporter: %v", err)
	}
	chosenExporter.Start()

	var exporterRegistered exporter.Exporter
	exporterRegistered, err = chosenExporter.Register(LineCountView)
	if err != nil {
		log.Fatalf("Failed to register the views: %v", err)
	}

	exporterRegistered.Start()

	for {
		err = metrics.CertMetrics()
		if err != nil {

			fmt.Println(err.Error())
			log.Fatalf("Error sending certificate metrics")
		}
		// TODO: Make the scrape interval configurable
		time.Sleep(61 * time.Second)
	}

}
