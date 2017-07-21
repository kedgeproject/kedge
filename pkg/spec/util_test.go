package spec

import (
	"reflect"
	"testing"
)

// No conflicting field
type MasterStruct1 struct {
	Ms1 string  `json:"ms1"`
	Ms2 int     `json:"ms2"`
	Ms3 float64 `json:"ms3"`
	EmbeddedStruct1
	EmbeddedStruct2
}

// Conflicting field with EmbeddedStruct1
type MasterStruct2 struct {
	Ms1  string  `json:"ms1"`
	Ms2  int     `json:"ms2"`
	Ms3  float64 `json:"ms3"`
	Es11 int     `json:"es11"`
	EmbeddedStruct1
	EmbeddedStruct2
}

// No conflicting field with embedded structs, but conflicting fields in embedded structs
type MasterStruct3 struct {
	Ms1 string  `json:"ms1"`
	Ms2 int     `json:"ms2"`
	Ms3 float64 `json:"ms3"`
	EmbeddedStruct3
	EmbeddedStruct4
}

type EmbeddedStruct1 struct {
	Es11 string  `json:"es11"`
	Es12 int     `json:"es12"`
	Es13 float64 `json:"es13"`
}

type EmbeddedStruct2 struct {
	Es21 string  `json:"es21"`
	Es22 int     `json:"es22"`
	Es23 float64 `json:"es23"`
}

type EmbeddedStruct3 struct {
	Es31 string  `json:"es31"`
	Es32 int     `json:"es32"`
	Es33 float64 `json:"es33"`
}

type EmbeddedStruct4 struct {
	Es31 string  `json:"es31"` // conflicting field
	Es41 string  `json:"es41"`
	Es42 int     `json:"es42"`
	Es43 float64 `json:"es43"`
}

// Conflicting fields in this struct and both the embedded structs
type MasterStruct4 struct {
	Ms1 string  `json:"ms1"` // conflicting field
	Ms2 int     `json:"ms2"`
	Ms3 float64 `json:"ms3"`
	EmbeddedStruct5
	EmbeddedStruct6
}

// Conflicting fields in this struct and both the embedded structs, with
// "conflicting" JSON tag
type MasterStruct5 struct {
	Ms1 string  `json:"ms1,conflicting"` // conflicting field
	Ms2 int     `json:"ms2"`
	Ms3 float64 `json:"ms3"`
	EmbeddedStruct5
	EmbeddedStruct6
}

type EmbeddedStruct5 struct {
	MasterField1 int     `json:"ms1"` // conflicting field
	Es51         string  `json:"es51"`
	Es52         int     `json:"es52"`
	Es53         float64 `json:"es53"`
}

type EmbeddedStruct6 struct {
	MasterFoo1 float64 `json:"ms1"` // conflicting field
	Es63       float64 `json:"es63"`
}

func TestFindConflictingYAMLTags(t *testing.T) {
	tests := []struct {
		Name            string
		InputStuct      interface{}
		ConflictingTags map[string][]string
		IsError         bool
	}{
		{
			Name:            "Pointer to struct not passed",
			InputStuct:      []int{42},
			ConflictingTags: nil,
			IsError:         true,
		},
		{
			Name:            "No conflicting field",
			InputStuct:      &MasterStruct1{},
			ConflictingTags: map[string][]string{},
			IsError:         false,
		},
		{
			Name:       "Conflicting field in top level struct and embedded struct",
			InputStuct: &MasterStruct2{},
			ConflictingTags: map[string][]string{
				"es11": {"spec.MasterStruct2", "spec.EmbeddedStruct1"},
			},
			IsError: false,
		},
		{
			Name:       "Conflicting field in embedded structs but not in top level struct",
			InputStuct: &MasterStruct3{},
			ConflictingTags: map[string][]string{
				"es31": {"spec.EmbeddedStruct3", "spec.EmbeddedStruct4"},
			},
			IsError: false,
		},
		{
			Name:       "Conflicting field in top level and all embedded structs",
			InputStuct: &MasterStruct4{},
			ConflictingTags: map[string][]string{
				"ms1": {"spec.MasterStruct4", "spec.EmbeddedStruct5", "spec.EmbeddedStruct6"},
			},
			IsError: false,
		},
		{
			Name:            "Conflicting tags should be empty when 'conflicting' JSON tag present",
			InputStuct:      &MasterStruct5{},
			ConflictingTags: map[string][]string{},
			IsError:         false,
		},

		// TODO: Conflicting tags should not be empty when 'conflicting' JSON tag present in embedded struct not in "spec" package
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			conflictingTags, err := findConflictingJSONTags(test.InputStuct)

			// Testing errors
			if test.IsError && err == nil {
				t.Fatal("Expected function to return an error, but no error returned")
			} else if !test.IsError && err != nil {
				t.Fatalf("No error expected, but got %v", err)
			}

			// Testing conflicting tags
			if !reflect.DeepEqual(conflictingTags, test.ConflictingTags) {
				t.Fatalf("Expected conflicting tags to be:\n%v\nBut got:\n%v\n", test.ConflictingTags, conflictingTags)
			}
		})
	}
}
