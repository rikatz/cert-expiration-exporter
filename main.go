package main

import (
	certmanager "cse/internal/pkg/cert-client/cert-manager"
	conf "cse/internal/pkg/config"
	"cse/internal/pkg/exporter"
	"cse/internal/pkg/exporter/prometheus"
	"cse/internal/pkg/exporter/stackdriver"
	"cse/internal/pkg/metrics"
	interfaceSchema "cse/internal/pkg/schema"
	"fmt"
	"log"
	"strconv"
	"strings"
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

func init() {
	interfaceSchema.RegisterExporter(&stackdriver.StackDriver{}, "stackdriver")
	interfaceSchema.RegisterExporter(&prometheus.Prometheus{}, "prometheus")
	interfaceSchema.RegisterCertClient(&certmanager.CertManager{}, "certmanager")

}

func main() {

	LineCountView = &view.View{
		Name:        "certificate_expiration_remaining_seconds",
		Measure:     conf.CertRemainingSeconds,
		Description: "The number of lines from standard input",
		Aggregation: view.LastValue(),
		TagKeys:     []tag.Key{conf.KeyNamespaceKey, conf.KeyCertNameKey, conf.KeyOwnerKey},
	}

	exporterChoice := conf.Exporter
	log.Println("Chosen Exporters:", exporterChoice)
	exporters := strings.Split(exporterChoice, ",")
	for _, ex := range exporters {
		log.Println("Trying to run exporter", ex)
		chosenExporter, err := interfaceSchema.GetExporter(ex)
		if err != nil {
			log.Fatalf("Failed to select exporter: %v", err)
		}

		var exporterRegistered exporter.Exporter

		exporterRegistered, err = chosenExporter.Register(LineCountView)

		if err != nil {
			log.Fatalf("Failed to register the views: %v", err)
		}
		exporterRegistered.Start()
	}

	for {
		err = metrics.CertMetrics()
		if err != nil {

			fmt.Println(err.Error())
			log.Fatalf("Error sending certificate metrics")
		}

		var refreshInterval int
		if conf.RefreshInterval == "" {
			refreshInterval = 61
		} else {
			refreshInterval, _ = strconv.Atoi(conf.RefreshInterval)
		}
		time.Sleep(time.Duration(refreshInterval) * time.Second)
	}

}
