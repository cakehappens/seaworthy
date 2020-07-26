package kubernetes

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/cakehappens/seaworthy/pkg/util/sh"
	"io"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"os"
	"strings"
)

type Kubectl struct {
	Binary string
}

type KubectlOption func(k *Kubectl)

func NewKubectl(options ...KubectlOption) *Kubectl {
	k := &Kubectl{}

	if env := os.Getenv("SEAWORTHY_KUBECTL_PATH"); env != "" {
		k.Binary = env
	} else {
		k.Binary = "kubectl"
	}

	for _, o := range options {
		o(k)
	}

	return k
}

func (k *Kubectl) run(ctx context.Context, args ...string) (stdout, stderr string, err error) {
	var out, outErr bytes.Buffer

	rState := sh.Run(ctx, k.Binary, func(opts *sh.RunOptions) {
		opts.Args = args
		opts.Stdout = &out
		opts.Stderr = &outErr
	})

	stdout = string(out.Bytes())
	stderr = string(outErr.Bytes())

	if !rState.Ran {
		err = rState.RunError
		return
	}

	if !rState.Success() {
		err = errors.New(stderr)
		return
	}

	return
}

type GetOptions struct {
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

	// IgnoreNotFound corresponds to the option: --ignore-not-found=false: If the requested object does not exist the command will return exit code 0.
	// Default false
	IgnoreNotFound bool
}

func (opts *GetOptions) ToArgs() []string {
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

	if opts.IgnoreNotFound {
		args = append(args, "--ignore-not-found")
	}

	return args
}

type GetOption func(opt *GetOptions)

func (k *Kubectl) Get(ctx context.Context, options ...GetOption) ([]unstructured.Unstructured, error) {
	opts := &GetOptions{}

	for _, o := range options {
		o(opts)
	}

	var args []string
	args = append(args, "get")

	args = append(args, opts.ToArgs()...)

	args = append(args, "--output")
	args = append(args, "json")

	stdout, stderr, err := k.run(ctx, args...)

	if err != nil {
		err = fmt.Errorf("%s: %w", stderr, err)
		return nil, fmt.Errorf("%s %s: %w", k.Binary, strings.Join(args, " "), err)
	}

	// fmt.Println("result from kubectl")
	// fmt.Printf("%+v\n", stdout)

	obj := &unstructured.Unstructured{}

	err = obj.UnmarshalJSON([]byte(stdout))

	if err != nil {
		if  err != io.EOF {
			return nil, fmt.Errorf("json unmarshal: %w", err)
		}
	}

	if obj.IsList() {
		objLs, err := obj.ToList()
		if err != nil {
			return nil, err
		}

		return objLs.Items, nil
	}

	return []unstructured.Unstructured{*obj}, nil
}
