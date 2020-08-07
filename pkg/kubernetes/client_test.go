package kubernetes

import (
	"reflect"
	"testing"
)

func TestResourcerOptions_GetCmdArgs(t *testing.T) {
	type fields struct {
		Name         string
		Type         string
		Namespace    string
		Selector     string
		Filename     string
		Recursive    bool
		rawResourcer rawResourcer
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
			opts := &ResourcerOptions{
				Name:         tt.fields.Name,
				Type:         tt.fields.Type,
				Namespace:    tt.fields.Namespace,
				Selector:     tt.fields.Selector,
				Filename:     tt.fields.Filename,
				Recursive:    tt.fields.Recursive,
				rawResourcer: tt.fields.rawResourcer,
			}
			if got := opts.GetCmdArgs(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCmdArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}
