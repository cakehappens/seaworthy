package v1beta1

import (
	extv1beta "k8s.io/api/extensions/v1beta1"
)

const (
	// GroupName is the K8s API Group
	GroupName = extv1beta.GroupName
	// Version is the K8s API Version
	// Unfortunately this isn't a constant in the upstream repo
	Version = "v1beta1"
)
