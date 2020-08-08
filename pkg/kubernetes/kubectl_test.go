package kubernetes

import (
	"context"
	"github.com/cakehappens/seaworthy/pkg/util/sh"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"reflect"
	"testing"
)

func TestGetEvents(t *testing.T) {
	t.Run("valid args provided", func(t *testing.T) {
		expectedUid := "abc123"
		expectedArgs := []string{
			"get",
			"events",
			"--field-selector",
			"involvedObject.uid=abc123",
			"--sort-by",
			"lastTimestamp",
			"--output",
			"json",
		}

		_, _ = GetEvents(context.Background(), expectedUid, func(option *EventerOptions) {
			option.rawResourcer = func(ctx context.Context, args ...string) ([]unstructured.Unstructured, error) {
				assert.Equal(t, expectedArgs, args)
				return nil, nil
			}
		})
	})

	type args struct {
		ctx         context.Context
		resourceUid string
		options     []EventerOption
	}
	tests := []struct {
		name        string
		args        args
		want        []corev1.Event
		wantErr     bool
		argAsserter func(t *testing.T, expected, actual []string)
	}{
		{
			name: "blah",
			args: args{
				ctx:         context.Background(),
				resourceUid: "abc123",
				options: []EventerOption{
					func(option *EventerOptions) {
						option.rawResourcer = func(ctx context.Context, args ...string) ([]unstructured.Unstructured, error) {
							return nil, nil
						}
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetEvents(tt.args.ctx, tt.args.resourceUid, tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEvents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetEvents() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetResources(t *testing.T) {
	type args struct {
		ctx     context.Context
		options []ResourcerOption
	}
	tests := []struct {
		name    string
		args    args
		want    []unstructured.Unstructured
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetResources(tt.args.ctx, tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetResources() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetResources() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKubeCtlRawResourcer(t *testing.T) {
	type args struct {
		ctx  context.Context
		args []string
	}
	tests := []struct {
		name    string
		args    args
		want    []unstructured.Unstructured
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := KubeCtlRawResourcer(tt.args.ctx, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("KubeCtlRawResourcer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KubeCtlRawResourcer() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_kubectlRawResourcer(t *testing.T) {
	type args struct {
		ctx       context.Context
		cmdRunner sh.CmdRunner
		args      []string
	}
	tests := []struct {
		name    string
		args    args
		want    []unstructured.Unstructured
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := kubectlRawResourcer(tt.args.ctx, tt.args.cmdRunner, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("kubectlRawResourcer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("kubectlRawResourcer() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_resourcesFromBytes(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []unstructured.Unstructured
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resourcesFromBytes(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("resourcesFromBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("resourcesFromBytes() got = %v, want %v", got, tt.want)
			}
		})
	}
}
