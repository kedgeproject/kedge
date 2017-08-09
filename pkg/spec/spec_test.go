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

package spec

import (
	"testing"

	"fmt"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

const (
	kedgeSpecPackagePath = "github.com/kedgeproject/kedge/pkg/spec"
)

func TestConflictingFields(t *testing.T) {
	// This array should ideally go away and every struct in spec.go should
	// get tested automatically, but for now, let's pass in all the structs
	// https://stackoverflow.com/questions/20803758/how-to-get-all-defined-struct-in-golang
	structsToTest := []interface{}{
		&VolumeClaim{},
		&ServicePortMod{},
		&ServiceSpecMod{},
		&IngressSpecMod{},
		&Container{},
		&ConfigMapMod{},
		&PodSpecMod{},
		&App{},
	}

	for _, inputStruct := range structsToTest {
		t.Run("Testing conflicting fields", func(t *testing.T) {

			// Checking if input is pointer to struct
			if err := checkTypePointerToStruct(inputStruct); err != nil {
				t.Error(errors.Wrap(err, "Input parameter type mismatch"))
			}

			conflictingTags, err := findConflictingJSONTags(reflect.ValueOf(inputStruct))
			if err != nil {
				t.Error(errors.Wrap(err, "Unable to find conflicting tags for spec.App"))
			}
			if len(conflictingTags) != 0 {
				t.Logf("The struct %v has unhandled conflicting JSON tags which exist in other structs.", reflect.Indirect(reflect.ValueOf(inputStruct)).Type().String())
				for tag, structs := range conflictingTags {
					t.Logf("The JSON tag '%v' exists in %v", tag, structs)
				}
				t.Fatal("Once you handle the above conflicting JSON tag, mark it as handled by adding a 'conflicting' JSON tag to its definition. e.g.\nContainers     []Container `json:\"containers,conflicting,omitempty\"`")
			}
		})
	}
}

// This function takes a map as an input which stores the JSON tag as key and
// the names of the structs holding that tag as the value. The passed in
// blacklisted tags are ignored and are not populated at all.
func getUnmarshalJSONTagsMap(inputMap map[string][]string, inputStruct reflect.Value, blacklistedTags []string) error {

	val := inputStruct.Elem()

StructFieldsLoop:
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)

		// We do not need to store JSON tags like "inline" or "omitempty",
		// instead we only need JSON tag which is used for unmarshalling the
		// YAML. It's safe to assume that, that JSON tag is the one that
		// is specified first.
		// Let's split the JSON tag on ",", and store the first element of the
		// resulting array.
		unmarshalJSONTag := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]

		if unmarshalJSONTag == "" {
			// It does not make any sense to store anything for JSON tags with no
			// first element, e.g. `json:",inline"`
			continue
		}

		// Populating only the non blacklisted tags
		for _, blacklistedTag := range blacklistedTags {
			if unmarshalJSONTag == blacklistedTag {
				continue StructFieldsLoop
			}
		}

		inputMap[unmarshalJSONTag] = append(inputMap[unmarshalJSONTag], val.Type().String())
	}
	return nil
}

// This function returns an array of pointers to all of the embedded structs
// in the input struct.
// The embedded structs in an embedded struct are also taken into consideration
// in a recursive manner.
func getEmbeddedStructs(inputStruct reflect.Value) ([]reflect.Value, error) {

	var embeddedStructs []reflect.Value

	val := inputStruct.Elem()

	for i := 0; i < val.NumField(); i++ {
		// Checking for embedded structs
		if val.Type().Field(i).Anonymous {

			// Appending pointer to the embedded/anonymous struct to
			// embeddedStructs
			embeddedStructs = append(embeddedStructs, val.Field(i).Addr())

			// Since the current field is an anonymous struct, we call the
			// function recursively to get embedded structs under it
			recursiveEmbeddedStructs, err := getEmbeddedStructs(val.Field(i).Addr())
			if err != nil {
				return nil, errors.Wrapf(err, "Unable to get embedded structs recursively from %v", val.Field(i).String())
			}

			// Adding the returned recursive structs to embeddedStructs
			embeddedStructs = append(embeddedStructs, recursiveEmbeddedStructs...)
		}
	}
	return embeddedStructs, nil
}

