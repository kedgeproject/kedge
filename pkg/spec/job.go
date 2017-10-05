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
	"fmt"
	"reflect"

	"github.com/davecgh/go-spew/spew"
	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"

	log "github.com/Sirupsen/logrus"
	api_v1 "k8s.io/client-go/pkg/api/v1"
	batch_v1 "k8s.io/client-go/pkg/apis/batch/v1"
)

func (job *JobSpecMod) Unmarshal(data []byte) error {
	err := yaml.Unmarshal(data, &job)
	if err != nil {
		return errors.Wrap(err, "could not unmarshal into internal struct")
	}
	log.Debugf("object unmarshalled: %#v\n", job)
	return nil
}

func (job *JobSpecMod) Fix() error {
	if err := job.ControllerFields.fixControllerFields(); err != nil {
		return errors.Wrap(err, "unable to fix ControllerFields")
	}

	job.ControllerFields.ObjectMeta.Labels = addKeyValueToMap(appLabelKey, job.ControllerFields.Name, job.ControllerFields.ObjectMeta.Labels)

	return nil
}

func (job *JobSpecMod) Transform() ([]runtime.Object, []string, error) {

	runtimeObjects, includeResources, err := job.CreateK8sObjects()
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create Kubernetes objects")
	}

	j, err := job.CreateK8sController()
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create Kubernetes Job controller")
	}

	if j != nil {
		runtimeObjects = append(runtimeObjects, j)
		log.Debug("job: &s, job: &s\n", j.Name, spew.Sprint(j))
	}

	// TODO: abstract out following code, but holding back since would be great
	// to make scheme controller specific
	if len(runtimeObjects) == 0 {
		return nil, nil, errors.New("No runtime objects created, possibly because not enough input data was passed")
	}

	scheme, err := GetScheme()
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to get scheme")
	}

	for _, runtimeObject := range runtimeObjects {
		if err := SetGVK(runtimeObject, scheme); err != nil {
			return nil, nil, errors.Wrap(err, "unable to set Group, Version and Kind for generated Kubernetes resources")
		}
	}

	return runtimeObjects, includeResources, nil
}

func (job *JobSpecMod) CreateK8sController() (*batch_v1.Job, error) {

	// We need to error out if both, job.PodSpec and job.JobSpec are empty
	if job.isJobSpecPodSpecEmpty() {
		log.Debug("Both, job.PodSpec and job.JobSpec are empty, not enough data to create a job.")
		return nil, nil
	}

	// Checking if PodSpec is specified at multiple levels
	if job.isMultiplePodSpecSpecified() {
		return nil, fmt.Errorf("Pod can't be specfied in two places. Use top level PodSpec or template.spec (JobSpec.Template.Spec) not both")
	}

	jobSpec := job.JobSpec

	if !reflect.DeepEqual(job.PodSpec, api_v1.PodSpec{}) {
		jobSpec.Template.Spec = job.PodSpec
	}

	// activeDeadlineSeconds is a conflicting field which exists in both,
	// v1.PodSpec and batch/v1.JobSpec, and both of these fields exist at the
	// top level of JobSpecMod.
	// So, whenever activeDeadlineSeconds field is passed, we will only
	// populate the JobSpec and not the PodSpec.
	// To populate PodSpec's activeDeadlineSeconds, the user will have to pass
	// this field the long way by defining PodSpec exclusively.
	if job.ActiveDeadlineSeconds != nil {
		jobSpec.ActiveDeadlineSeconds = job.ActiveDeadlineSeconds
	}

	return &batch_v1.Job{
		ObjectMeta: job.ObjectMeta,
		Spec:       jobSpec,
	}, nil
}

func (job *JobSpecMod) Validate() error {

	// validate controller fields
	if err := job.ControllerFields.validateControllerFields(); err != nil {
		return errors.Wrap(err, "unable to validate controller fields")
	}

	return nil
}

func (job *JobSpecMod) isJobSpecPodSpecEmpty() bool {
	return reflect.DeepEqual(job.PodSpec, api_v1.PodSpec{}) && reflect.DeepEqual(job.JobSpec, batch_v1.JobSpec{})
}

func (job *JobSpecMod) isMultiplePodSpecSpecified() bool {
	return !(reflect.DeepEqual(job.JobSpec.Template.Spec, api_v1.PodSpec{}) || reflect.DeepEqual(job.PodSpec, api_v1.PodSpec{}))
}
