package encoding

import (
	"github.com/ghodss/yaml"
	"github.com/pkg/errors"

	log "github.com/Sirupsen/logrus"

	"github.com/surajssd/kapp/pkg/spec"
)

func Decode(data []byte) (*spec.App, error) {

	var app spec.App
	err := yaml.Unmarshal(data, &app)
	if err != nil {
		return nil, errors.Wrap(err, "could not unmarshal into internal struct")
	}
	log.Debugf("object unmrashalled: %#v\n", app)
	return &app, nil

}
