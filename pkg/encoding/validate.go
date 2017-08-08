package encoding

import (
	"fmt"

	"github.com/kedgeproject/kedge/pkg/spec"

	"github.com/pkg/errors"
)

func validateVolumeClaims(vcs []spec.VolumeClaim) error {
	// find the duplicate volume claim names, if found any then error out
	vcmap := make(map[string]interface{})
	for _, vc := range vcs {
		if _, ok := vcmap[vc.Name]; !ok {
			// value here does not matter
			vcmap[vc.Name] = nil
		} else {
			return fmt.Errorf("duplicate entry of volume claim %q", vc.Name)
		}
	}
	return nil
}

func validateApp(app *spec.App) error {

	// validate volumeclaims
	if err := validateVolumeClaims(app.VolumeClaims); err != nil {
		return errors.Wrap(err, "error validating volume claims")
	}

	return nil
}
