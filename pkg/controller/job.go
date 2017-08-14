package controller

import (
	"github.com/kedgeproject/kedge/pkg/transform/kubernetes"

	log "github.com/Sirupsen/logrus"
	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

func (job *jobSpecMod) Unmarshal(data []byte) error {
	err := yaml.Unmarshal(data, &job)
	if err != nil {
		return errors.Wrap(err, "could not unmarshal into internal struct")
	}
	log.Debugf("object unmarshalled: %#v\n", job)
	return nil
}

func (job *jobSpecMod) Transform() ([]runtime.Object, []string, error) {
	kJob := kubernetes.JobSpecMod(*job)
	return kJob.Transform()
}

func (job *jobSpecMod) Validate() error {
	return nil
}

func (job *jobSpecMod) Fix() error {
	return nil
}
