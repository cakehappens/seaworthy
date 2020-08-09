package v1

import (
	"reflect"
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/cakehappens/seaworthy/pkg/kubernetes/health"
)

func TestDeploymentHealth(t *testing.T) {
	type args struct {
		obj unstructured.Unstructured
	}
	tests := []struct {
		name        string
		args        args
		wantHstatus health.Status
		wantErr     bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHstatus, err := DeploymentHealth(tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeploymentHealth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotHstatus, tt.wantHstatus) {
				t.Errorf("DeploymentHealth() gotHstatus = %v, want %v", gotHstatus, tt.wantHstatus)
			}
		})
	}
}
