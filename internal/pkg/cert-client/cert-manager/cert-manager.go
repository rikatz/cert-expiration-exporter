package certmanager

import (
	"context"
	certclient "cse/internal/pkg/cert-client"
	conf "cse/internal/pkg/config"
	"fmt"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1"

	"flag"

	cmclient "github.com/jetstack/cert-manager/pkg/client/clientset/versioned"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// CertManager CertClient implementation
type CertManager struct {
	Client *cmclient.Clientset
}

// New returns a new StackDriver exporter implementation
func (sd *CertManager) New() (certclient.CertClient, error) {

	flag.Parse()
	var err error
	var config *rest.Config
	if !conf.InCluster {
		config, err = clientcmd.BuildConfigFromFlags("", conf.KubeconfigPath)
	} else {
		config, err = rest.InClusterConfig()
	}

	if err != nil {
		log.Fatalf("Unable to generate Kubernetes Configuration: %s", err)
	}

	var client *cmclient.Clientset
	client, err = cmclient.NewForConfig(config)

	if err != nil {
		log.Fatalf("Unable to connect to Kubernetes Cluster: %s", err)
	}

	sd.Client = client

	return sd, nil
}

// GetCertList returns a list of Cert Manager Certificates
func (sd *CertManager) GetCertList(ctx context.Context, opts metav1.ListOptions) (*v1.CertificateList, error) {

	client := sd.Client
	certificateList, err := client.CertmanagerV1().Certificates("").List(context.TODO(), metav1.ListOptions{})
	if apierrors.IsForbidden(err) {
		return certificateList, fmt.Errorf("Permission denied while getting the certificates")
	}

	if err != nil {
		return certificateList, fmt.Errorf("Failed to communicate with Kubernetes cluster")
	}

	if len(certificateList.Items) < 1 {
		return certificateList, nil
	}
	return certificateList, nil
}
