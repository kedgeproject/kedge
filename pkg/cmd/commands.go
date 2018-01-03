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
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/kedgeproject/kedge/pkg/spec"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
)

// GenerateArtifacts either writes to file or uses kubectl/oc to deploy.
// TODO: Refactor into two separate functions (remove `generate bool`).
func CreateArtifacts(paths []string, generate bool, args ...string) error {

	files, err := GetAllYAMLFiles(paths)
	if err != nil {
		return errors.Wrap(err, "unable to get YAML files")
	}

	inputs, err := getApplicationsFromFiles(files)
	if err != nil {
		return errors.Wrap(err, "unable to get kedge definitions from input files")
	}

	for _, input := range inputs {

		// Substitute variables
		// We do this on raw Kedge file before unmarshalling, because it would be
		// complicated to go through all different go structs.
		kedgeData, err := SubstituteVariables(input.data)
		if err != nil {
			return errors.Wrap(err, "failed to replace variables")
		}

		ros, includeResources, err := spec.CoreOperations(kedgeData)
		if err != nil {
			return errors.Wrap(err, "unable to perform controller operations")
		}

		// decide between kubectl and oc
		useOC := false
		for _, runtimeObject := range ros {
			switch runtimeObject.GetObjectKind().GroupVersionKind().Kind {
			// If there is at least one OpenShift resource use oc
			case "DeploymentConfig", "Route", "ImageStream", "BuildConfig":
				useOC = true
				break
			}
		}

		for _, runtimeObject := range ros {

			// Unmarshal said object
			data, err := yaml.Marshal(runtimeObject)
			if err != nil {
				return errors.Wrap(err, "failed to marshal object")
			}

			// Write to file if generate = true
			if generate {
				err = writeObject(data)
				if err != nil {
					return errors.Wrap(err, "failed to write object")
				}
			} else {
				// We need to add "-f -" at the end of the command passed to us to
				// pass the generated files.
				// e.g. If the command and arguments are "apply --namespace staging", then the
				// final command becomes "kubectl apply --namespace staging -f -"
				arguments := append(args, "-f", "-")
				err = RunClusterCommand(arguments, data, useOC)
				if err != nil {
					return errors.Wrap(err, "failed to execute command")
				}
			}

		}

		for _, file := range includeResources {
			// change the file name to absolute file name
			file = findAbsPath(input.fileName, file)

			if generate {
				data, err := ioutil.ReadFile(file)
				if err != nil {
					return errors.Wrap(err, "file reading failed")
				}
				err = writeObject(data)
				if err != nil {
					return errors.Wrap(err, "failed to write object")
				}
			} else {

				// We need to add "-f absolute-filename" at the end of the command passed to us to
				// pass the generated files.
				// e.g. If the command and arguments are "apply --namespace staging", then the
				// final command becomes "kubectl apply --namespace staging -f absolute-filename"
				arguments := append(args, "-f", file)
				err = RunClusterCommand(arguments, nil, useOC)
				if err != nil {
					return errors.Wrap(err, "failed to execute command")
				}
			}
		}
	}
	return nil
}

// runClusterCommand calls kubectl or oc binary.
// Boolean flag useOC controls if oc or kubectl will be used
func RunClusterCommand(args []string, data []byte, useOC bool) error {

	// Use kubectl by default, oc if useOC bool is true (in cases such as DeploymentConfig, ImageStream, etc.)
	executable := "kubectl"
	if useOC {
		executable = "oc"
	} else {
		if _, err := exec.LookPath("kubectl"); err != nil {
			log.Debug("kubectl is unavailable, using oc")
			executable = "oc"
		}
	}

	// If oc is used, error out if it's not available
	if executable == "oc" {
		if _, err := exec.LookPath("oc"); err != nil {
			return errors.New("Unable to find oc command. Please install oc to your system")
		}
	}

	// Create the executable command
	cmd := exec.Command(executable, args...)

	// Read from stdin
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return errors.Wrap(err, "can't get stdinPipe for kubectl")
	}

	// Write to stdin
	go func() {
		defer stdin.Close()
		_, err := io.WriteString(stdin, string(data))
		if err != nil {
			fmt.Printf("can't write to stdin %v\n", err)
		}
	}()

	// Execute the actual command
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("%s", string(out))
		return errors.Wrap(err, "failed to execute command")
	}

	fmt.Printf("%s", string(out))
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
