package fixtures

import (
	"github.com/surajssd/kapp/pkg/spec"
	api_v1 "k8s.io/client-go/pkg/api/v1"
)

var SinglePersistentVolumeApp spec.App = spec.App{
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
			Name: "test",
			Ports: []spec.ServicePortMod{
				{
					ServicePort: api_v1.ServicePort{
						Port: 8080,
					},
				},
			},
		},
	},
	PersistentVolumes: []spec.PersistentVolume{
		{
			Name: "test",
			Size: "500Mi",
		},
	},
}
