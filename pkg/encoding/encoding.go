package encoding

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ghodss/yaml"
	"github.com/pkg/errors"

	"github.com/kedgeproject/kedge/pkg/spec"
)

func Decode(data []byte) (*spec.App, error) {

	var app spec.App
	err := yaml.Unmarshal(data, &app)
	if err != nil {
		return nil, errors.Wrap(err, "could not unmarshal into internal struct")
	}
	log.Debugf("object unmrashalled: %#v\n", app)
	if err := fixApp(&app); err != nil {
		return nil, errors.Wrapf(err, "Unable to fix app %q", app.Name)
	}
	return &app, nil

}
