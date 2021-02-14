// Package schema probably need a better name for this :(
package schema

import (
	certclient "cse/internal/pkg/cert-client"
	"cse/internal/pkg/exporter"
	"fmt"
	"sync"
)

var exporterBuilder map[string]exporter.Exporter
var exporterBuildlock sync.RWMutex

var certClientBuilder map[string]certclient.CertClient
var certClientBuildlock sync.RWMutex

func init() {
	exporterBuilder = make(map[string]exporter.Exporter)
	certClientBuilder = make(map[string]certclient.CertClient)
}

// RegisterExporter registers a exporter
func RegisterExporter(s exporter.Exporter, exporterName string) {
	exporterBuildlock.Lock()
	defer exporterBuildlock.Unlock()

	_, exists := exporterBuilder[exporterName]
	if exists {
		panic(fmt.Sprintf("exporter %q already registered", exporterName))
	}

	exporterBuilder[exporterName] = s
}

// GetExporter returns the registered exporter
func GetExporter(exporterName string) (exporter.Exporter, error) {
	exporterBuildlock.RLock()
	f, ok := exporterBuilder[exporterName]
	exporterBuildlock.RUnlock()

	if !ok {
		return nil, fmt.Errorf("failed to find registered exporter name: %s", exporterName)
	}
	return f, nil
}

// RegisterCertClient registers a CertClient
func RegisterCertClient(s certclient.CertClient, certClientName string) {
	certClientBuildlock.Lock()
	defer certClientBuildlock.Unlock()

	_, exists := certClientBuilder[certClientName]
	if exists {
		panic(fmt.Sprintf("Cert Client %q already registered", certClientName))
	}

	certClientBuilder[certClientName] = s
}

// GetCertClient returns the registered Cert Client
func GetCertClient(certClientName string) (certclient.CertClient, error) {
	certClientBuildlock.RLock()
	f, ok := certClientBuilder[certClientName]
	certClientBuildlock.RUnlock()

	if !ok {
		return nil, fmt.Errorf("failed to find registered Cert Client name: %s", certClientName)
	}
	return f, nil
}
