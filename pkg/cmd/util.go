package cmd

import (
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

type inputData struct {
	data []byte
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
					data: []byte(app),
				})
			}
		}
	}
	return appData, nil
}
