package kubernetes

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// resourcer describes how to get the resources args and then returns the resources as unstructured objects
type rawResourcer func(ctx context.Context, args ...string) ([]unstructured.Unstructured, error)

// ResourcerOptions is part of the functional API for Resourcer
type ResourcerOptions struct {
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

	rawResourcer
}

// GetCmdArgs is part of the functional API for Resourcer
func (opts *ResourcerOptions) GetCmdArgs() []string {
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

// ResourcerOption is part of the functional API for Resourcer
type ResourcerOption func(opt *ResourcerOptions)

// Resourcer should return kubernetes resources
type Resourcer func(ctx context.Context, options ...ResourcerOption) ([]unstructured.Unstructured, error)

// EventerOptions is part of the functional API for Eventer
type EventerOptions struct {
	rawResourcer
}

// EventerOption is part of the functional API for Eventer
type EventerOption func(option *EventerOptions)

// Eventer should return kubernetes events
type Eventer func(ctx context.Context, resource unstructured.Unstructured, options ...EventerOption) ([]corev1.Event, error)
