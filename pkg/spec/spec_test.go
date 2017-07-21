package spec

import (
	"testing"

	"github.com/pkg/errors"
	"reflect"
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
		&EnvFromSource{},
		&ConfigMapEnvSource{},
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

			conflictingTags, err := findConflictingJSONTags(inputStruct)
			if err != nil {
				t.Error(errors.Wrap(err, "Unable to find conflicting tags for spec.App"))
			}
			if len(conflictingTags) != 0 {
				t.Fatalf("The struct %v has unhandled conflicting JSON tags which exist in other structs.\n%v", reflect.Indirect(reflect.ValueOf(inputStruct)).Type().String(), conflictingTags)
			}
		})
	}
}
