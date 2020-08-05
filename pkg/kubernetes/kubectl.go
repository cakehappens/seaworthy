package kubernetes

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	pkgerrors "github.com/cakehappens/seaworthy/pkg/errors"
	"github.com/cakehappens/seaworthy/pkg/util/sh"
	"io"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/kubectl/pkg/scheme"
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

func (k *Kubectl) run(ctx context.Context, args ...string) (stdout, stderr []byte, err error) {
	var out, outErr bytes.Buffer

	rState := sh.Run(ctx, k.Binary, func(opts *sh.RunOptions) {
		opts.Args = args
		opts.Stdout = &out
		opts.Stderr = &outErr
	})

	stdout = out.Bytes()
	stderr = outErr.Bytes()

	if !rState.Ran {
		err = rState.RunError
		return
	}

	if !rState.Success() {
		err = errors.New(string(stderr))
		err = fmt.Errorf("%s %s: %w", k.Binary, strings.Join(args, " "), err)
		return
	}

	return
}

func (k *Kubectl) GetResources(ctx context.Context, options ...GetResourceOption) ([]unstructured.Unstructured, error) {
	opts := &GetResourceOptions{}

	for _, o := range options {
		o(opts)
	}

	var args []string
	args = append(args, "get")

	args = append(args, opts.ToArgs()...)

	args = append(args, "--output", "json")

	stdout, _, err := k.run(ctx, args...)

	if err != nil {
		return nil, err
	}

	return resourcesFromBytes(stdout)
}

func (k *Kubectl) GetEvents(ctx context.Context, resourceUid string) ([]corev1.Event, error) {
	var args []string

	args = append(args, "get", "events")

	var fieldSelectors []string

	fieldSelectors = append(fieldSelectors, fmt.Sprintf("involvedObject.uid=%s", resourceUid))

	args = append(args, "--field-selector", strings.Join(fieldSelectors, ","))
	args = append(args, "--sort-by", "lastTimestamp")
	args = append(args, "--output", "json")

	stdout, _, err := k.run(ctx, args...)

	if err != nil {
		return nil, err
	}

	objs, err := resourcesFromBytes(stdout)
	if err != nil {
		return nil, err
	}

	var events []corev1.Event
	var errors pkgerrors.MultiError

	for _, obj := range objs {
		event := &corev1.Event{}

		err = scheme.Scheme.Convert(&obj, event, nil)
		if err != nil {
			errors.Add(err)
			events = append(events, *event)
		}
	}

	return events, errors.Return()
}

const resourcesFromBytesErrorFmt = "failed to convert output to resource list: %w"

func resourcesFromBytes(b []byte) ([]unstructured.Unstructured, error) {
	obj := &unstructured.Unstructured{}

	err := obj.UnmarshalJSON(b)

	if err != nil {
		if  err != io.EOF {
			err = fmt.Errorf("json unmarshal: %w", err)
			return nil, fmt.Errorf(resourcesFromBytesErrorFmt, err)
		}
	}

	if obj.IsList() {
		objLs, err := obj.ToList()
		if err != nil {
			err = fmt.Errorf("isList() = true, but toList() failed: %w", err)
			return nil, fmt.Errorf(resourcesFromBytesErrorFmt, err)
		}

		return objLs.Items, nil
	}

	return []unstructured.Unstructured{*obj}, nil
}
