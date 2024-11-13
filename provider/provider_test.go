package provider

import (
	"reflect"
	"testing"

	provider "github.com/pulumi/pulumi-go-provider"
)

func TestProvider(t *testing.T) {
	tests := []struct {
		name string
		want provider.Provider
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Provider(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Provider() = %v, want %v", got, tt.want)
			}
		})
	}
}
