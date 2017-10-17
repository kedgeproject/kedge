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
	log "github.com/Sirupsen/logrus"
	"github.com/davecgh/go-spew/spew"
	"github.com/ghodss/yaml"
	os_deploy_v1 "github.com/openshift/origin/pkg/deploy/apis/apps/v1"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	kapi "k8s.io/kubernetes/pkg/api/v1"
)

// Unmarshal the Kedge YAML file
func (deploymentConfig *DeploymentConfigSpecMod) Unmarshal(data []byte) error {
	err := yaml.Unmarshal(data, &deploymentConfig)
	if err != nil {
		return errors.Wrap(err, "could not unmarshal into internal struct")
	}
	log.Debugf("object unmarshalled: %#v\n", deploymentConfig)
	return nil
}

// Validate all portions of the file
func (deploymentConfig *DeploymentConfigSpecMod) Validate() error {

	if err := validateVolumeClaims(deploymentConfig.VolumeClaims); err != nil {
		return errors.Wrap(err, "error validating volume claims")
	}

	return nil
}

// Fix all services / volume claims / configmaps that are applied
// TODO: abstract out this code when more controllers are added
func (deploymentConfig *DeploymentConfigSpecMod) Fix() error {

	if err := deploymentConfig.ControllerFields.fixControllerFields(); err != nil {
		return errors.Wrap(err, "unable to fix ControllerFields")
	}

	// Fix DeploymentConfig
	deploymentConfig.fixDeploymentConfig()

	return nil
}

func (deploymentConfig *DeploymentConfigSpecMod) fixDeploymentConfig() {
	deploymentConfig.ControllerFields.ObjectMeta.Labels = addKeyValueToMap(appLabelKey,
		deploymentConfig.ControllerFields.Name,
		deploymentConfig.ControllerFields.ObjectMeta.Labels)

	// If the replicas are not specified at all, we need to set the value as 1
	if deploymentConfig.Replicas == nil {
		deploymentConfig.Replicas = getInt32Addr(1)
	}

	// Since we have unmarshalled replicas in a custom defined field, we need
	// to substitute the unmarshalled (and fixed) value in the internal
	// DeploymentConfigSpec struct
	deploymentConfig.DeploymentConfigSpec.Replicas = *deploymentConfig.Replicas
}

func (deploymentConfig *DeploymentConfigSpecMod) Transform() ([]runtime.Object, []string, error) {

	// Create Kubernetes objects (since OpenShift uses Kubernetes underneath, no need to refactor
	// this portion
	runtimeObjects, extraResources, err := deploymentConfig.CreateK8sObjects()
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create Kubernetes objects")
	}

	// Create the DeploymentConfig controller
	deploy, err := deploymentConfig.createOpenShiftController()
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create DeploymentConfig controller")
	}

	// adding controller objects
	// deploymentConfig will be nil if no deploymentConfig is generated and no error occurs,
	// so we only need to append this when a legit deploymentConfig resource is returned
	if deploy != nil {
		runtimeObjects = append(runtimeObjects, deploy)
		log.Debugf("deploymentConfig: %s, deploymentConfig: %s\n", deploy.Name, spew.Sprint(deploymentConfig))
	}

	if len(runtimeObjects) == 0 {
		return nil, nil, errors.New("No runtime objects created, possibly because not enough input data was passed")
	}

	scheme, err := GetScheme()
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to get scheme")
	}

	// Set's the appropriate GVK
	for _, runtimeObject := range runtimeObjects {
		if err := SetGVK(runtimeObject, scheme); err != nil {
			return nil, nil, errors.Wrap(err, "unable to set Group, Version and Kind for generated resources")
		}
	}

	return runtimeObjects, extraResources, nil
}

// Created the OpenShift DeploymentConfig controller
func (deploymentConfig *DeploymentConfigSpecMod) createOpenShiftController() (*os_deploy_v1.DeploymentConfig, error) {

	dcSpec := deploymentConfig.DeploymentConfigSpec
	dcSpec.Template = &kapi.PodTemplateSpec{
		Spec:       deploymentConfig.PodSpec,
		ObjectMeta: deploymentConfig.ControllerFields.ObjectMeta,
	}

	return &os_deploy_v1.DeploymentConfig{
		ObjectMeta: deploymentConfig.ObjectMeta,
		Spec:       dcSpec,
	}, nil
}
