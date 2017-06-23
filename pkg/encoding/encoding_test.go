package encoding

import (
	"testing"

	"github.com/surajssd/kapp/pkg/encoding/fixtures"
	"github.com/surajssd/kapp/pkg/spec"

	"reflect"

	"github.com/davecgh/go-spew/spew"
)

func TestDecode(t *testing.T) {
	tests := []struct {
		Name string
		Data []byte
		App  *spec.App
	}{
		{
			Name: "One container mentioned in the spec",
			Data: fixtures.SingleContainer,
			App:  &fixtures.SingleContainerApp,
		},
		{
			Name: "One persistent volume mentioned in the spec",
			Data: fixtures.SinglePersistentVolume,
			App:  &fixtures.SinglePersistentVolumeApp,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			app, err := Decode(test.Data)
			if err != nil {
				t.Fatalf("Unable to run Decode(), and error occurred: %v", err)
			}

			if !reflect.DeepEqual(test.App, app) {
				t.Fatalf("Expected:\n%v\nGot:\n%v", spew.Sprint(test.App), spew.Sprint(app))
			}
		})
	}
}
