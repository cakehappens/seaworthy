package kubernetes

import (
	"context"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type ResourceGetter interface {
	Get(ctx context.Context, options ...GetOption) ([]unstructured.Unstructured, error)
}

type Client interface {
	ResourceGetter
}
