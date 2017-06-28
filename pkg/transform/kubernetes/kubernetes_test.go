package kubernetes

import (
	"testing"

	"reflect"

	encodingFixtures "github.com/kedgeproject/kedge/pkg/encoding/fixtures"
	"github.com/kedgeproject/kedge/pkg/spec"
	transformFixtures "github.com/kedgeproject/kedge/pkg/transform/fixtures"
	"k8s.io/client-go/pkg/runtime"
)

func TestCreateServices(t *testing.T) {
	tests := []struct {
		Name    string
		App     *spec.App
		Objects []runtime.Object
	}{
		{
			"Single container specified",
			&encodingFixtures.SingleContainerApp,
			append(make([]runtime.Object, 0), transformFixtures.SingleContainerService),
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			object, err := createServices(test.App)
			if err != nil {
				t.Fatalf("Creating services failed: %v", err)
			}
			if !reflect.DeepEqual(test.Objects, object) {
				t.Fatalf("Expected:\n%v\nGot:\n%v", test.Objects, object)
			}
		})
	}
}

// TODO: add test for auto naming of single persistent volume
