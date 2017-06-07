package main

import (
	"os"

	"github.com/surajssd/opencomposition/cmd"
)

func main() {
	cmd := cmd.NewRootCommand()
	if err := cmd.Execute(); err != nil {
		//log.Println(err)
		os.Exit(-1)
	}
}
