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

package v1beta1

import (
	"fmt"
	"github.com/cakehappens/seaworthy/pkg/kubernetes/health"
	extv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/kubernetes/scheme"
)

const (
	IngressKind = "Ingress"
)

func IngressHealth(obj unstructured.Unstructured) (hstatus health.Status, err error) {
	ingress := &extv1beta1.Ingress{}
	err = scheme.Scheme.Convert(&obj, ingress, nil)
	if err != nil {
		err = fmt.Errorf("failed to convert %T to %T: %v", obj, ingress, err)
		return
	}

	if len(ingress.Status.LoadBalancer.Ingress) <= 0 {
		hstatus = health.Status{
			Code:    health.Progressing,
			Message: "Working on it...",
		}
		return
	}

	hstatus = health.NewHealthyHealthStatus()
	return
}
