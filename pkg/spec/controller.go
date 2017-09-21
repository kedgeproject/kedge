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
	"fmt"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

// Every controller that Kedge supports is required to implement this interface
type ControllerInterface interface {
	// Unmarshals input YAML data to the corresponding Kedge controller spec
	Unmarshal(data []byte) error

	// Validates the unmarshalled data
	Validate() error

	// Fixes the unmarshalled data, e.g. auto population/generation of fields
	Fix() error

	// Transforms the data in Kedge spec to Kubernetes' resource objects
	Transform() ([]runtime.Object, []string, error)
}

// GetController takes in raw input data, and returns the intended controller
// defined in the Kedge definition.
// Returns an error if the controller is not supported by Kedge
func GetController(data []byte) (ControllerInterface, error) {
	var specController Controller
	yaml.Unmarshal(data, &specController)

	switch specController.Controller {
	// If no controller is defined, we default to deployment controller
	case "", "deployment":
		return &DeploymentSpecMod{}, nil
	case "job":
		return &JobSpecMod{}, nil
	case "deploymentconfig":
		return &DeploymentConfigSpecMod{}, nil
	default:
		return nil, fmt.Errorf("invalid controller: %v", specController.Controller)
	}
}

// CoreOperations takes in the raw input data and extracts the controller
// information, and proceeds to run the controller specific operations on the
// parsed data.

// The "core operations" is important to Kedge as it unmarshals, validates the artifacts
// as well as return the correct transformation of said artifact.
//
// Returns the converted Kubernetes objects, extra resources and an error, if any.
func CoreOperations(data []byte) ([]runtime.Object, []string, error) {

	// Retrieve the selected controller
	kController, err := GetController(data)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to get Kubernetes controller information from Kedge definition")
	}

	// Unmarshal all data
	if err := kController.Unmarshal(data); err != nil {
		return nil, nil, errors.Wrap(err, "unable to unmarshal data")
	}

	// Validate said data
	if err := kController.Validate(); err != nil {
		return nil, nil, errors.Wrap(err, "unable to validate data")
	}

	// Fix any problems that may occur
	if err := kController.Fix(); err != nil {
		return nil, nil, errors.Wrap(err, "unable to fix data")
	}

	// Transform! Here we shall transform all the data to their Kubernetes (or OpenShift) object equivilant.
	artifacts, extraResources, err := kController.Transform()
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to transform data")
	}

	// In the end we return both the transformed data as well as any extra resources passed in.
	return artifacts, extraResources, nil
}
