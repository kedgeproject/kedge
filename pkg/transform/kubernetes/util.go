package kubernetes

import (
	"github.com/kedgeproject/kedge/pkg/spec"

	api_v1 "k8s.io/client-go/pkg/api/v1"
)

// This function will search in the pod level volumes
// and see if the volume with given name is defined
func isVolumeDefined(volumes []api_v1.Volume, name string) bool {
	for _, v := range volumes {
		if v.Name == name {
			return true
		}
	}
	return false
}

// search through all the persistent volumes defined in the root level
func isPVCDefined(volumes []spec.VolumeClaim, name string) bool {
	for _, v := range volumes {
		if v.Name == name {
			return true
		}
	}
	return false
}
