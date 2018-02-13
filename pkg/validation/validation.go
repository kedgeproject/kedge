package validation

import (
	"fmt"
	"log"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"github.com/xeipuuv/gojsonschema"
)

// https://github.com/xeipuuv/gojsonschema#formats
type ValidFormat struct{}

// IsFormat always returns true and meets the
// gojsonschema.FormatChecker interface
func (f ValidFormat) IsFormat(input interface{}) bool {
	return true
}

// Based on https://stackoverflow.com/questions/40737122/convert-yaml-to-json-without-struct-golang
// We unmarshal yaml into a value of type interface{},
// go through the result recursively, and convert each encountered
// map[interface{}]interface{} to a map[string]interface{} value
// required to marshall to JSON.
// Reference: https://github.com/garethr/kubeval/blob/master/kubeval/utils.go#L8
func convertToStringKeys(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = convertToStringKeys(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = convertToStringKeys(v)
		}
	}
	return i
}

// validate function will validate input kedgefile against JSON schema provided in pkg/validation
func Validate(p []byte) {
	// To detect extra field
	// if we let yaml marshaller first deserialize this, all extra fields will be ignored so,
	// we are unmarshalling input into interface so that we can detect extra fields
	var speco interface{}
	err := yaml.Unmarshal(p, &speco)
	if err != nil {
		fmt.Printf("Error with Unmarhsalling")
	}
	body := convertToStringKeys(speco)
	loader := gojsonschema.NewGoLoader(body)

	var schema gojsonschema.JSONLoader
	// we are retrieving SchemaJSON from `pkg/validation/kedgeschema.go`
	// kedgeschema.go is automatically generated when using `make update-schema`
	schema = gojsonschema.NewStringLoader(SchemaJson)
	// Without forcing these types the schema fails to load
	//Reference: https://github.com/xeipuuv/gojsonschema#formats
	gojsonschema.FormatCheckers.Add("int64", ValidFormat{})
	gojsonschema.FormatCheckers.Add("byte", ValidFormat{})
	gojsonschema.FormatCheckers.Add("int32", ValidFormat{})
	gojsonschema.FormatCheckers.Add("int-or-string", ValidFormat{})
	result, err := gojsonschema.Validate(schema, loader)
	if err != nil {
		errors.Wrap(err, "Error while validating schema")
	}
	if !result.Valid() {
		error := ""
		for _, err := range result.Errors() {
			error = error + err.String() + "\n"
		}
		log.Fatal("The kedgefile is not valid. see errors :\n", error)
	}

}
