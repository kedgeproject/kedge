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
	"path/filepath"

	"fmt"
	"log"

	"github.com/pkg/errors"
)

// Common global variables being used for kedge subcommands are declared here.
// Before adding anything here, make sure that the subcommands using these
// variables are mutually exclusive.
// e.g. only one of `kedge generate` or `kedge create` can be run at a time,
// so it makes sense to use the common InputFiles variable in both of those
// commands.
var (
	InputFiles     []string
	Namespace      string
	SkipValidation bool
)

func ifFilesPassed(files []string) error {
	if len(files) == 0 {
		return errors.New("No files were passed. Please pass file(s) using '-f' or '--files'")
	}
	return nil
}

// removeDuplicateFiles function will see input files and will remove duplicate entries
func removeDuplicateFiles(inputFiles []string) []string {
	files := map[string]bool{}
	fileSet := []string{}
	for v := range inputFiles {
		absPath, err := filepath.Abs(inputFiles[v])
		if err != nil {
			log.Fatalf("Given file %s is not present", inputFiles[v])
		}
		if !files[absPath] {
			files[absPath] = true
			fileSet = append(fileSet, inputFiles[v])
		}
	}
	return fileSet
}

// UsageError will show errors with suggestion to see help
func usageError(cmdPath string, err error) {
	fmt.Printf("%s\nSee '%s -h' for help and examples.\n", err, cmdPath)
}
