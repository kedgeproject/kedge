package main

import (
	"fmt"
	"io/ioutil"
	//"github.com/ghodss/yaml"
	"github.com/tidwall/gjson"
)

func main() {
	filename := "swagger.json"

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("cannot read file %q: %v\n", filename, err)
	}

	//y, err := yaml.JSONToYAML(content)
	//fmt.Printf("%s", string(y))

	fmt.Println(gjson.Get(string(content), "definitions/io.k8s.apimachinery.pkg.util.intstr.IntOrString"))
}
