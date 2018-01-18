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

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	api_v1 "k8s.io/kubernetes/pkg/api/v1"
	ext_v1beta1 "k8s.io/kubernetes/pkg/apis/extensions/v1beta1"

	"k8s.io/apimachinery/pkg/runtime"
)

func (app *App) validateDeployments() error {

	// TODO: v2
	return nil
}

func (app *App) fixDeployments() error {
	// auto populate name only if one deployment is specified without any name
	if len(app.Deployments) == 1 && app.Deployments[0].ObjectMeta.Name == "" {
		app.Deployments[0].ObjectMeta.Name = app.Name
	}

	for i := range app.Deployments {
		// copy root labels (already has "app: <app.Name>" label)
		for key, value := range app.ObjectMeta.Labels {
			app.Deployments[i].ObjectMeta.Labels = addKeyValueToMap(key, value, app.Deployments[i].ObjectMeta.Labels)
		}

		// copy root annotations (already has "appVersion: <app.AppVersion>" annotation)
		for key, value := range app.ObjectMeta.Annotations {
			app.Deployments[i].ObjectMeta.Annotations = addKeyValueToMap(key, value, app.Deployments[i].ObjectMeta.Annotations)
		}

		var err error
		app.Deployments[i].InitContainers, err = fixContainers(app.Deployments[i].InitContainers, app.Name)
		if err != nil {
			return errors.Wrap(err, "unable to fix init-containers")
		}

		app.Deployments[i].Containers, err = fixContainers(app.Deployments[i].Containers, app.Name)
		if err != nil {
			return errors.Wrap(err, "unable to fix containers")
		}

		vols, err := populateVolumes(app.Deployments[i].Containers, app.VolumeClaims, app.Deployments[i].Volumes)
		if err != nil {
			return errors.Wrapf(err, "unable to populate Volumes for deployment %q", app.Deployments[i].Name)
		}
		app.Deployments[i].Volumes = append(app.Deployments[i].Volumes, vols...)
	}
	return nil
}

// Creates a Deployment Kubernetes resource. The returned Deployment resource
// will be nil if it could not be generated due to insufficient input data.
func (app *App) createDeployments() ([]runtime.Object, error) {
	var deployments []runtime.Object

	for _, deployment := range app.Deployments {

		// We need to error out if both, deployment.PodSpec and deployment.DeploymentSpec are empty

		if deployment.isDeploymentSpecPodSpecEmpty() {
			log.Debug("Both, deployment.PodSpec and deployment.DeploymentSpec are empty, not enough data to create a deployment.")
			return nil, nil
		}

		// We are merging whole DeploymentSpec with PodSpec.
		// This means that someone could specify containers in template.spec and also in top level PodSpec.
		// This stupid check is supposed to make sure that only one of them set.
		// TODO: merge DeploymentSpec.Template.Spec and top level PodSpec
		if deployment.isMultiplePodSpecSpecified() {
			return nil, fmt.Errorf("Pod can't be specfied in two places. Use top level PodSpec or template.spec (DeploymentSpec.Template.Spec) not both")
		}

		deploymentSpec := deployment.DeploymentSpec

		// top level PodSpecMod is not empty, use it for deployment template
		// we already know that if deployment.PodSpec is not empty deployment.DeploymentSpec.Template.Spec is empty
		if !reflect.DeepEqual(deployment.PodSpecMod, PodSpecMod{}) {

			//copy over regular podSpec fields
			deploymentSpec.Template.Spec = deployment.PodSpec

			// our customized fields
			var err error
			deploymentSpec.Template.Spec.Containers, err = populateContainers(deployment.PodSpecMod.Containers, app.ConfigMaps, app.Secrets)
			if err != nil {
				return nil, errors.Wrapf(err, "deployment %q", app.Name)
			}
			log.Debugf("object after population: %#v\n", app)

			deploymentSpec.Template.Spec.InitContainers, err = populateContainers(deployment.PodSpecMod.InitContainers, app.ConfigMaps, app.Secrets)
			if err != nil {
				return nil, errors.Wrapf(err, "deployment %q", app.Name)
			}
			log.Debugf("object after population: %#v\n", app)
		}

		// TODO: check if this wasn't set by user, in that case we shouldn't overwrite it
		deploymentSpec.Template.ObjectMeta.Name = deployment.Name

		// TODO: merge with already existing labels and avoid duplication
		deploymentSpec.Template.ObjectMeta.Labels = deployment.Labels

		deploymentSpec.Template.ObjectMeta.Annotations = deployment.Annotations

		deployments = append(deployments, &ext_v1beta1.Deployment{
			ObjectMeta: deployment.ObjectMeta,
			Spec:       deploymentSpec,
		})
	}
	return deployments, nil
}

func (deployment *DeploymentSpecMod) isDeploymentSpecPodSpecEmpty() bool {
	return reflect.DeepEqual(deployment.PodSpecMod, PodSpecMod{}) && reflect.DeepEqual(deployment.DeploymentSpec, ext_v1beta1.DeploymentSpec{})
}

func (deployment *DeploymentSpecMod) isMultiplePodSpecSpecified() bool {
	return !(reflect.DeepEqual(deployment.DeploymentSpec.Template.Spec, api_v1.PodSpec{}) || reflect.DeepEqual(deployment.PodSpecMod, PodSpecMod{}))
}
