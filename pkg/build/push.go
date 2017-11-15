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

package build

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"
)

/*
PushImage push a Docker image via the docker client. Takes the image name
as input.
*/
func PushImage(fullImageName string) error {
	log.Infof("Pushing image %q", fullImageName)

	command := []string{"docker", "push", fullImageName}
	if err := RunCommand(command); err != nil {
		return err
	}
	log.Infof("Successfully pushed image %q", fullImageName)

	return nil
}

func RunCommand(command []string) error {
	cmd := exec.Command(command[0], command[1:]...)
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("%s, %s", strings.TrimSpace(stderr.String()), err)
	}

	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("%s, %s", strings.TrimSpace(stderr.String()), err)
	}

	return nil
}
