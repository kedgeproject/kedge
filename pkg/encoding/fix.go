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

package encoding

import (
	"fmt"
	"strconv"

	"github.com/kedgeproject/kedge/pkg/spec"

	"github.com/pkg/errors"
)

func fixApp(app *spec.App) error {

	// fix app.Services
	if err := fixServices(app); err != nil {
		return errors.Wrap(err, "Unable to fix services")
	}

	// fix app.VolumeClaims
	if err := fixVolumeClaims(app); err != nil {
		return errors.Wrap(err, "Unable to fix persistentVolume")
	}

	// fix app.configMaps
	if err := fixConfigMaps(app); err != nil {
		return errors.Wrap(err, "unable to fix configMaps")
	}

	if err := fixContainers(app); err != nil {
		return errors.Wrap(err, "unable to fix containers")
	}

	if err := fixSecrets(app); err != nil {
		return errors.Wrap(err, "unable to fix secrets")
	}

	return nil
}

func fixServices(app *spec.App) error {
	for i, service := range app.Services {
		// auto populate service name if only one service is specified
		if service.Name == "" {
			if len(app.Services) == 1 {
				service.Name = app.Name
			} else {
				return errors.New("More than one service mentioned, please specify name for each one")
			}
		}
		app.Services[i] = service

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

func fixVolumeClaims(app *spec.App) error {
	for i, pVolume := range app.VolumeClaims {
		if pVolume.Name == "" {
			if len(app.VolumeClaims) == 1 {
				pVolume.Name = app.Name
			} else {
				return errors.New("More than one persistent volume mentioned, please specify name for each one")
			}
		}
		app.VolumeClaims[i] = pVolume
	}
	return nil
}

func fixConfigMaps(app *spec.App) error {
	// if only one configMap is defined and its name is not specified
	if len(app.ConfigMaps) == 1 && app.ConfigMaps[0].Name == "" {
		app.ConfigMaps[0].Name = app.Name
	} else if len(app.ConfigMaps) > 1 {
		// if multiple configMaps is defined then each should have a name
		for cdn, cd := range app.ConfigMaps {
			if cd.Name == "" {
				return fmt.Errorf("name not specified for app.configMaps[%d]", cdn)
			}
		}
	}
	return nil
}

func fixSecrets(app *spec.App) error {
	// populate secret name only if one secret is specified
	if len(app.Secrets) == 1 && app.Secrets[0].Name == "" {
		app.Secrets[0].Name = app.Name
	} else if len(app.Secrets) > 1 {
		for i, sec := range app.Secrets {
			if sec.Name == "" {
				return fmt.Errorf("name not specified for app.secrets[%d]", i)
			}
		}
	}
	return nil
}

func fixContainers(app *spec.App) error {
	// if only one container set name of it as app name
	if len(app.Containers) == 1 && app.Containers[0].Name == "" {
		app.Containers[0].Name = app.Name
	} else if len(app.Containers) > 1 {
		// check if all the containers have a name
		// if not fail giving error
		for cn, c := range app.Containers {
			if c.Name == "" {
				return fmt.Errorf("app %q: container name not defined for app.containers[%d]", app.Name, cn)
			}
		}
	}
	return nil
}
