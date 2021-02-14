package config

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"go.opencensus.io/stats"
	"go.opencensus.io/tag"

	"k8s.io/client-go/util/homedir"
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

// CertRemainingSeconds expiration time for certiicates
var CertRemainingSeconds *stats.Float64Measure

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

	if KubeconfigPath == "" {
		if home := homedir.HomeDir(); home != "" {
			flag.StringVar(&KubeconfigPath, "kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			flag.StringVar(&KubeconfigPath, "kubeconfig", "", "absolute path to the kubeconfig file")
		}

	}
	CertRemainingSeconds = stats.Float64("certificate_expiration", "The remaining seconds of a certificate before expiring", "s")
}
