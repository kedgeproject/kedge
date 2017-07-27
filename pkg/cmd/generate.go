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

package cmd

import (
	"fmt"
	"os"

	"github.com/kedgeproject/kedge/pkg/encoding"
	"github.com/kedgeproject/kedge/pkg/transform/kubernetes"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

func Generate(paths []string) error {

	files, err := GetAllYAMLFiles(paths)
	if err != nil {
		return errors.Wrap(err, "unable to get YAML files")
	}

	inputs, err := getApplicationsFromFiles(files)
	if err != nil {
		return errors.Wrap(err, "unable to get kedge definitions from input files")
	}

	for _, input := range inputs {
		app, err := encoding.Decode(input.data)
		if err != nil {
			return errors.Wrap(err, "unable to unmarshal data")
		}

		ros, err := kubernetes.Transform(app)
		if err != nil {
			return errors.Wrap(err, "unable to transform data")
		}

		for _, runtimeObject := range ros {

			data, err := yaml.Marshal(runtimeObject)
			if err != nil {
				return errors.Wrap(err, "failed to marshal object")
			}

			writeObject := func(o runtime.Object, data []byte) error {
				_, err := fmt.Fprintln(os.Stdout, "---")
				if err != nil {
					return errors.Wrap(err, "could not print to STDOUT")
				}

				_, err = os.Stdout.Write(data)
				return errors.Wrap(err, "could not write to STDOUT")
			}

			err = writeObject(runtimeObject, data)
			if err != nil {
				return errors.Wrap(err, "failed to write object")
			}
		}
	}
	return nil
}
