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

package encoding

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/kedgeproject/kedge/pkg/spec"

	api_v1 "k8s.io/client-go/pkg/api/v1"
)

func TestDecode(t *testing.T) {
	tests := []struct {
		Name string
		Data []byte
		App  *spec.App
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
			App: &spec.App{
				Name: "test",
				PodSpecMod: spec.PodSpecMod{
					Containers: []spec.Container{{Container: api_v1.Container{Name: "test", Image: "nginx"}}},
				},
				Services: []spec.ServiceSpecMod{
					{Name: "test", Ports: []spec.ServicePortMod{{ServicePort: api_v1.ServicePort{Port: 8080}}}},
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
			App: &spec.App{
				Name: "test",
				PodSpecMod: spec.PodSpecMod{
					Containers: []spec.Container{{Container: api_v1.Container{Name: "test", Image: "nginx"}}},
				},
				Services: []spec.ServiceSpecMod{
					{Name: "test", Ports: []spec.ServicePortMod{{ServicePort: api_v1.ServicePort{Port: 8080}}}},
				},
				VolumeClaims: []spec.VolumeClaim{{Name: "test", Size: "500Mi"}},
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
			App: &spec.App{
				Name: "test",
				PodSpecMod: spec.PodSpecMod{
					Containers: []spec.Container{{Container: api_v1.Container{Name: "test", Image: "nginx"}}},
				},
				Services: []spec.ServiceSpecMod{
					{
						Name: "nginx",
						Ports: []spec.ServicePortMod{
							{ServicePort: api_v1.ServicePort{Port: 8080, Name: "nginx-8080"}},
							{ServicePort: api_v1.ServicePort{Port: 8081, Name: "nginx-8081"}},
							{ServicePort: api_v1.ServicePort{Port: 8082, Name: "nginx-8082"}},
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
			App: &spec.App{
				Name: "test",
				PodSpecMod: spec.PodSpecMod{
					Containers: []spec.Container{{Container: api_v1.Container{Name: "test", Image: "nginx"}}},
				},
				Services: []spec.ServiceSpecMod{
					{
						Name: "test",
						Ports: []spec.ServicePortMod{
							{ServicePort: api_v1.ServicePort{Port: 8080, Name: "port-1"}},
							{ServicePort: api_v1.ServicePort{Port: 8081, Name: "port-2"}},
							{ServicePort: api_v1.ServicePort{Port: 8082, Name: "test-8082"}},
							{ServicePort: api_v1.ServicePort{Port: 8083, Name: "test-8083"}},
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
			App: &spec.App{
				Name: "test",
				PodSpecMod: spec.PodSpecMod{
					Containers: []spec.Container{{Container: api_v1.Container{Name: "test", Image: "nginx"}}},
				},
				Services: []spec.ServiceSpecMod{
					{
						Name: "test",
						Ports: []spec.ServicePortMod{
							{ServicePort: api_v1.ServicePort{Port: 8080, Name: "port-1"}},
							{ServicePort: api_v1.ServicePort{Port: 8081, Name: "port-2"}},
							{ServicePort: api_v1.ServicePort{Port: 8082, Name: "port-3"}},
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
			App: &spec.App{
				Name: "test",
				PodSpecMod: spec.PodSpecMod{
					Containers: []spec.Container{{Container: api_v1.Container{Name: "test", Image: "nginx"}}},
				},
				Services: []spec.ServiceSpecMod{
					{Name: "test", Ports: []spec.ServicePortMod{{ServicePort: api_v1.ServicePort{Port: 8080}}}},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			app, err := Decode(test.Data)
			if err != nil {
				t.Fatalf("Unable to run Decode(), and error occurred: %v", err)
			}

			if !reflect.DeepEqual(test.App, app) {
				t.Fatalf("==> Expected:\n%v\n==> Got:\n%v", prettyPrintObjects(test.App), prettyPrintObjects(app))
			}
		})
	}
}

func prettyPrintObjects(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}
