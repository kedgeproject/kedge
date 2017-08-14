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
	"reflect"
	"testing"

	"github.com/kedgeproject/kedge/pkg/spec"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	api_v1 "k8s.io/client-go/pkg/api/v1"
)

func TestCreateServices(t *testing.T) {
	tests := []struct {
		Name    string
		App     *spec.DeploymentSpecMod
		Objects []runtime.Object
	}{
		{
			"Single container specified",
			&spec.DeploymentSpecMod{
				Name: "test",
				PodSpecMod: spec.PodSpecMod{
					Containers: []spec.Container{{Container: api_v1.Container{Image: "nginx"}}},
				},
				Services: []spec.ServiceSpecMod{
					{Name: "test", Ports: []spec.ServicePortMod{{ServicePort: api_v1.ServicePort{Port: 8080}}}},
				},
			},
			append(make([]runtime.Object, 0), &api_v1.Service{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec:       api_v1.ServiceSpec{Ports: []api_v1.ServicePort{{Port: 8080}}},
			}),
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
