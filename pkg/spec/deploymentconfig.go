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
	os_v1 "github.com/kedgeproject/kedge/pkg/hack/openshift"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// Unmarshal the Kedge YAML file
func (deployment *DeploymentConfigSpecMod) Unmarshal(data []byte) error {
	err := yaml.Unmarshal(data, &deployment)
	if err != nil {
		return errors.Wrap(err, "could not unmarshal into internal struct")
	}
	log.Debugf("object unmarshalled: %#v\n", deployment)
	return nil
}

// Validate all portions of the file
func (deployment *DeploymentConfigSpecMod) Validate() error {

	if err := validateVolumeClaims(deployment.VolumeClaims); err != nil {
		return errors.Wrap(err, "error validating volume claims")
	}

	return nil
}

// Fix all services / volume claims / configmaps that are applied
// TODO: abstract out this code when more controllers are added
func (deployment *DeploymentConfigSpecMod) Fix() error {

	var err error

	// fix deployment.Services
	deployment.Services, err = fixServices(deployment.Services, deployment.Name)
	if err != nil {
		return errors.Wrap(err, "Unable to fix services")
	}

	// fix deployment.VolumeClaims
	deployment.VolumeClaims, err = fixVolumeClaims(deployment.VolumeClaims, deployment.Name)
	if err != nil {
		return errors.Wrap(err, "Unable to fix persistentVolume")
	}

	// fix deployment.configMaps
	deployment.ConfigMaps, err = fixConfigMaps(deployment.ConfigMaps, deployment.Name)
	if err != nil {
		return errors.Wrap(err, "unable to fix configMaps")
	}

	deployment.Containers, err = fixContainers(deployment.Containers, deployment.Name)
	if err != nil {
		return errors.Wrap(err, "unable to fix containers")
	}

	deployment.InitContainers, err = fixContainers(deployment.InitContainers, deployment.Name)
	if err != nil {
		return errors.Wrap(err, "unable to fix init-containers")
	}

	deployment.Secrets, err = fixSecrets(deployment.Secrets, deployment.Name)
	if err != nil {
		return errors.Wrap(err, "unable to fix secrets")
	}

	return nil
}

func (deployment *DeploymentConfigSpecMod) Transform() ([]runtime.Object, []string, error) {

	// Create Kubernetes objects (since OpenShift uses Kubernetes underneath, no need to refactor
	// this portion
	runtimeObjects, extraResources, err := deployment.CreateK8sObjects()
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create Kubernetes objects")
	}

	// Set appropriate GVK BEFORE adding DeploymentConfig controller
	// as OpenShift controllers are not available in the Kubernetes controller / setGVK check
	scheme, err := GetScheme()
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to get scheme")
	}

	// Set's the appropriate GVK
	for _, runtimeObject := range runtimeObjects {
		if err := SetGVK(runtimeObject, scheme); err != nil {
			return nil, nil, errors.Wrap(err, "unable to set Group, Version and Kind for generated Kubernetes resources")
		}
	}

	// Create the DeploymentConfig controller!
	deploy, err := deployment.createDeploymentConfigController()
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

	return runtimeObjects, extraResources, nil
}

// TODO: Rather than return blank, return fully usable DeploymentConfig
func (deployment *DeploymentConfigSpecMod) createDeploymentConfigController() (*os_v1.DeploymentConfig, error) {
	return &os_v1.DeploymentConfig{
		TypeMeta: metav1.TypeMeta{
			Kind:       "DeploymentConfig",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   deployment.Name,
			Labels: deployment.Labels,
		},
		Spec: os_v1.DeploymentConfigSpec{},
	}, nil
}
