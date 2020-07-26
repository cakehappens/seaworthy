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

package v1

import (
	"fmt"
	"github.com/cakehappens/seaworthy/pkg/kubernetes/health"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/kubectl/pkg/scheme"
)

const (
	DeploymentKind = "Deployment"
)

func DeploymentHealth(obj unstructured.Unstructured) (hstatus health.Status, err error) {
	deployment := &appsv1.Deployment{}
	err = scheme.Scheme.Convert(&obj, deployment, nil)
	if err != nil {
		err = fmt.Errorf("failed to convert %T to %T: %w", obj, deployment, err)
		return
	}
	if deployment.Spec.Paused {
		hstatus = health.Status{
			Code:    health.Suspended,
			Message: "Deployment is paused",
		}

		return
	}

	if deployment.Generation <= deployment.Status.ObservedGeneration {
		cond := getDeploymentCondition(deployment.Status, v1.DeploymentProgressing)
		if cond != nil && cond.Reason == "ProgressDeadlineExceeded" {
			hstatus = health.Status{
				Code:    health.Degraded,
				Message: fmt.Sprintf("Deployment %q exceeded its progress deadline", obj.GetName()),
			}

			return
		} else if deployment.Spec.Replicas != nil && deployment.Status.UpdatedReplicas < *deployment.Spec.Replicas {
			hstatus = health.Status{
				Code:    health.Progressing,
				Message: fmt.Sprintf("Waiting for rollout to finish: %d out of %d new replicas have been updated...", deployment.Status.UpdatedReplicas, *deployment.Spec.Replicas),
			}

			return
		} else if deployment.Status.Replicas > deployment.Status.UpdatedReplicas {
			hstatus = health.Status{
				Code:    health.Progressing,
				Message: fmt.Sprintf("Waiting for rollout to finish: %d old replicas are pending termination...", deployment.Status.Replicas-deployment.Status.UpdatedReplicas),
			}

			return
		} else if deployment.Status.AvailableReplicas < deployment.Status.UpdatedReplicas {
			hstatus = health.Status{
				Code:    health.Progressing,
				Message: fmt.Sprintf("Waiting for rollout to finish: %d of %d updated replicas are available...", deployment.Status.AvailableReplicas, deployment.Status.UpdatedReplicas),
			}

			return
		}
	} else {
		hstatus = health.Status{
			Code:    health.Progressing,
			Message: "Waiting for rollout to finish: observed deployment generation less then desired generation",
		}

		return
	}

	hstatus = health.NewHealthyHealthStatus()

	return
}

func getDeploymentCondition(status v1.DeploymentStatus, condType v1.DeploymentConditionType) *v1.DeploymentCondition {
	for i := range status.Conditions {
		c := status.Conditions[i]
		if c.Type == condType {
			return &c
		}
	}
	return nil
}
