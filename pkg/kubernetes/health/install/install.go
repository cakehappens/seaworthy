package install

import (
	"github.com/cakehappens/seaworthy/pkg/kubernetes/health"
	appsv1 "github.com/cakehappens/seaworthy/pkg/kubernetes/health/apis/apps/v1"
	extv1beta1 "github.com/cakehappens/seaworthy/pkg/kubernetes/health/apis/extensions/v1beta1"
)

func init() {
	health.MustRegisterCheckFunc(appsv1.GroupName, appsv1.Version, appsv1.DeploymentKind, appsv1.DeploymentHealth)
	health.MustRegisterCheckFunc(extv1beta1.GroupName, extv1beta1.Version, extv1beta1.IngressKind, extv1beta1.IngressHealth)
	health.MustRegisterCheckFunc(extv1beta1.GroupName, extv1beta1.Version, appsv1.DeploymentKind, appsv1.DeploymentHealth)
}
