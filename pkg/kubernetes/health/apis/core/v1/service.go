package v1

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/kubectl/pkg/scheme"

	"github.com/cakehappens/seaworthy/pkg/kubernetes/health"
)

// ServiceHealth checks the health of a Service
func ServiceHealth(obj *unstructured.Unstructured) (health.Status, error) {
	service := &corev1.Service{}
	err := scheme.Scheme.Convert(obj, service, nil)
	if err != nil {
		err = fmt.Errorf("failed to convert %T to %T: %w", obj, service, err)
		return health.Status{
			Code:    health.Unknown,
			Message: err.Error(),
		}, err
	}

	if service.Spec.Type == corev1.ServiceTypeLoadBalancer {
		if len(service.Status.LoadBalancer.Ingress) > 0 {
			return health.NewHealthyHealthStatus(), nil
		}

		return health.Status{
			Code: health.Progressing,
		}, nil
	}
	return health.Status{
		Code: health.Unknown,
	}, nil
}
