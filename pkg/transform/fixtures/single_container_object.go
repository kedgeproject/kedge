package fixtures

import api_v1 "k8s.io/client-go/pkg/api/v1"

var SingleContainerService *api_v1.Service = &api_v1.Service{
	ObjectMeta: api_v1.ObjectMeta{
		Name: "test",
	},
	Spec: api_v1.ServiceSpec{
		Ports: []api_v1.ServicePort{
			{
				Port: 8080,
			},
		},
	},
}
