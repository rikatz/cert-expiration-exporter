package certclient

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/rest"

	v1 "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1"
)

// CertClient to manage certificates
type CertClient interface {
	New(config *rest.Config) (CertClient, error)
	//TODO: generalize this for other clients later
	GetCertList(ctx context.Context, opts metav1.ListOptions) (*v1.CertificateList, error)
}
