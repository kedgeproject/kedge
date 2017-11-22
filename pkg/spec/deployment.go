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
	api_v1 "k8s.io/kubernetes/pkg/api/v1"
	ext_v1beta1 "k8s.io/kubernetes/pkg/apis/extensions/v1beta1"

	"github.com/davecgh/go-spew/spew"
	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

func (deployment *DeploymentSpecMod) Unmarshal(data []byte) error {
	err := yaml.Unmarshal(data, &deployment)
	if err != nil {
		return errors.Wrap(err, "Deployment could not unmarshal into internal struct")
	}
	log.Debugf("object unmarshalled: %#v\n", deployment)
	return nil
}

func (deployment *DeploymentSpecMod) Validate() error {

	// validate controller fields
	if err := deployment.ControllerFields.validateControllerFields(); err != nil {
		return errors.Wrap(err, "unable to validate controller fields")
	}

	return nil
}

func (deployment *DeploymentSpecMod) Fix() error {
	if err := deployment.ControllerFields.fixControllerFields(); err != nil {
		return errors.Wrap(err, "unable to fix ControllerFields")
	}

	deployment.ControllerFields.ObjectMeta.Labels = addKeyValueToMap(appLabelKey, deployment.ControllerFields.Name, deployment.ControllerFields.ObjectMeta.Labels)

	return nil
}

// Transform function if given DeploymentSpecMod data creates the versioned
// kubernetes objects and returns them in list of runtime.Object
// And if the field in DeploymentSpecMod called 'includeResources' is used
// then it returns the filenames mentioned there as list of string
func (deployment *DeploymentSpecMod) Transform() ([]runtime.Object, []string, error) {

	runtimeObjects, includeResources, err := deployment.CreateK8sObjects()
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create Kubernetes objects")
	}

	deploy, err := deployment.createKubernetesController()
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create Kubernetes Deployment controller")
	}

	// adding controller objects
	// deployment will be nil if no deployment is generated and no error occurs,
	// so we only need to append this when a legit deployment resource is returned
	if deploy != nil {
		runtimeObjects = append(runtimeObjects, deploy)
		log.Debugf("deployment: %s, deployment: %s\n", deploy.Name, spew.Sprint(deployment))

	}

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

// Creates a Deployment Kubernetes resource. The returned Deployment resource
// will be nil if it could not be generated due to insufficient input data.
func (deployment *DeploymentSpecMod) createKubernetesController() (*ext_v1beta1.Deployment, error) {

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

	// top level PodSpec is not empty, use it for deployment template
	// we already know that if deployment.PodSpec is not empty deployment.DeploymentSpec.Template.Spec is empty
	if !reflect.DeepEqual(deployment.PodSpec, api_v1.PodSpec{}) {
		deploymentSpec.Template.Spec = deployment.PodSpec
	}

	// TODO: check if this wasn't set by user, in that case we shouldn't overwrite it
	deploymentSpec.Template.ObjectMeta.Name = deployment.Name

	// TODO: merge with already existing labels and avoid duplication
	deploymentSpec.Template.ObjectMeta.Labels = deployment.Labels

	return &ext_v1beta1.Deployment{
		ObjectMeta: deployment.ObjectMeta,
		Spec:       deploymentSpec,
	}, nil
}

func (deployment *DeploymentSpecMod) isDeploymentSpecPodSpecEmpty() bool {
	return reflect.DeepEqual(deployment.PodSpec, api_v1.PodSpec{}) && reflect.DeepEqual(deployment.DeploymentSpec, ext_v1beta1.DeploymentSpec{})
}

func (deployment *DeploymentSpecMod) isMultiplePodSpecSpecified() bool {
	return !(reflect.DeepEqual(deployment.DeploymentSpec.Template.Spec, api_v1.PodSpec{}) || reflect.DeepEqual(deployment.PodSpec, api_v1.PodSpec{}))
}
