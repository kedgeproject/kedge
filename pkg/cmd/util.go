package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"time"
)

type inputData struct {
	fileName string
	data     []byte
}

const (
	retryAttempts = 3
)

func getApplicationsFromFiles(files []string) ([]inputData, error) {
	var appData []inputData
	var data []byte
	var err error
	for _, file := range files {
		if checkIfURL(file) {
			data, err = GetURLData(file, retryAttempts)
		} else {
			data, err = ioutil.ReadFile(file)
			if err != nil {
				return nil, errors.Wrap(err, "file reading failed")
			}
			file, err := filepath.Abs(file)
			if err != nil {
				return nil, errors.Wrapf(err, "cannot determine the absolute file path of %q", file)
			}
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
		if checkIfURL(path) {
			files = append(files, path)
		} else {
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
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("no *.yml or *.yaml files were found")
	}
	return files, nil
}

// This function validates the URL, then tries to fetch the data with retries
// and then reads and returns the data as []byte
// Returns an error if the URL is invalid, or fetching the data failed or
// if reading the response body fails.
func GetURLData(urlString string, attempts int) ([]byte, error) {
	// Validate URL
	_, err := url.ParseRequestURI(urlString)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %v", urlString)
	}

	// Fetch the URL and store the response body
	data, err := FetchURLWithRetries(urlString, attempts, 1*time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed fetching data from the URL %v: %s", urlString, err)
	}

	return data, nil
}

// Try to fetch the given url string, and make <attempts> attempts at it.
// Wait for <duration> time between each try.
// This returns the data from  the response body []byte upon successful fetch
// The passed URL is not validated, so validate the URL before passing to this
// function
func FetchURLWithRetries(url string, attempts int, duration time.Duration) ([]byte, error) {
	var data []byte
	var err error

	for i := 0; i < attempts; i++ {
		var response *http.Response

		// sleep for <duration> seconds before trying again
		if i > 0 {
			time.Sleep(duration)
		}

		// retry if http.Get fails
		// if all the retries fail, then return statement at the end of the
		// function will return this err received from http.Get
		response, err = http.Get(url)
		if err != nil {
			continue
		}
		defer response.Body.Close()

		// if the status code is not 200 OK, return an error
		if response.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("unable to fetch %v, server returned status code %v", url, response.StatusCode)
		}

		// Read from the response body, ioutil.ReadAll will return []byte
		data, err = ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, fmt.Errorf("reading from response body failed: %s", err)
		}
		break
	}

	return data, err
}

// checkIfURL will return True if given path is an URL
func checkIfURL(path string) bool {
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return true
	}
	return false

}
