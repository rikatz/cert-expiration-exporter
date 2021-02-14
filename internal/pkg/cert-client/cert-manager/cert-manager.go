package certmanager

import (
	"context"
	certclient "cse/internal/pkg/cert-client"
	"cse/internal/pkg/schema"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1"

	"k8s.io/client-go/rest"
)

// CertManager CertClient implementation
type CertManager struct {
}

// New returns a new StackDriver exporter implementation
func (sd *CertManager) New(config *rest.Config) (certclient.CertClient, error) {

	return sd, nil
}

// GetCertList returns a list of Cert Manager Certificates
func (sd *CertManager) GetCertList(ctx context.Context, opts metav1.ListOptions) (*v1.CertificateList, error) {
	certList := &v1.CertificateList{}
	return certList, nil
}

func init() {
	schema.RegisterCertClient(&CertManager{}, "certmanager")
}
