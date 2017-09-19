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
	"os"

	"text/template"

	"github.com/spf13/cobra"
)

type templateData struct {
	Name  string
	Image string
	Ports []int
}

var (
	fileName, image, name string
	port                  []int
)

const (
	boilerplate = `name: {{.Name}}
containers:
- image: {{.Image}}
{{if .Ports}}services:
- ports:{{block "list" .Ports}}
{{range .}}{{print "  - port: " .}}
{{end}}{{end}}{{end}}`
)

// initCmd represents the version command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a Kedge file",
	Run: func(cmd *cobra.Command, args []string) {
		data := templateData{Name: name, Image: image, Ports: port}
		if name != "" && image != "" {
			masterTmpl, err := template.New("master").Parse(boilerplate)
			if err != nil {
				fmt.Println("failed to create template")
				os.Exit(-1)
			}
			_, err = os.Stat(fileName)
			if err != nil {
				f, err := os.Create(fileName)
				if err != nil {
					fmt.Println(err, "failed to create file")
					os.Exit(-1)
				}
				err = masterTmpl.Execute(f, data)
				if err != nil {
					fmt.Println(err, "failed to write file")
					os.Exit(-1)
				}
				defer f.Close()
				fmt.Println("file", fileName, "created")
			} else {
				fmt.Println(fileName, "is already present")
				os.Exit(-1)
			}
		} else {
			fmt.Println("--name and --image are mandatory flags, Please provide these flags")
			os.Exit(-1)
		}
	},
}

func init() {
	initCmd.Flags().StringVarP(&fileName, "out", "o", "kedge.yml", "Output filename")
	initCmd.Flags().StringVarP(&name, "name", "n", "", "The name of service")
	initCmd.Flags().StringVarP(&image, "image", "i", "", "The image for the container to run")
	initCmd.Flags().IntSliceVarP(&port, "port", "p", []int{}, "The ports that this container exposes")
	RootCmd.AddCommand(initCmd)
}
