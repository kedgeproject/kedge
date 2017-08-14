package controller

import (
	"fmt"
	"strconv"

	"github.com/kedgeproject/kedge/pkg/spec"
	"github.com/kedgeproject/kedge/pkg/transform/kubernetes"

	log "github.com/Sirupsen/logrus"
	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

func (deployment *deploymentSpecMod) Unmarshal(data []byte) error {
	err := yaml.Unmarshal(data, &deployment)
	if err != nil {
		return errors.Wrap(err, "could not unmarshal into internal struct")
	}
	log.Debugf("object unmarshalled: %#v\n", deployment)
	return nil
}

func (deployment *deploymentSpecMod) Validate() error {
	// validate volumeclaims
	if err := validateVolumeClaims(deployment.VolumeClaims); err != nil {
		return errors.Wrap(err, "error validating volume claims")
	}

	return nil
}

func (deployment *deploymentSpecMod) Fix() error {

	if err := fixServices(deployment); err != nil {
		return errors.Wrap(err, "Unable to fix services")
	}

	if err := fixVolumeClaims(deployment); err != nil {
		return errors.Wrap(err, "Unable to fix persistentVolume")
	}

	if err := fixConfigMaps(deployment); err != nil {
		return errors.Wrap(err, "unable to fix configMaps")
	}

	if err := fixContainers(deployment); err != nil {
		return errors.Wrap(err, "unable to fix containers")
	}

	if err := fixSecrets(deployment); err != nil {
		return errors.Wrap(err, "unable to fix secrets")
	}

	return nil
}

func (deployment *deploymentSpecMod) Transform() ([]runtime.Object, []string, error) {
	kDeployment := kubernetes.DeploymentSpecMod(*deployment)
	return kDeployment.Transform()
}

func validateVolumeClaims(vcs []spec.VolumeClaim) error {
	// find the duplicate volume claim names, if found any then error out
	vcmap := make(map[string]interface{})
	for _, vc := range vcs {
		if _, ok := vcmap[vc.Name]; !ok {
			// value here does not matter
			vcmap[vc.Name] = nil
		} else {
			return fmt.Errorf("duplicate entry of volume claim %q", vc.Name)
		}
	}
	return nil
}

func fixServices(deployment *deploymentSpecMod) error {
	for i, service := range deployment.Services {
		// auto populate service name if only one service is specified
		if service.Name == "" {
			if len(deployment.Services) == 1 {
				service.Name = deployment.Name
			} else {
				return errors.New("More than one service mentioned, please specify name for each one")
			}
		}
		deployment.Services[i] = service

		for i, servicePort := range service.Ports {
			// auto populate port names if not specified
			if len(service.Ports) > 1 && servicePort.Name == "" {
				servicePort.Name = service.Name + "-" + strconv.FormatInt(int64(servicePort.Port), 10)
			}
			service.Ports[i] = servicePort
		}
	}
	return nil
}

func fixVolumeClaims(deployment *deploymentSpecMod) error {
	for i, pVolume := range deployment.VolumeClaims {
		if pVolume.Name == "" {
			if len(deployment.VolumeClaims) == 1 {
				pVolume.Name = deployment.Name
			} else {
				return errors.New("More than one persistent volume mentioned, please specify name for each one")
			}
		}
		deployment.VolumeClaims[i] = pVolume
	}
	return nil
}

func fixConfigMaps(deployment *deploymentSpecMod) error {
	// if only one configMap is defined and its name is not specified
	if len(deployment.ConfigMaps) == 1 && deployment.ConfigMaps[0].Name == "" {
		deployment.ConfigMaps[0].Name = deployment.Name
	} else if len(deployment.ConfigMaps) > 1 {
		// if multiple configMaps is defined then each should have a name
		for cdn, cd := range deployment.ConfigMaps {
			if cd.Name == "" {
				return fmt.Errorf("name not specified for deployment.configMaps[%d]", cdn)
			}
		}
	}
	return nil
}

func fixSecrets(deployment *deploymentSpecMod) error {
	// populate secret name only if one secret is specified
	if len(deployment.Secrets) == 1 && deployment.Secrets[0].Name == "" {
		deployment.Secrets[0].Name = deployment.Name
	} else if len(deployment.Secrets) > 1 {
		for i, sec := range deployment.Secrets {
			if sec.Name == "" {
				return fmt.Errorf("name not specified for deployment.secrets[%d]", i)
			}
		}
	}
	return nil
}

func fixContainers(deployment *deploymentSpecMod) error {
	// if only one container set name of it as deployment name
	if len(deployment.Containers) == 1 && deployment.Containers[0].Name == "" {
		deployment.Containers[0].Name = deployment.Name
	} else if len(deployment.Containers) > 1 {
		// check if all the containers have a name
		// if not fail giving error
		for cn, c := range deployment.Containers {
			if c.Name == "" {
				return fmt.Errorf("deployment %q: container name not defined for deployment.containers[%d]", deployment.Name, cn)
			}
		}
	}
	return nil
}
