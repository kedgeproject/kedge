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

package spec

import (
	"encoding/json"
	"reflect"
	"testing"

	api_v1 "k8s.io/client-go/pkg/api/v1"
)

func TestUnmarshalValidateFixControllerOperations(t *testing.T) {
	tests := []struct {
		Name string
		Data []byte
		App  *DeploymentSpecMod
	}{
		{
			Name: "One container mentioned in the spec",
			Data: []byte(`
name: test
containers:
 - image: nginx
services:
  - ports:
    - port: 8080`),
			App: &DeploymentSpecMod{
				ControllerFields: ControllerFields{
					Name: "test",
					PodSpecMod: PodSpecMod{
						Containers: []Container{{Container: api_v1.Container{Name: "test", Image: "nginx"}}},
					},
					Services: []ServiceSpecMod{
						{Name: "test", Ports: []ServicePortMod{{ServicePort: api_v1.ServicePort{Port: 8080}}}},
					},
				},
			},
		},
		{
			Name: "One persistent volume mentioned in the spec",
			Data: []byte(`
name: test
containers:
 - image: nginx
services:
  - ports:
    - port: 8080
volumeClaims:
- size: 500Mi`),
			App: &DeploymentSpecMod{
				ControllerFields: ControllerFields{

					Name: "test",
					PodSpecMod: PodSpecMod{
						Containers: []Container{{Container: api_v1.Container{Name: "test", Image: "nginx"}}},
					},
					Services: []ServiceSpecMod{
						{Name: "test", Ports: []ServicePortMod{{ServicePort: api_v1.ServicePort{Port: 8080}}}},
					},
					VolumeClaims: []VolumeClaim{{Name: "test", Size: "500Mi"}},
				},
			},
		},
		{
			Name: "Multiple ports specified with any names",
			Data: []byte(`
name: test
containers:
 - image: nginx
services:
- name: nginx
  ports:
  - port: 8080
  - port: 8081
  - port: 8082`),
			App: &DeploymentSpecMod{
				ControllerFields: ControllerFields{

					Name: "test",
					PodSpecMod: PodSpecMod{
						Containers: []Container{{Container: api_v1.Container{Name: "test", Image: "nginx"}}},
					},
					Services: []ServiceSpecMod{
						{
							Name: "nginx",
							Ports: []ServicePortMod{
								{ServicePort: api_v1.ServicePort{Port: 8080, Name: "nginx-8080"}},
								{ServicePort: api_v1.ServicePort{Port: 8081, Name: "nginx-8081"}},
								{ServicePort: api_v1.ServicePort{Port: 8082, Name: "nginx-8082"}},
							},
						},
					},
				},
			},
		},
		{
			Name: "Multiple ports, some with names specified, others with no names",
			Data: []byte(`
name: test
containers:
 - image: nginx
services:
- ports:
  - port: 8080
    name: port-1
  - port: 8081
    name: port-2
  - port: 8082
  - port: 8083`),
			App: &DeploymentSpecMod{
				ControllerFields: ControllerFields{

					Name: "test",
					PodSpecMod: PodSpecMod{
						Containers: []Container{{Container: api_v1.Container{Name: "test", Image: "nginx"}}},
					},
					Services: []ServiceSpecMod{
						{
							Name: "test",
							Ports: []ServicePortMod{
								{ServicePort: api_v1.ServicePort{Port: 8080, Name: "port-1"}},
								{ServicePort: api_v1.ServicePort{Port: 8081, Name: "port-2"}},
								{ServicePort: api_v1.ServicePort{Port: 8082, Name: "test-8082"}},
								{ServicePort: api_v1.ServicePort{Port: 8083, Name: "test-8083"}},
							},
						},
					},
				},
			},
		},
		{
			Name: "Multiple ports, all with names",
			Data: []byte(`
name: test
containers:
 - image: nginx
services:
- ports:
  - port: 8080
    name: port-1
  - port: 8081
    name: port-2
  - port: 8082
    name: port-3`),
			App: &DeploymentSpecMod{
				ControllerFields: ControllerFields{

					Name: "test",
					PodSpecMod: PodSpecMod{
						Containers: []Container{{Container: api_v1.Container{Name: "test", Image: "nginx"}}},
					},
					Services: []ServiceSpecMod{
						{
							Name: "test",
							Ports: []ServicePortMod{
								{ServicePort: api_v1.ServicePort{Port: 8080, Name: "port-1"}},
								{ServicePort: api_v1.ServicePort{Port: 8081, Name: "port-2"}},
								{ServicePort: api_v1.ServicePort{Port: 8082, Name: "port-3"}},
							},
						},
					},
				},
			},
		},
		{
			Name: "Single port, without any name",
			Data: []byte(`
name: test
containers:
 - image: nginx
services:
- ports:
  - port: 8080`),
			App: &DeploymentSpecMod{
				ControllerFields: ControllerFields{

					Name: "test",
					PodSpecMod: PodSpecMod{
						Containers: []Container{{Container: api_v1.Container{Name: "test", Image: "nginx"}}},
					},
					Services: []ServiceSpecMod{
						{Name: "test", Ports: []ServicePortMod{{ServicePort: api_v1.ServicePort{Port: 8080}}}},
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {

			kController, err := GetController(test.Data)

			if err != nil {
				t.Fatalf("unable to get Kubernetes controller information from Kedge definition - %v", err)
			}

			if err := kController.Unmarshal(test.Data); err != nil {
				t.Fatalf("unable to unmarshal data - %v", err)
			}

			if err := kController.Validate(); err != nil {
				t.Fatalf("unable to validate data - %v", err)
			}

			if err := kController.Fix(); err != nil {
				t.Fatalf("unable to fix data - %v", err)
			}

			if !reflect.DeepEqual(test.App, kController) {
				t.Fatalf("==> Expected:\n%v\n==> Got:\n%v", prettyPrintObjects(test.App), prettyPrintObjects(kController))
			}
		})
	}
}

func prettyPrintObjects(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}
