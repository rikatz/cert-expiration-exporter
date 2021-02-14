package metrics

import (
	"context"
	"fmt"
	"time"

	conf "cse/internal/pkg/config"
	interfaceSchema "cse/internal/pkg/schema"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"go.opencensus.io/stats"
	"go.opencensus.io/tag"
)

// CertMetrics method to grab expiration seconds
func CertMetrics() error {
	//TODO: change so it is not hardcoded
	certmanager, _ := interfaceSchema.GetCertClient("certmanager")

	certmanager, err := certmanager.New()
	if err != nil {
		return fmt.Errorf("Failed start new Cert Manager")
	}
	certificateList, err := certmanager.GetCertList(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("Failed to get cert list")
	}

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
		ctx, err := tag.New(context.Background(), tag.Insert(conf.KeyNamespaceKey, ns),
			tag.Insert(conf.KeyCertNameKey, name), tag.Insert(conf.KeyOwnerKey, annotations["owner"]))
		if err != nil {
			return fmt.Errorf("Error adding tags to metrics: %s", err)
		}
		secondsValid := float64(certificate.Status.NotAfter.Sub(time.Now()) / time.Second)
		stats.Record(ctx, conf.CertRemainingSeconds.M(secondsValid))
	}

	return nil
}

func init() {
}