// This function returns an array of all the JSON tags used for unmarshalling
// (which are the first ones specified) which have a JSON tag "conflicting", in
// all of the input structs.
// It's made sure that only the tags in spec.go are checked and all other
// packages are ignored. This is done because we are expecting the structs to
// be handled, i.e. marked as "conflicting" only in spec.go, and we do not want
// to populate the list with JSON tags from some other package which had the
// "conflicting" JSON tag for some other reason
func getMarkedAsConflictingJSONUnmarshalTags(inputStructs []reflect.Value) ([]string, error) {
	var blacklistedTags []string
	for _, inputStruct := range inputStructs {

		val := inputStruct.Elem()

		// Proceeding only if the struct belongs to the kedge's spec package
		if val.Type().PkgPath() == kedgeSpecPackagePath {
			for i := 0; i < val.NumField(); i++ {
				field := val.Type().Field(i)

				fieldJSONTags := strings.Split(field.Tag.Get("json"), ",")

				// Adding fields that have "conflicting" as a JSON tag to the
				// array. We need not check the first tag that is used for
				// unmarshalling.
				for _, tag := range fieldJSONTags[1:] {
					if tag == "conflicting" {
						blacklistedTags = append(blacklistedTags, fieldJSONTags[0])
						break
					}
				}
			}
		}
	}
	return blacklistedTags, nil
}

// This function checks if the input parameter has the pointer to a struct as
// its underlying type.
// Returns "nil" if true, or an error if false
func checkTypePointerToStruct(input interface{}) error {
	if reflect.Indirect(reflect.ValueOf(input)).Kind().String() != "struct" {
		return fmt.Errorf("Expected pointer to struct, got %v", reflect.ValueOf(input).Kind().String())
	}
	return nil
}

// This function takes in pointer to a struct as input, and returns a map of
// conflicting JSON tags used for unmarshaling.
// All of the embedded structs are taken into account, and the fields with a
// JSON tag "conflicting" are assumed to be handled.
// The returned map contains the JSON tags as the keys, and the values are an
// array of struct names which contain those fields, without being handled.
func findConflictingJSONTags(inputStruct reflect.Value) (map[string][]string, error) {

	// We need to check that no JSON tag is specified more than once at the
	// top level of the struct.
	// This means that the struct fields, and the JSON tags of all the
	// fields of all the embedded structs need to be unique.
	// We accomplish this over a 4 step process.

	// Step 1: Get all the target structs.
	// Look at tags in the fields of the input struct as well as the embedded
	// structs, so both of them become our target structs.

	var targetStructs []reflect.Value
	// The first target struct is the input struct to the function
	targetStructs = append(targetStructs, inputStruct)

	// Getting all embedded structs for the input struct and appending to the
	// list of target structs.
	embeddedStructs, err := getEmbeddedStructs(inputStruct)
	if err != nil {
		return nil, errors.Wrapf(err, "Error getting embedded structs for struct: %v", inputStruct)
	}
	targetStructs = append(targetStructs, embeddedStructs...)

	// Step 2: Get blacklisted tags.
	// Get the fields which have already been handled and were previously
	// conflicting. This needs to be done only for the structs in spec.go,
	// and with fields marked as "conflicting" using a JSON tag

	blacklistedTags, err := getMarkedAsConflictingJSONUnmarshalTags(targetStructs)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get blacklisted tags for target structs")
	}

	// Step 3: Get JSON tags from target structs.
	// Get the JSON tags from all of the target structs which we got in Step 1
	// except for the blacklisted tags, along with the names of the structs
	// which hold the field with that tag.

	tagStructMap := make(map[string][]string)
	for _, targetStruct := range targetStructs {
		err := getUnmarshalJSONTagsMap(tagStructMap, targetStruct, blacklistedTags)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to get JSON tags from the struct: %v", inputStruct)
		}
	}

	// Step 4: Check if more than one struct holds a tag.
	// Delete all the tags in Step 3 which are only help by one struct, because
	// it means that the JSON tag is not conficting, and finally return the
	// remaining tags which are held by more than one struct, and are not
	// handled, hence conflicting tags

	for tag, containingStructs := range tagStructMap {
		if len(containingStructs) == 1 {
			delete(tagStructMap, tag)
		}
	}

	return tagStructMap, nil
}
