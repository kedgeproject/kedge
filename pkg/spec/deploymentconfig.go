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
	"k8s.io/apimachinery/pkg/runtime"

	os_deploy_v1 "github.com/openshift/origin/pkg/apps/apis/apps/v1"
	api_v1 "k8s.io/kubernetes/pkg/api/v1"
)

func (deploymentConfig *DeploymentConfigSpecMod) validateDeploymentConfig() error {
	// TODO: v2
	return nil
}

func (app *App) fixDeploymentConfigs() error {
	// auto populate name only if one deployment is specified without any name
	if len(app.DeploymentConfigs) == 1 && app.DeploymentConfigs[0].ObjectMeta.Name == "" {
		app.DeploymentConfigs[0].ObjectMeta.Name = app.Name
	}

	for i := range app.DeploymentConfigs {
		// If the replicas are not specified at all, we need to set the value as 1
		if app.DeploymentConfigs[i].Replicas == nil {
			app.DeploymentConfigs[i].Replicas = getInt32Addr(1)
		}

		// Since we have unmarshalled replicas in a custom defined field, we need
		// to substitute the unmarshalled (and fixed) value in the internal
		// DeploymentConfigSpec struct
		app.DeploymentConfigs[i].DeploymentConfigSpec.Replicas = *app.DeploymentConfigs[i].Replicas

		// copy root labels (already has "app: <app.Name>" label)
		for key, value := range app.ObjectMeta.Labels {
			app.DeploymentConfigs[i].ObjectMeta.Labels = addKeyValueToMap(key, value, app.DeploymentConfigs[i].ObjectMeta.Labels)
		}

		// copy root annotations (already has "appVersion: <app.AppVersion>" annotation)
		for key, value := range app.ObjectMeta.Annotations {
			app.DeploymentConfigs[i].ObjectMeta.Annotations = addKeyValueToMap(key, value, app.DeploymentConfigs[i].ObjectMeta.Annotations)
		}

		var err error
		app.DeploymentConfigs[i].InitContainers, err = fixContainers(app.DeploymentConfigs[i].InitContainers, app.Name)
		if err != nil {
			return errors.Wrap(err, "unable to fix init-containers")
		}

		app.DeploymentConfigs[i].Containers, err = fixContainers(app.DeploymentConfigs[i].Containers, app.Name)
		if err != nil {
			return errors.Wrap(err, "unable to fix containers")
		}

		vols, err := populateVolumes(app.DeploymentConfigs[i].Containers, app.VolumeClaims, app.DeploymentConfigs[i].Volumes)
		if err != nil {
			return errors.Wrapf(err, "unable to populate Volumes for deploymentConfig %q", app.DeploymentConfigs[i].Name)
		}
		app.DeploymentConfigs[i].Volumes = append(app.DeploymentConfigs[i].Volumes, vols...)
	}

	return nil
}

// Created the OpenShift DeploymentConfig controller
func (app *App) createDeploymentConfigs() ([]runtime.Object, error) {
	var deploymentConfigs []runtime.Object

	for _, deploymentConfig := range app.DeploymentConfigs {

		// We need to error out if both, deployment.PodSpec and deployment.DeploymentSpec are empty

		if deploymentConfig.isDeploymentConfigSpecPodSpecEmpty() {
			log.Debug("Both, deploymentConfig.PodSpec and deploymentConfig.DeploymentConfigSpec are empty, not enough data to create a deployment.")
			return nil, nil
		}

		// We are merging whole DeploymentSpec with PodSpec.
		// This means that someone could specify containers in template.spec and also in top level PodSpec.
		// This stupid check is supposed to make sure that only one of them set.
		// TODO: merge DeploymentSpec.Template.Spec and top level PodSpec
		if deploymentConfig.isMultiplePodSpecSpecified() {
			return nil, fmt.Errorf("Pod can't be specfied in two places. Use top level PodSpec or template.spec (DeploymentSpec.Template.Spec) not both")
		}

		deploymentConfigSpec := deploymentConfig.DeploymentConfigSpec

		// top level PodSpecMod is not empty, use it for deployment template
		// we already know that if deployment.PodSpec is not empty deployment.DeploymentSpec.Template.Spec is empty
		if !reflect.DeepEqual(deploymentConfig.PodSpecMod, PodSpecMod{}) {

			if deploymentConfigSpec.Template == nil {
				deploymentConfigSpec.Template = &api_v1.PodTemplateSpec{}
			}

			//copy over regular podSpec fields
			deploymentConfigSpec.Template.Spec = deploymentConfig.PodSpec

			// our customized fields
			var err error
			deploymentConfigSpec.Template.Spec.Containers, err = populateContainers(deploymentConfig.PodSpecMod.Containers, app.ConfigMaps, app.Secrets)
			if err != nil {
				return nil, errors.Wrapf(err, "deployment %q", app.Name)
			}
			log.Debugf("object after population: %#v\n", app)

			deploymentConfigSpec.Template.Spec.InitContainers, err = populateContainers(deploymentConfig.PodSpecMod.InitContainers, app.ConfigMaps, app.Secrets)
			if err != nil {
				return nil, errors.Wrapf(err, "deployment %q", app.Name)
			}
			log.Debugf("object after population: %#v\n", app)
		}

		// TODO: check if this wasn't set by user, in that case we shouldn't overwrite it
		deploymentConfigSpec.Template.ObjectMeta.Name = deploymentConfig.Name

		// TODO: merge with already existing labels and avoid duplication
		deploymentConfigSpec.Template.ObjectMeta.Labels = deploymentConfig.Labels

		deploymentConfigSpec.Template.ObjectMeta.Annotations = deploymentConfig.Annotations

		deploymentConfigs = append(deploymentConfigs, &os_deploy_v1.DeploymentConfig{
			ObjectMeta: deploymentConfig.ObjectMeta,
			Spec:       deploymentConfigSpec,
		})
	}
	return deploymentConfigs, nil
}

func (deploymentConfig *DeploymentConfigSpecMod) isDeploymentConfigSpecPodSpecEmpty() bool {
	return reflect.DeepEqual(deploymentConfig.PodSpecMod, PodSpecMod{}) && reflect.DeepEqual(deploymentConfig.DeploymentConfigSpec, os_deploy_v1.DeploymentConfigSpec{})
}

func (deploymentConfig *DeploymentConfigSpecMod) isMultiplePodSpecSpecified() bool {
	// OpenShift DeploymentConfigSpec.Template is pointer to PodTemplateSpec
	// First we need to check if it is not nil
	if deploymentConfig.DeploymentConfigSpec.Template == nil {
		return false
	}
	return !(reflect.DeepEqual(deploymentConfig.DeploymentConfigSpec.Template.Spec, api_v1.PodSpec{}) || reflect.DeepEqual(deploymentConfig.PodSpecMod, PodSpecMod{}))
}
