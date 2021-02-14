package config

import (
	"log"
	"os"
	"strconv"

	"go.opencensus.io/tag"
)

// Exporter to be used (prometheus or stackdriver)
var Exporter string

// KubeconfigPath to use (can be nil to use default)
var KubeconfigPath string

// InCluster to decide whether to use kubeconfig or service acc
var InCluster bool

// KeyNamespace tag
var KeyNamespace string

// KeyCertName tag
var KeyCertName string

// KeyOwner tag
var KeyOwner string

// RefreshInterval for sending stats
var RefreshInterval string

// KeyNamespaceKey tag key
var KeyNamespaceKey tag.Key

// KeyCertNameKey tag key
var KeyCertNameKey tag.Key

// KeyOwnerKey tag key
var KeyOwnerKey tag.Key

var err error

func init() {
	Exporter = os.Getenv("EXPORTER")
	KubeconfigPath = os.Getenv("KUBECONFIG_PATH")
	InCluster, _ = strconv.ParseBool(os.Getenv("IN_CLUSTER"))
	KeyNamespace = os.Getenv("KEY_NAMESPACE")
	KeyCertName = os.Getenv("KEY_CERTNAME")
	KeyOwner = os.Getenv("KEY_OWNER")
	RefreshInterval = os.Getenv("REFRESH_INTERVAL")

	KeyNamespaceKey, err = tag.NewKey("namespace")
	if err != nil {
		log.Fatalf("Failed to generate key namespace: %s", err)
	}

	KeyCertNameKey, err = tag.NewKey("certname")
	if err != nil {
		log.Fatalf("Failed to generate key certname: %s", err)
	}

	KeyOwnerKey, err = tag.NewKey("owner")
	if err != nil {
		log.Fatalf("Failed to generate key owner: %s", err)
	}
}
