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
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/kedgeproject/kedge/pkg/spec"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
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

		ros, extraResources, err := spec.CoreOperations(input.data)
		if err != nil {
			return errors.Wrap(err, "unable to perform controller operations")
		}

		// write all the kubernetes objects that were generated
		for _, runtimeObject := range ros {

			data, err := yaml.Marshal(runtimeObject)
			if err != nil {
				return errors.Wrap(err, "failed to marshal object")
			}

			err = writeObject(data)
			if err != nil {
				return errors.Wrap(err, "failed to write object")
			}
		}

		for _, file := range extraResources {
			// change the file name to absolute file name
			// then read the file and then pass it to writeObject
			file = findAbsPath(input.fileName, file)
			data, err := ioutil.ReadFile(file)
			if err != nil {
				return errors.Wrap(err, "file reading failed")
			}
			err = writeObject(data)
			if err != nil {
				return errors.Wrap(err, "failed to write object")
			}
		}
	}
	return nil
}

func writeObject(data []byte) error {
	_, err := fmt.Fprintln(os.Stdout, "---")
	if err != nil {
		return errors.Wrap(err, "could not print to STDOUT")
	}

	_, err = os.Stdout.Write(data)
	return errors.Wrap(err, "could not write to STDOUT")
}

func findAbsPath(baseFilePath, path string) string {
	// TODO: if the baseFilePath is empty then just take the
	// pwd as basefilePath, here we will force user to
	// use the kedge binary from the directory that has files
	// otherwise there is no way of knowing where the files will be
	// this condition will happen when we add support for reading from the stdin
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(filepath.Dir(baseFilePath), path)
}
