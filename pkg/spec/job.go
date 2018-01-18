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

	"github.com/pkg/errors"

	log "github.com/Sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
	api_v1 "k8s.io/kubernetes/pkg/api/v1"
	batch_v1 "k8s.io/kubernetes/pkg/apis/batch/v1"
)

func (app *App) fixJobs() error {

	// auto populate name only if one deployment is specified without any name
	if len(app.Jobs) == 1 && app.Jobs[0].ObjectMeta.Name == "" {
		app.Jobs[0].ObjectMeta.Name = app.Name
	}

	for i := range app.Jobs {

		// copy root labels (already has "app: <app.Name>" label)
		for key, value := range app.ObjectMeta.Labels {
			app.Jobs[i].ObjectMeta.Labels = addKeyValueToMap(key, value, app.Jobs[i].ObjectMeta.Labels)
		}

		// copy root annotations (already has "appVersion: <app.AppVersion>" annotation)
		for key, value := range app.ObjectMeta.Annotations {
			app.Jobs[i].ObjectMeta.Annotations = addKeyValueToMap(key, value, app.Jobs[i].ObjectMeta.Annotations)
		}

		var err error
		app.Jobs[i].InitContainers, err = fixContainers(app.Jobs[i].InitContainers, app.Name)
		if err != nil {
			return errors.Wrap(err, "unable to fix init-containers")
		}

		app.Jobs[i].Containers, err = fixContainers(app.Jobs[i].Containers, app.Name)
		if err != nil {
			return errors.Wrap(err, "unable to fix containers")
		}

		vols, err := populateVolumes(app.Jobs[i].Containers, app.VolumeClaims, app.Jobs[i].Volumes)
		if err != nil {
			return errors.Wrapf(err, "unable to populate Volumes for Job %q", app.Jobs[i].Name)
		}
		app.Jobs[i].Volumes = append(app.Jobs[i].Volumes, vols...)

		// if RestartPolicy is not set by user default it to 'OnFailure'
		if app.Jobs[0].RestartPolicy == "" {
			app.Jobs[0].RestartPolicy = api_v1.RestartPolicyOnFailure
		}
	}
	return nil
}

func (app *App) createJobs() ([]runtime.Object, error) {

	var jobs []runtime.Object

	for _, job := range app.Jobs {

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

			// our customized fields
			var err error
			jobSpec.Template.Spec.Containers, err = populateContainers(job.PodSpecMod.Containers, app.ConfigMaps, app.Secrets)
			if err != nil {
				return nil, errors.Wrapf(err, "job %q", app.Name)
			}
			log.Debugf("object after population: %#v\n", app)

			jobSpec.Template.Spec.InitContainers, err = populateContainers(job.PodSpecMod.InitContainers, app.ConfigMaps, app.Secrets)
			if err != nil {
				return nil, errors.Wrapf(err, "deployment %q", app.Name)
			}
			log.Debugf("object after population: %#v\n", app)
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

		jobs = append(jobs, &batch_v1.Job{
			ObjectMeta: job.ObjectMeta,
			Spec:       jobSpec,
		})
	}
	return jobs, nil
}

func (app *App) validateJobs() error {
	for _, job := range app.Jobs {
		if job.RestartPolicy == api_v1.RestartPolicyAlways {
			return fmt.Errorf("the Job %q is invalid: restartPolicy: unsupported value: \"Always\": supported values: OnFailure, Never", job.Name)
		}
	}
	return nil
}

func (job *JobSpecMod) isJobSpecPodSpecEmpty() bool {
	return reflect.DeepEqual(job.PodSpec, api_v1.PodSpec{}) && reflect.DeepEqual(job.JobSpec, batch_v1.JobSpec{})
}

func (job *JobSpecMod) isMultiplePodSpecSpecified() bool {
	return !(reflect.DeepEqual(job.JobSpec.Template.Spec, api_v1.PodSpec{}) || reflect.DeepEqual(job.PodSpec, api_v1.PodSpec{}))
}
