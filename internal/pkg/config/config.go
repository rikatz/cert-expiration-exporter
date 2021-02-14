package config

import (
	"os"
)

// Exporter to be used (prometheus or stackdriver)
var Exporter string

// KubeconfigPath to use (can be nil to use default)
var KubeconfigPath string

// InCluster to decide whether to use kubeconfig or service acc
var InCluster string

// KeyNamespace tag
var KeyNamespace string

// KeyCertName tag
var KeyCertName string

// KeyOwner tag
var KeyOwner string

// RefreshInterval for sending stats
var RefreshInterval string

func init() {
	Exporter = os.Getenv("EXPORTER")
	KubeconfigPath = os.Getenv("KUBECONFIG_PATH")
	InCluster = os.Getenv("IN_CLUSTER")
	KeyNamespace = os.Getenv("KEY_NAMESPACE")
	KeyCertName = os.Getenv("KEY_CERTNAME")
	KeyOwner = os.Getenv("KEY_OWNER")
	RefreshInterval = os.Getenv("REFRESH_INTERVAL")
}
