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
	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"

	log "github.com/Sirupsen/logrus"
)

// CoreOperations takes in the raw input data and extracts the controller
// information, and proceeds to run the controller specific operations on the
// parsed data.

// The "core operations" is important to Kedge as it unmarshals, validates the artifacts
// as well as return the correct transformation of said artifact.
//
// Returns the converted Kubernetes objects, extra resources and an error, if any.
func CoreOperations(data []byte) ([]runtime.Object, []string, error) {

	var app App
	err := app.LoadData(data)
	if err != nil {
		return nil, nil, errors.Wrap(err, "App could not be loaded into internal struct")
	}

	// Validate said data
	if err := app.Validate(); err != nil {
		return nil, nil, errors.Wrap(err, "unable to validate data")
	}

	// Fix any problems that may occur
	if err := app.Fix(); err != nil {
		return nil, nil, errors.Wrap(err, "unable to fix data")
	}

	ros, err := app.CreateK8sObjects()
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to transform data")
	}

	scheme, err := GetScheme()
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to get scheme")
	}

	for _, runtimeObject := range ros {
		if err := SetGVK(runtimeObject, scheme); err != nil {
			return nil, nil, errors.Wrap(err, "unable to set Group, Version and Kind for generated Kubernetes resources")
		}
	}

	return ros, app.IncludeResources, nil
}

// LoadData - unmarshal data into App struct
func (app *App) LoadData(data []byte) error {
	err := yaml.Unmarshal(data, app)
	if err != nil {
		return errors.Wrap(err, "App could not be unmarshaled into internal struct")
	}
	log.Debugf("object unmarshalled: %v\n", PrettyPrintObjects(app))
	return nil
}
