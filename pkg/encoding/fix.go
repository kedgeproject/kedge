package encoding

import (
	"github.com/pkg/errors"
	"github.com/surajssd/kapp/pkg/spec"
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

	return nil
}

func fixServices(app *spec.App) error {
	for i, service := range app.Services {
		if service.Name == "" {
			if len(app.Services) == 1 {
				service.Name = app.Name
			} else {
				return errors.New("More than one service mentioned, please specify name for each one")
			}
		}
		app.Services[i] = service
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
