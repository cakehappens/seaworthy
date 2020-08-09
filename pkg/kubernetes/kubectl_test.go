package kubernetes

import (
	"context"
	"github.com/cakehappens/seaworthy/pkg/util/sh"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"reflect"
	"sigs.k8s.io/yaml"
	"testing"
)

func loadTestResource(yamlPath string) *unstructured.Unstructured {
	yamlBytes, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		panic(err)
	}
	var obj unstructured.Unstructured
	err = yaml.Unmarshal(yamlBytes, &obj)
	if err != nil {
		panic(err)
	}
	return &obj
}

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
	type want struct {
		apiVersion string
		kind       string
		name       string
		uid        string
	}
	tests := []struct {
		name    string
		args    args
		wants   []want
		wantErr bool
	}{
		{
			name: "given many events__expect list returned",
			args: args{
				ctx:         context.Background(),
				resourceUid: "abc123",
				options: []EventerOption{
					func(option *EventerOptions) {
						option.rawResourcer = func(ctx context.Context, args ...string) ([]unstructured.Unstructured, error) {
							obj := loadTestResource("./test_data/events_many.yml")
							objLs, err := obj.ToList()
							if err != nil {
								panic(err)
							}

							return objLs.Items, nil
						}
					},
				},
			},
			wants: []want{
				{
					apiVersion: "v1",
					kind:       "Event",
					name:       "ssm-agent.16293412df341d17",
					uid:        "aa37e1c0-d938-11ea-9375-02ecadc8ef22",
				},
				{
					apiVersion: "v1",
					kind:       "Event",
					name:       "ssm-agent.16293412e6b5bf59",
					uid:        "aa4c1f84-d938-11ea-9375-02ecadc8ef33",
				},
				{
					apiVersion: "v1",
					kind:       "Event",
					name:       "ssm-agent.16293412e6b5bf59",
					uid:        "aa4c1f84-d938-11ea-9375-02ecadc8ef44",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetEvents(tt.args.ctx, tt.args.resourceUid, tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEvents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			require.Nil(t, err)
			require.Equal(t, len(tt.wants), len(got))

			t.Logf("%+v", got)

			for i, w := range tt.wants {
				assert.Equal(t, w.name, got[i].Name)
				assert.Equal(t, w.apiVersion, got[i].APIVersion)
				assert.Equal(t, w.kind, got[i].Kind)
				assert.Equal(t, w.uid, string(got[i].UID))
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
