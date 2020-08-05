package kubernetes

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type GetResourceOptions struct {
	// Name corresponds to `NAME` within `kubectl get TYPE NAME`
	// Cannot be combined with Filename option
	Name string

	// Type corresponds to `pod` within `kubectl get pod`
	// Cannot be combined with Filename option
	Type string

	// Namespace corresponds to the option: -n, --namespace='': If present, the namespace scope for this CLI request
	// Cannot be combined with Filename option
	Namespace string

	// Selector corresponds to the option: -l, --selector='': Selector (label query) to filter on, supports '=', '==', and '!='.(e.g. -l key1=value1,key2=value2)
	// Cannot be combined with Filename option
	Selector string

	// Filename corresponds to the option: -f, --filename=[]: Filename, directory, or URL to files identifying the resource to get from a server.
	// Cannot be combined with Name, Type, Namespace or Selector Options
	Filename string

	// Recursive corresponds to the option: -R, --recursive=false: Process the directory used in -f, --filename recursively. Useful when you want to manage
	// related manifests organized within the same directory.
	// Only valid when using the Filename Option
	Recursive bool
}

func (opts *GetResourceOptions) ToArgs() []string {
	var args []string

	if len(opts.Type) > 0 {
		args = append(args, opts.Type)
	}

	if len(opts.Name) > 0 {
		args = append(args, opts.Name)
	}

	if len(opts.Namespace) > 0 {
		args = append(args, "--namespace")
		args = append(args, opts.Namespace)
	}

	if len(opts.Selector) > 0 {
		args = append(args, "--selector")
		args = append(args, opts.Selector)
	}

	if len(opts.Filename) > 0 {
		args = append(args, "--filename")
		args = append(args, opts.Filename)
	}

	if opts.Recursive {
		args = append(args, "--recursive")
	}

	return args
}

type GetResourceOption func(opt *GetResourceOptions)

type ResourceGetter interface {
	GetResources(ctx context.Context, options ...GetResourceOption) ([]unstructured.Unstructured, error)
}

type EventGetter interface {
	GetEvents(ctx context.Context, resource unstructured.Unstructured) ([]corev1.Event, error)
}

type Client interface {
	ResourceGetter
	EventGetter
}