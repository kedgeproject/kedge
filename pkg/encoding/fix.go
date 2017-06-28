package encoding

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	"github.com/kedgeproject/kedge/pkg/spec"
)

func fixApp(app *spec.App) error {

	// fix app.Services
	if err := fixServices(app); err != nil {
		return errors.Wrap(err, "Unable to fix services")
	}

	// fix app.PersistentVolumes
	if err := fixPersistentVolumes(app); err != nil {
		return errors.Wrap(err, "Unable to fix persistentVolume")
	}

	// fix app.configMaps
	if err := fixConfigMaps(app); err != nil {
		return errors.Wrap(err, "unable to fix configMaps")
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

func fixPersistentVolumes(app *spec.App) error {
	for i, pVolume := range app.PersistentVolumes {
		if pVolume.Name == "" {
			if len(app.PersistentVolumes) == 1 {
				pVolume.Name = app.Name
			} else {
				return errors.New("More than one persistent volume mentioned, please specify name for each one")
			}
		}
		app.PersistentVolumes[i] = pVolume
	}
	return nil
}

func fixConfigMaps(app *spec.App) error {
	// if only one configMap is defined and it's name is not specified
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
