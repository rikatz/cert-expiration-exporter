package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"contrib.go.opencensus.io/exporter/prometheus"
	"contrib.go.opencensus.io/exporter/stackdriver"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"

	cmclient "github.com/jetstack/cert-manager/pkg/client/clientset/versioned"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

var (
	client cmclient.Interface

	kubeconfig string
	inCluster  bool
	err        error

	KeyNamespace, KeyCertName, KeyOwner tag.Key
	LineCountView                       *view.View

	certRemainingSeconds = stats.Float64("certificate_expiration", "The remaining seconds of a certificate before expiring", "s")
)

func main() {

	// TODO: Blah, this shouldn't be at main ;P And the Tag should be a const.
	KeyNamespace, err = tag.NewKey("namespace")
	if err != nil {
		log.Fatalf("Failed to generate key namespace: %s", err)
	}

	KeyCertName, err = tag.NewKey("certname")
	if err != nil {
		log.Fatalf("Failed to generate key certname: %s", err)
	}

	KeyOwner, err = tag.NewKey("owner")
	if err != nil {
		log.Fatalf("Failed to generate key owner: %s", err)
	}

	LineCountView = &view.View{
		Name:        "certificate_expiration_remaining_seconds",
		Measure:     certRemainingSeconds,
		Description: "The number of lines from standard input",
		Aggregation: view.LastValue(),
		TagKeys:     []tag.Key{KeyNamespace, KeyCertName, KeyOwner},
	}

	flag.Parse()
	var config *rest.Config
	var err error
	if !inCluster {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	} else {
		config, err = rest.InClusterConfig()
	}

	if err != nil {
		log.Fatalf("Unable to generate Kubernetes Configuration: %s", err)
	}
	client, err = cmclient.NewForConfig(config)
	if err != nil {
		log.Fatalf("Unable to connect to Kubernetes Cluster: %s", err)
	}

	if err := view.Register(LineCountView); err != nil {
		log.Fatalf("Failed to register the views: %v", err)
	}

	// TODO: Define each exporter in a separate file/package
	exporter, err := stackdriver.NewExporter(stackdriver.Options{})
	if err != nil {
		log.Fatal(err)
	}

	defer exporter.Flush()

	pe, err := prometheus.NewExporter(prometheus.Options{})
	if err != nil {
		log.Fatalf("Failed to create the Prometheus stats exporter: %v", err)
	}

	if err := exporter.StartMetricsExporter(); err != nil {
		log.Fatalf("Error starting metric exporter: %v", err)
	}
	defer exporter.StopMetricsExporter()

	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", pe)
		if err := http.ListenAndServe(":8888", mux); err != nil {
			log.Fatalf("Failed to run Prometheus scrape endpoint: %v", err)
		}
	}()

	for {
		err = CertMetrics()
		if err != nil {
			log.Fatalf("Error sending certificate metrics")
		}
		// TODO: Make the scrape interval configurable
		time.Sleep(61 * time.Second)
	}

}

// TODO: Move this function to other package and create unit tests
func CertMetrics() error {

	certificateList, err := client.CertmanagerV1().Certificates("").List(context.TODO(), metav1.ListOptions{})

	if apierrors.IsForbidden(err) {
		return fmt.Errorf("Permission denied while getting the certificates")
	}

	if err != nil {
		return fmt.Errorf("Failed to communicate with Kubernetes cluster")
	}

	if len(certificateList.Items) < 1 {
		return nil
	}

	for _, certificate := range certificateList.Items {
		ns, name := certificate.GetNamespace(), certificate.GetName()
		fmt.Println(name, ns)
		annotations := certificate.GetAnnotations()

		// If Owner is empty, should fill with something, otherwise
		// it generates two metrics on stackdriver.
		// also should have const with label definition
		ctx, err := tag.New(context.Background(), tag.Insert(KeyNamespace, ns),
			tag.Insert(KeyCertName, name), tag.Insert(KeyOwner, annotations["owner"]))
		if err != nil {
			return fmt.Errorf("Error adding tags to metrics: %s", err)
		}
		secondsValid := float64(certificate.Status.NotAfter.Sub(time.Now()) / time.Second)
		stats.Record(ctx, certRemainingSeconds.M(secondsValid))
	}

	return nil
}

func init() {

	if home := homedir.HomeDir(); home != "" {
		flag.StringVar(&kubeconfig, "kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		flag.StringVar(&kubeconfig, "kubeconfig", "", "absolute path to the kubeconfig file")
	}
}
