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
	"strconv"

	"github.com/pkg/errors"
)

func FixApp(app *App) error {
	var err error

	// fix app.Services
	app.Services, err = fixServices(app.Services, app.Name)
	if err != nil {
		return errors.Wrap(err, "Unable to fix services")
	}

	// fix app.VolumeClaims
	app.VolumeClaims, err = fixVolumeClaims(app.VolumeClaims, app.Name)
	if err != nil {
		return errors.Wrap(err, "Unable to fix persistentVolume")
	}

	// fix app.configMaps
	app.ConfigMaps, err = fixConfigMaps(app.ConfigMaps, app.Name)
	if err != nil {
		return errors.Wrap(err, "unable to fix configMaps")
	}

	app.Containers, err = fixContainers(app.Containers, app.Name)
	if err != nil {
		return errors.Wrap(err, "unable to fix containers")
	}

	app.InitContainers, err = fixContainers(app.InitContainers, app.Name)
	if err != nil {
		return errors.Wrap(err, "unable to fix init-containers")
	}

	app.Secrets, err = fixSecrets(app.Secrets, app.Name)
	if err != nil {
		return errors.Wrap(err, "unable to fix secrets")
	}

	return nil
}

func fixServices(services []ServiceSpecMod, appName string) ([]ServiceSpecMod, error) {
	for i, service := range services {
		// auto populate service name if only one service is specified
		if service.Name == "" {
			if len(services) == 1 {
				service.Name = appName
			} else {
				return nil, errors.New("More than one service mentioned, please specify name for each one")
			}
		}

		for i, servicePort := range service.Ports {
			// auto populate port names if not specified
			if len(service.Ports) > 1 && servicePort.Name == "" {
				servicePort.Name = service.Name + "-" + strconv.FormatInt(int64(servicePort.Port), 10)
			}
			service.Ports[i] = servicePort
		}

		// this should be the last statement in this for loop
		services[i] = service
	}
	return services, nil
}

func fixVolumeClaims(volumeClaims []VolumeClaim, appName string) ([]VolumeClaim, error) {
	for i, pVolume := range volumeClaims {
		if pVolume.Name == "" {
			if len(volumeClaims) == 1 {
				pVolume.Name = appName
			} else {
				return nil, errors.New("More than one persistent volume mentioned," +
					" please specify name for each one")
			}
		}
		volumeClaims[i] = pVolume
	}
	return volumeClaims, nil
}

func fixConfigMaps(configMaps []ConfigMapMod, appName string) ([]ConfigMapMod, error) {
	// if only one configMap is defined and its name is not specified
	if len(configMaps) == 1 && configMaps[0].Name == "" {
		configMaps[0].Name = appName
	} else if len(configMaps) > 1 {
		// if multiple configMaps is defined then each should have a name
		for cdn, cd := range configMaps {
			if cd.Name == "" {
				return nil, fmt.Errorf("name not specified for app.configMaps[%d]", cdn)
			}
		}
	}
	return configMaps, nil
}

func fixSecrets(secrets []SecretMod, appName string) ([]SecretMod, error) {
	// populate secret name only if one secret is specified
	if len(secrets) == 1 && secrets[0].Name == "" {
		secrets[0].Name = appName
	} else if len(secrets) > 1 {
		for i, sec := range secrets {
			if sec.Name == "" {
				return nil, fmt.Errorf("name not specified for app.secrets[%d]", i)
			}
		}
	}
	return secrets, nil
}

func fixContainers(containers []Container, appName string) ([]Container, error) {
	// if only one container set name of it as app name
	if len(containers) == 1 && containers[0].Name == "" {
		containers[0].Name = appName
	} else if len(containers) > 1 {
		// check if all the containers have a name
		// if not fail giving error
		for cn, c := range containers {
			if c.Name == "" {
				return nil, fmt.Errorf("app %q: container name not defined for app.containers[%d]", appName, cn)
			}
		}
	}
	return containers, nil
}
