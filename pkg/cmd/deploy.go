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
	"os/exec"

	"github.com/kedgeproject/kedge/pkg/encoding"
	"github.com/kedgeproject/kedge/pkg/transform/kubernetes"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
)

func Deploy(files []string) error {
	for _, file := range files {

		data, err := ioutil.ReadFile(file)
		if err != nil {
			return errors.Wrap(err, "file reading failed")
		}

		app, err := encoding.Decode(data)
		if err != nil {
			return errors.Wrap(err, "unable to unmarshal data")
		}

		ros, err := kubernetes.Transform(app)
		if err != nil {
			return errors.Wrap(err, "unable to convert data")
		}

		for _, o := range ros {
			data, err := yaml.Marshal(o)
			if err != nil {
				return errors.Wrap(err, "failed to marshal object")
			}

			cmd := exec.Command("kubectl", "create", "-f", "-")

			stdin, err := cmd.StdinPipe()
			if err != nil {
				return errors.Wrap(err, "can't get stdinPipe for kubectl")
			}

			go func() {
				defer stdin.Close()
				_, err := io.WriteString(stdin, string(data))
				if err != nil {
					fmt.Printf("can't write to stdin %v\n", err)
				}
			}()

			out, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Printf("%s", string(out))
				return errors.Wrap(err, "failed to execute command")
			}
			fmt.Printf("%s", string(out))
		}
	}

	return nil
}
