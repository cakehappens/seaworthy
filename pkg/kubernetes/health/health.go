/*
Copyright 2020 argoproj/gitops-engine Authors.

Modifications Copyright 2020 cakehappens/seaworthy Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package health

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// group/version/kind
var checkFuncs = make(map[string]map[string]map[string]CheckFunc)

// MustRegisterCheckFunc registers new resource check funcs or panics if already registered
func MustRegisterCheckFunc(group, version, kind string, fn CheckFunc) {
	err := RegisterCheckFunc(group, version, kind, fn)
	if err != nil {
		panic(err)
	}
}

// RegisterCheckFunc registers new resource check funcs or returns an error if already registered
func RegisterCheckFunc(group, version, kind string, fn CheckFunc) error {
	if vMap, ok := checkFuncs[group]; ok {
		if kMap, ok := vMap[version]; ok {
			if registeredFn, ok := kMap[kind]; ok && registeredFn != nil {
				return fmt.Errorf("already registered: %s/%s, kind: %s", group, version, kind)
			}
		} else {
			checkFuncs[group][version] = make(map[string]CheckFunc)
		}
	} else {
		checkFuncs[group] = make(map[string]map[string]CheckFunc)
		checkFuncs[group][version] = make(map[string]CheckFunc)
	}

	checkFuncs[group][version][kind] = fn
	return nil
}

// GetCheckFunc returns the resource health check function registered for the given group, version and kind
// or an error if no check function exists
func GetCheckFunc(group, version, kind string) (CheckFunc, error) {
	if vMap, ok := checkFuncs[group]; ok {
		if kMap, ok := vMap[version]; ok {
			if registeredFn, ok := kMap[kind]; ok && registeredFn != nil {
				return registeredFn, nil
			}
		}
	}

	return nil, fmt.Errorf("unregistered resource: %s/%s, kind: %s", group, version, kind)
}

// CheckFunc describes the function signature for all check functions
type CheckFunc func(obj unstructured.Unstructured) (Status, error)

// StatusCode represents resource health status
type StatusCode string

const (
	// Unknown indicates that health assessment failed and actual health status is unknown
	Unknown StatusCode = "Unknown"
	// Progressing health status means that resource is not healthy but still have a chance to reach healthy state
	Progressing StatusCode = "Progressing"
	// Healthy indicates the resource is 100% healthy
	Healthy StatusCode = "Healthy"
	// Suspended indicates the resource is suspended or paused. The typical example is a
	// [suspended](https://kubernetes.io/docs/tasks/job/automated-tasks-with-cron-jobs/#suspend) CronJob.
	Suspended StatusCode = "Suspended"
	// Degraded status is used if resource status indicates failure or resource could not reach healthy state
	// within some timeout.
	Degraded StatusCode = "Degraded"
	// Unsupported indicates that the resource does have a health check available to run
	Unsupported StatusCode = "Unsupported"
	// Missing indicates that resource is missing from the cluster.
	Missing StatusCode = "Missing"
)

// Status holds health assessment results
type Status struct {
	Code    StatusCode `json:"status,omitempty"`
	Message string     `json:"message,omitempty"`
}

// NewHealthyHealthStatus returns a healthy Status struct
func NewHealthyHealthStatus() Status {
	return Status{
		Code:    Healthy,
		Message: "All systems nominal",
	}
}

// codeOrder is a list of health codes in order of most healthy to least healthy
var codeOrder = []StatusCode{
	Healthy,
	Suspended,
	Progressing,
	Degraded,
	Missing,
	Unsupported,
	Unknown,
}

// GetCodeOrder returns a list of health codes in order of most healthy to least healthy
func GetCodeOrder() []StatusCode {
	return codeOrder
}

// IsWorse returns whether or not the new health status code is a worser condition than the current
func IsWorse(current, new StatusCode) bool {
	currentIndex := 0
	newIndex := 0
	for i, code := range codeOrder {
		if current == code {
			currentIndex = i
		}
		if new == code {
			newIndex = i
		}
	}
	return newIndex > currentIndex
}

// ResourceHealthOptions is part of the functional API for ResourceHealth
type ResourceHealthOptions struct {
	Override CheckFunc
}

// ResourceHealthOption is part of the functional API for ResourceHealth
type ResourceHealthOption func(options *ResourceHealthOptions)

// ResourceHealth returns the health of a k8s resource
func ResourceHealth(obj unstructured.Unstructured, options ...ResourceHealthOption) Status {
	if obj.GetDeletionTimestamp() != nil {
		return Status{
			Code:    Progressing,
			Message: "Pending deletion",
		}
	}

	opts := &ResourceHealthOptions{}

	for _, o := range options {
		o(opts)
	}

	var checkFunc CheckFunc
	var err error

	if opts.Override != nil {
		checkFunc = opts.Override
	} else {
		gvk := obj.GroupVersionKind()

		checkFunc, err = GetCheckFunc(gvk.Group, gvk.Version, gvk.Kind)
		if err != nil {
			return Status{
				Code:    Unsupported,
				Message: err.Error(),
			}
		}
	}

	health, err := checkFunc(obj)
	if err != nil {
		return Status{
			Code:    Unknown,
			Message: err.Error(),
		}
	}

	return health
}
