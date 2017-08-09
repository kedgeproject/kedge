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
