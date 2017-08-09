/*
Copyright 2017 The Kedge Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
