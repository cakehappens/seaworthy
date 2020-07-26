package v1beta1

import (
	extv1beta "k8s.io/api/extensions/v1beta1"
)

const (
	GroupName = extv1beta.GroupName
	// Unfortunately this isn't a constant in the upstream repo
	Version = "v1beta1"
)
