package cmd

import (
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

func getApplicationsFromFiles(files []string) ([][]byte, error) {
	var appData [][]byte
	for _, file := range files {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, errors.Wrap(err, "file reading failed")
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
				appData = append(appData, []byte(app))
			}
		}
	}
	return appData, nil
}
