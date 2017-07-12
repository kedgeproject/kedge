package cmd

import "github.com/pkg/errors"

func ifFilesPassed(files []string) error {
	if len(files) == 0 {
		return errors.New("No files were passed. Please pass file(s) using '-f' or '--files'")
	}
	return nil
}
