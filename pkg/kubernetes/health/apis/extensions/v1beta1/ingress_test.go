package v1beta1

import (
	"reflect"
	"testing"

	"github.com/cakehappens/seaworthy/pkg/kubernetes/health"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestIngressHealth(t *testing.T) {
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
			gotHstatus, err := IngressHealth(tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("IngressHealth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotHstatus, tt.wantHstatus) {
				t.Errorf("IngressHealth() gotHstatus = %v, want %v", gotHstatus, tt.wantHstatus)
			}
		})
	}
}
