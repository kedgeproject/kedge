package encoding

import (
	"testing"

	"github.com/kedgeproject/kedge/pkg/spec"
)

func TestValidateVolumeClaims(t *testing.T) {

	failingTest := []spec.VolumeClaim{{Name: "foo"}, {Name: "bar"}, {Name: "foo"}}

	err := validateVolumeClaims(failingTest)
	if err == nil {
		t.Errorf("should have failed but passed for input: %+v", failingTest)
	} else {
		t.Logf("failed with error: %v", err)
	}

}
