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
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	log "github.com/sirupsen/logrus"
	api_v1 "k8s.io/client-go/pkg/api/v1"

	"github.com/pkg/errors"
)

func populateServicePortNames(serviceName string, servicePorts []api_v1.ServicePort) {
	// auto populate port names if more than 1 port specified
	if len(servicePorts) > 1 {
		for i := range servicePorts {
			// Only populate if the port name is not already specified
			if len(servicePorts[i].Name) == 0 {
				servicePorts[i].Name = serviceName + "-" + strconv.FormatInt(int64(servicePorts[i].Port), 10)
			}
		}
	}
}

func populateProbes(c Container) (Container, error) {
	// check if health and liveness given together
	if c.Health != nil && (c.ReadinessProbe != nil || c.LivenessProbe != nil) {
		return c, fmt.Errorf("cannot define field 'health' and " +
			"'livenessProbe' or 'readinessProbe' together")
	}
	if c.Health != nil {
		c.LivenessProbe = c.Health
		c.ReadinessProbe = c.Health
		c.Health = nil
	}
	return c, nil
}

func searchConfigMap(cms []ConfigMapMod, name string) (ConfigMapMod, error) {
	for _, cm := range cms {
		if cm.Name == name {
			return cm, nil
		}
	}
	return ConfigMapMod{}, fmt.Errorf("configMap %q not found", name)
}

func getSecretDataKeys(secrets []SecretMod, name string) ([]string, error) {
	var dataKeys []string
	for _, secret := range secrets {
		if secret.Name == name {
			for dk := range secret.Data {
				dataKeys = append(dataKeys, dk)
			}
			for sdk := range secret.StringData {
				dataKeys = append(dataKeys, sdk)
			}
			return dataKeys, nil
		}
	}
	return nil, fmt.Errorf("secret %q not found", name)
}

func getMapKeys(m map[string]string) []string {
	var d []string
	for k := range m {
		d = append(d, k)
	}
	return d
}

func convertEnvFromToEnvs(envFrom []api_v1.EnvFromSource, cms []ConfigMapMod, secrets []SecretMod) ([]api_v1.EnvVar, error) {
	var envs []api_v1.EnvVar

	// we will iterate on all envFroms
	for ei, e := range envFrom {
		if e.ConfigMapRef != nil {

			cmName := e.ConfigMapRef.Name

			// see if the configMap name which is given actually exists
			cm, err := searchConfigMap(cms, cmName)
			if err != nil {
				return nil, errors.Wrapf(err, "envFrom[%d].configMapRef.name", ei)
			}
			// once that configMap is found extract all data from it and create a env out of it
			configMapKeys := getMapKeys(cm.Data)
			sort.Strings(configMapKeys)
			for _, key := range configMapKeys {
				envs = append(envs, api_v1.EnvVar{
					Name: key,
					ValueFrom: &api_v1.EnvVarSource{
						ConfigMapKeyRef: &api_v1.ConfigMapKeySelector{
							LocalObjectReference: api_v1.LocalObjectReference{
								Name: cmName,
							},
							Key: key,
						},
					},
				})
			}
		}

		if e.SecretRef != nil {
			rootSecretDataKeys, err := getSecretDataKeys(secrets, e.SecretRef.Name)
			if err != nil {
				return nil, errors.Wrapf(err, "envFrom[%d].secretRef.name", ei)
			}

			sort.Strings(rootSecretDataKeys)
			for _, secretDataKey := range rootSecretDataKeys {
				envs = append(envs, api_v1.EnvVar{
					Name: secretDataKey,
					ValueFrom: &api_v1.EnvVarSource{
						SecretKeyRef: &api_v1.SecretKeySelector{
							LocalObjectReference: api_v1.LocalObjectReference{
								Name: e.SecretRef.Name,
							},
							Key: secretDataKey,
						},
					},
				})
			}
		}
	}
	return envs, nil
}

func populateEnvFrom(c Container, cms []ConfigMapMod, secrets []SecretMod) (Container, error) {
	// now do the env from
	envs, err := convertEnvFromToEnvs(c.EnvFrom, cms, secrets)
	if err != nil {
		return c, err
	}
	// Since we are not supporting envFrom in our generated Kubernetes
	// artifacts right now, we need to set it as nil for every container.
	// This makes sure that Kubernetes artifacts do not contain an
	// envFrom field.
	// This is safe to set since all of the data from envFrom has been
	// extracted till this point.
	c.EnvFrom = nil
	// we collect all the envs from configMap before
	// envs provided inside the container
	envs = append(envs, c.Env...)
	c.Env = envs
	return c, nil
}

func populateContainers(containers []Container, cms []ConfigMapMod, secrets []SecretMod) ([]api_v1.Container, error) {
	var cnts []api_v1.Container

	for cn, c := range containers {
		// process health field
		c, err := populateProbes(c)
		if err != nil {
			return cnts, errors.Wrapf(err, "error converting 'health' to 'probes', app.containers[%d]", cn)
		}

		// process envFrom field
		c, err = populateEnvFrom(c, cms, secrets)
		if err != nil {
			return cnts, fmt.Errorf("error converting 'envFrom' to 'envs', app.containers[%d].%s", cn, err.Error())
		}

		// this is where we are only taking apart upstream container
		// and not our own remix of containers
		cnts = append(cnts, c.Container)
	}

	b, _ := json.MarshalIndent(cnts, "", "  ")
	log.Debugf("containers after populating health: %s", string(b))
	return cnts, nil
}

// Since we are automatically creating pvc from
// root level persistent volume and entry in the container
// volume mount, we also need to update the pod's volume field
func populateVolumes(containers []api_v1.Container, volumeClaims []VolumeClaim,
	volumes []api_v1.Volume) ([]api_v1.Volume, error) {
	var newPodVols []api_v1.Volume

	for cn, c := range containers {
		for vn, vm := range c.VolumeMounts {
			if isPVCDefined(volumeClaims, vm.Name) && !isVolumeDefined(volumes, vm.Name) {
				newPodVols = append(newPodVols, api_v1.Volume{
					Name: vm.Name,
					VolumeSource: api_v1.VolumeSource{
						PersistentVolumeClaim: &api_v1.PersistentVolumeClaimVolumeSource{
							ClaimName: vm.Name,
						},
					},
				})
			} else if !isVolumeDefined(volumes, vm.Name) {
				// pvc is not defined so we need to check if the entry is made in the pod volumes
				// since a volumeMount entry without entry in pod level volumes might cause failure
				// while deployment since that would not be a complete configuration
				return nil, fmt.Errorf("neither root level Persistent Volume"+
					" nor Volume in pod spec defined for %q, "+
					"in app.containers[%d].volumeMounts[%d]", vm.Name, cn, vn)
			}
		}
	}
	return newPodVols, nil
}
