package spec

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	os_api "github.com/openshift/origin/pkg/apps/apis/apps/v1"

	"github.com/davecgh/go-spew/spew"
	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/pkg/api"
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

	for _, runtimeObject := range runtimeObjects {

		gvk, isUnversioned, err := api.Scheme.ObjectKind(runtimeObject)
		if err != nil {
			return nil, nil, errors.Wrap(err, "ConvertToVersion failed")
		}
		if isUnversioned {
			return nil, nil, fmt.Errorf("ConvertToVersion failed: can't output unversioned type: %T", runtimeObject)
		}

		runtimeObject.GetObjectKind().SetGroupVersionKind(gvk)
	}

	return runtimeObjects, extraResources, nil
}

func (deployment *DeploymentConfigSpecMod) createDeploymentConfigController() (*os_api.DeploymentConfig, error) {

	return &os_api.DeploymentConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:   deployment.Name,
			Labels: deployment.Labels,
		},
		Spec: os_api.DeploymentConfigSpec{},
	}, nil

	/*
		return &ext_v1beta1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:   deployment.Name,
				Labels: deployment.Labels,
			},
			Spec: deploymentSpec,
		}, nil
	*/
}
