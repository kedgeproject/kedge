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
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

var matchFiles = []string{
	"test.yml",
	"foo/bar.yml",
	"foo/bar.yaml",
}

// createMatchFiles crates empty files in tmpDir as defined in matchFiles
func createMatchFiles(tmpDir string) error {
	for _, file := range matchFiles {
		fileName := filepath.Join(tmpDir, file)
		// check if parrent directory exists
		if _, err := os.Stat(filepath.Dir(fileName)); os.IsNotExist(err) {
			// create parrent directory
			err = os.MkdirAll(filepath.Dir(fileName), os.ModePerm)
			if err != nil {
				return err
			}
		}
		// create empty files
		_, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0666)
		if err != nil {
			return err
		}
	}
	return nil
}

var matchTests = []struct {
	paths  []string
	result []string
}{
	{[]string{"test.yml"}, []string{"test.yml"}},
	{[]string{"foo"}, []string{"foo/bar.yml", "foo/bar.yaml"}},
	{[]string{"test.yml", "foo"}, []string{"test.yml", "foo/bar.yml", "foo/bar.yaml"}},
}

func TestGetAllYMLFiles(t *testing.T) {

	// create temporary dir where all test files will be created
	tmpDir, err := ioutil.TempDir("", "matchTest")
	if err != nil {
		t.Fatal("creating temp dir:", err)
	}
	defer os.RemoveAll(tmpDir)

	createMatchFiles(tmpDir)

	for _, test := range matchTests {
		var paths []string
		var result []string
		// prefix all test path with tmpDir
		for _, p := range test.paths {
			paths = append(paths, filepath.Join(tmpDir, p))
		}
		// prefix all expected results with tmpDir
		for _, p := range test.result {
			result = append(result, filepath.Join(tmpDir, p))
		}

		out, err := GetAllYAMLFiles(paths)
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(out, result) {
			t.Errorf("output doesn't match expected output\n output  : %#v \n expected: %#v \n", out, result)
		}
	}

}
