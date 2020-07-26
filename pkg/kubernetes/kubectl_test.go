package kubernetes

import (
	"context"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"reflect"
	"testing"
)

func TestGetOptions_ToArgs(t *testing.T) {
	type fields struct {
		Name           string
		Type           string
		Namespace      string
		Selector       string
		Filename       string
		Recursive      bool
		IgnoreNotFound bool
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &GetOptions{
				Name:           tt.fields.Name,
				Type:           tt.fields.Type,
				Namespace:      tt.fields.Namespace,
				Selector:       tt.fields.Selector,
				Filename:       tt.fields.Filename,
				Recursive:      tt.fields.Recursive,
				IgnoreNotFound: tt.fields.IgnoreNotFound,
			}
			if got := opts.ToArgs(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKubectl_Get(t *testing.T) {
	type fields struct {
		Binary string
	}
	type args struct {
		ctx     context.Context
		options []GetOption
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []unstructured.Unstructured
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &Kubectl{
				Binary: tt.fields.Binary,
			}
			got, err := k.Get(tt.args.ctx, tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKubectl_run(t *testing.T) {
	type fields struct {
		Binary string
	}
	type args struct {
		ctx  context.Context
		args []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		want1   string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &Kubectl{
				Binary: tt.fields.Binary,
			}
			got, got1, err := k.run(tt.args.ctx, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("run() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("run() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestNewKubectl(t *testing.T) {
	type args struct {
		options []KubectlOption
	}
	tests := []struct {
		name string
		args args
		want *Kubectl
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewKubectl(tt.args.options...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewKubectl() = %v, want %v", got, tt.want)
			}
		})
	}
}
