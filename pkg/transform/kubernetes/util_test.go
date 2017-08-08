package kubernetes

import (
	"testing"

	"github.com/kedgeproject/kedge/pkg/spec"

	api_v1 "k8s.io/client-go/pkg/api/v1"
)

func TestIsVolumeDefined(t *testing.T) {
	volumes := []api_v1.Volume{{Name: "foo"}, {Name: "bar"}, {Name: "baz"}}

	tests := []struct {
		Search string
		Output bool
	}{
		{Search: "bar", Output: true},
		{Search: "fooo", Output: false},
		{Search: "baz", Output: true},
	}

	t.Logf("volumes: %+v", volumes)
	for _, test := range tests {
		if test.Output != isVolumeDefined(volumes, test.Search) {
			t.Errorf("expected output to match, but did not match for"+
				" volumes %+v and search query %q", volumes, test.Search)
		} else {
			t.Logf("test passed for search query %q", test.Search)
		}
	}

}

func TestIsPVCDefined(t *testing.T) {
	volumes := []spec.VolumeClaim{{Name: "foo"}, {Name: "bar"}, {Name: "baz"}}

	tests := []struct {
		Search string
		Output bool
	}{
		{Search: "bar", Output: true},
		{Search: "fooo", Output: false},
		{Search: "baz", Output: true},
	}

	t.Logf("volumes: %+v", volumes)
	for _, test := range tests {
		if test.Output != isPVCDefined(volumes, test.Search) {
			t.Errorf("expected output to match, but did not match for"+
				" volumes %+v and search query %q", volumes, test.Search)
		} else {
			t.Logf("test passed for search query %q", test.Search)
		}
	}
}
