package v1

import (
	appsv1 "k8s.io/api/apps/v1"
)

const (
	// GroupName is the K8s API Group
	GroupName = appsv1.GroupName
	// Version is the K8s API Version
	// Unfortunately this isn't a constant in the upstream repo
	Version = "v1"
)
