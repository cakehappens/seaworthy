package v1

import (
	appsv1 "k8s.io/api/apps/v1"
)

const (
	GroupName = appsv1.GroupName
	// Unfortunately this isn't a constant in the upstream repo
	Version = "v1"
)
