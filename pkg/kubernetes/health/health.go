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
	"errors"
	"fmt"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// group/version/kind
var checkFuncs = make(map[string]map[string]map[string]CheckFunc)

func MustRegisterCheckFunc(group, version, kind string, fn CheckFunc) {
	err := RegisterCheckFunc(group, version, kind, fn)
	if err != nil {
		panic(err)
	}
}

func RegisterCheckFunc(group, version, kind string, fn CheckFunc) error {
	if vMap, ok := checkFuncs[group]; ok {
		if kMap, ok := vMap[version]; ok {
			if registeredFn, ok := kMap[kind]; ok && registeredFn != nil {
				return errors.New(fmt.Sprintf("already registered: %s/%s, kind: %s", group, version, kind))
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

func GetCheckFunc(group, version, kind string) (CheckFunc, error) {
	if vMap, ok := checkFuncs[group]; ok {
		if kMap, ok := vMap[version]; ok {
			if registeredFn, ok := kMap[kind]; ok && registeredFn != nil {
				return registeredFn, nil
			}
		}
	}

	return nil, errors.New(fmt.Sprintf("unregistered resource: %s/%s, kind: %s", group, version, kind))
}

type CheckFunc func(obj unstructured.Unstructured) (Status, error)

// Represents resource health status
type StatusCode string

const (
	// Indicates that health assessment failed and actual health status is unknown
	Unknown StatusCode = "Unknown"
	// Progressing health status means that resource is not healthy but still have a chance to reach healthy state
	Progressing StatusCode = "Progressing"
	// Resource is 100% healthy
	Healthy StatusCode = "Healthy"
	// Assigned to resources that are suspended or paused. The typical example is a
	// [suspended](https://kubernetes.io/docs/tasks/job/automated-tasks-with-cron-jobs/#suspend) CronJob.
	Suspended StatusCode = "Suspended"
	// Degrade status is used if resource status indicates failure or resource could not reach healthy state
	// within some timeout.
	Degraded StatusCode = "Degraded"
	// Indicates that the resource does have a health check available to run
	Unsupported StatusCode = "Unsupported"
	// Indicates that resource is missing in the cluster.
	Missing StatusCode = "Missing"
)

// Holds health assessment results
type Status struct {
	Code    StatusCode `json:"status,omitempty"`
	Message string     `json:"message,omitempty"`
}

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

type ResourceHealthOptions struct {
	Override CheckFunc
}

type ResourceHealthOption func(options *ResourceHealthOptions)

// GetResourceHealth returns the health of a k8s resource
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
