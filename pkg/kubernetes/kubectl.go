package kubernetes

import (
	"bytes"
	"context"
	"fmt"
	pkgerrors "github.com/cakehappens/seaworthy/pkg/errors"
	"github.com/cakehappens/seaworthy/pkg/util/sh"
	"io"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"os"
	"strings"
)

const resourcesFromBytesErrorFmt = "failed to convert output to resource list: %w"

func resourcesFromBytes(b []byte) ([]unstructured.Unstructured, error) {
	obj := &unstructured.Unstructured{}

	err := obj.UnmarshalJSON(b)

	if err != nil {
		if err != io.EOF {
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

func kubectlRawResourcer(ctx context.Context, cmdRunner sh.CmdRunner, args ...string) ([]unstructured.Unstructured, error) {
	var binary string
	if env := os.Getenv("SEAWORTHY_KUBECTL_PATH"); env != "" {
		binary = env
	} else {
		binary = "kubectl"
	}

	var out, outErr bytes.Buffer

	exitCode, err := cmdRunner(ctx, binary, func(opts *sh.RunOptions) {
		opts.Args = args
		opts.Stdout = &out
		opts.Stderr = &outErr
	})

	stdout := out.Bytes()
	stderr := outErr.Bytes()

	if err != nil {
		err = fmt.Errorf("unable to start/complete command: %s %s: %w", binary, strings.Join(args, " "), err)
		return nil, err
	}

	if exitCode != 0 {
		err = fmt.Errorf("command failed: %s %s: %s", binary, strings.Join(args, " "), stderr)
		return nil, err
	}

	return resourcesFromBytes(stdout)
}

func KubeCtlRawResourcer(ctx context.Context, args ...string) ([]unstructured.Unstructured, error) {
	return kubectlRawResourcer(ctx, sh.Run, args...)
}

func GetResources(ctx context.Context, options ...ResourcerOption) ([]unstructured.Unstructured, error) {
	opts := &ResourcerOptions{
		rawResourcer: KubeCtlRawResourcer,
	}

	for _, ofn := range options {
		ofn(opts)
	}

	var args []string
	args = append(args, "get")

	args = append(args, opts.GetCmdArgs()...)

	args = append(args, "--output", "json")

	return opts.rawResourcer(ctx, args...)
}

func GetEvents(ctx context.Context, resourceUid string, options ...EventerOption) ([]corev1.Event, error) {
	opts := &EventerOptions{
		rawResourcer: KubeCtlRawResourcer,
	}

	for _, ofn := range options {
		ofn(opts)
	}

	var args []string

	args = append(args, "get", "events")

	var fieldSelectors []string

	fieldSelectors = append(fieldSelectors, fmt.Sprintf("involvedObject.uid=%s", resourceUid))

	args = append(args, "--field-selector", strings.Join(fieldSelectors, ","))
	args = append(args, "--sort-by", "lastTimestamp")
	args = append(args, "--output", "json")

	objs, err := opts.rawResourcer(ctx, args...)

	if err != nil {
		return nil, err
	}

	var events []corev1.Event
	var errors pkgerrors.MultiError

	for _, obj := range objs {
		event := &corev1.Event{}

		err = runtime.DefaultUnstructuredConverter.FromUnstructured(obj.Object, &event)
		if err != nil {
			errors.Add(err)
		} else {
			events = append(events, *event)
		}
	}

	return events, errors.Return()
}
