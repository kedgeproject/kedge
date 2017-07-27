package cmd

import "github.com/pkg/errors"

// Common global variables being used for kedge subcommands are declared here.
// Before adding anything here, make sure that the subcommands using these
// variables are mutually exclusive.
// e.g. only one of `kedge generate` or `kedge create` can be run at a time,
// so it makes sense to use the common InputFiles variable in both of those
// commands.
var (
	InputFiles []string
	Namespace  string
)

func ifFilesPassed(files []string) error {
	if len(files) == 0 {
		return errors.New("No files were passed. Please pass file(s) using '-f' or '--files'")
	}
	return nil
}
