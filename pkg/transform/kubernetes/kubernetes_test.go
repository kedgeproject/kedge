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

	"reflect"

	encodingFixtures "github.com/kedgeproject/kedge/pkg/encoding/fixtures"
	"github.com/kedgeproject/kedge/pkg/spec"
	transformFixtures "github.com/kedgeproject/kedge/pkg/transform/fixtures"
	"k8s.io/client-go/pkg/runtime"
)

func TestCreateServices(t *testing.T) {
	tests := []struct {
		Name    string
		App     *spec.App
		Objects []runtime.Object
	}{
		{
			"Single container specified",
			&encodingFixtures.SingleContainerApp,
			append(make([]runtime.Object, 0), transformFixtures.SingleContainerService),
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			object, err := createServices(test.App)
			if err != nil {
				t.Fatalf("Creating services failed: %v", err)
			}
			if !reflect.DeepEqual(test.Objects, object) {
				t.Fatalf("Expected:\n%v\nGot:\n%v", test.Objects, object)
			}
		})
	}
}

// TODO: add test for auto naming of single persistent volume
