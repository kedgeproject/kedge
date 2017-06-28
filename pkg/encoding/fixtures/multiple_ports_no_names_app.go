package fixtures

import (
	"github.com/kedgeproject/kedge/pkg/spec"
	api_v1 "k8s.io/client-go/pkg/api/v1"
)

var MultiplePortsNoNamesApp spec.App = spec.App{
	Name: "test",
	Containers: []spec.Container{
		{
			Container: api_v1.Container{
				Image: "nginx",
			},
		},
	},
	Services: []spec.ServiceSpecMod{
		{
			Name: "nginx",
			Ports: []spec.ServicePortMod{
				{
					ServicePort: api_v1.ServicePort{
						Port: 8080,
						Name: "nginx-8080",
					},
				},
				{
					ServicePort: api_v1.ServicePort{
						Port: 8081,
						Name: "nginx-8081",
					},
				},
				{
					ServicePort: api_v1.ServicePort{
						Port: 8082,
						Name: "nginx-8082",
					},
				},
			},
		},
	},
}
