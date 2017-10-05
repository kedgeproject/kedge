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
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

type inputData struct {
	fileName string
	data     []byte
}

func getApplicationsFromFiles(files []string) ([]inputData, error) {
	var appData []inputData

	for _, file := range files {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, errors.Wrap(err, "file reading failed")
		}
		file, err := filepath.Abs(file)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot determine the absolute file path of %q", file)
		}

		// The regular expression takes care of when triple dashes are in the
		// starting of the file or when they are as a separate line somewhere
		// in the middle of the file or at the end. Ideally this should be taken
		// care by the yaml library since this is valid YAML syntax anyway,
		// but right now neither the go-yaml/yaml (issue #232) nor the one
		// that we use supports the multiple document structure, so yeah!
		apps := regexp.MustCompile("(^|\n)---\n").Split(string(data), -1)
		for _, app := range apps {
			// strings.TrimSpace will remove all the extra whitespaces and
			// newline characters, and then proceed only when the length of the
			// string is more than 0
			// this avoids passing empty input further in the program in cases
			// like -
			// ---			# avoids empty input here
			// ---
			// name: abc
			// containers:
			// ...
			// ---
			//				# avoids empty input here
			// ---			# avoids empty input here
			if len(strings.TrimSpace(app)) > 0 {
				appData = append(appData, inputData{
					fileName: file,
					data:     []byte(app),
				})
			}
		}
	}
	return appData, nil
}

// GetAllYAMLFiles if path in argument is directory get all *.yml and *.yaml files
// in that directory. If path is file just add it to output list as it is.
func GetAllYAMLFiles(paths []string) ([]string, error) {
	var files []string
	for _, path := range paths {
		fileInfo, err := os.Stat(path)
		if err != nil {
			return nil, errors.Wrapf(err, "can't get file info about %s", path)
		}
		if fileInfo.IsDir() {
			ymlFiles, err := filepath.Glob(filepath.Join(path, "*.yml"))
			if err != nil {
				return nil, errors.Wrapf(err, "can't list *.yml files in %s", path)
			}
			files = append(files, ymlFiles...)
			yamlFiles, err := filepath.Glob(filepath.Join(path, "*.yaml"))
			if err != nil {
				return nil, errors.Wrapf(err, "can't list *.yaml files in %s", path)
			}
			files = append(files, yamlFiles...)
		} else {
			// path is regular file, do nothing and just add it to list of files
			files = append(files, path)
		}
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("no *.yml or *.yaml files were found")
	}
	return files, nil
}

// SubstituteVariables handles environment variable substitution
// Expression [[ NAME ]] or [[NAME]] is used for defining variables.
// [[NAME]] will be replaced with the value from NAME environment variable.
// Returns new byte array with all variables replaced with values.
// If no matching variable is found, an error is raised.
func SubstituteVariables(kedgeFileContent []byte) ([]byte, error) {
	// all occurrences  [[ something ]]
	// (there can be any number of white characters after '[[' and before ']]' )
	re := regexp.MustCompile("\\[\\[\\s*(.*?)\\s*\\]\\]")
	newContent := re.ReplaceAllFunc(kedgeFileContent, replaceWithEnv)

	// check for any remaining variables in newContent
	remains := re.FindAllSubmatch(newContent, -1)
	remainingVariables := []string{}

	for _, r := range remains {
		name := string(r[1][:])
		remainingVariables = append(remainingVariables, name)
	}

	if len(remainingVariables) > 0 {
		return nil, fmt.Errorf("undefined variable(s): %s", strings.Join(remainingVariables, ", "))
	}

	return newContent, nil
}

// replaceWithEnv is used with regexp.ReplaceAllFunc to replace variable
// with value from environment
func replaceWithEnv(in []byte) []byte {
	// get just a name from variable (withot '[[' and ']]')
	re := regexp.MustCompile("\\[\\[\\s*(.*?)\\s*\\]\\]")
	groups := re.FindSubmatch(in)
	// index 0 contains full match, 1 contains first group
	name := string(groups[1][:])

	value, found := os.LookupEnv(name)
	if !found {
		// If there is no corresponding env variable, return original string.
		// We will use this to detect remaining variables in the input file.
		// We can't return error from this function.
		return in
	}
	return []byte(value)
}
