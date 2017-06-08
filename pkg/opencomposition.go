package pkg

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/resource"
	"k8s.io/client-go/pkg/runtime"
	"k8s.io/client-go/pkg/util/intstr"

	log "github.com/Sirupsen/logrus"
	api_v1 "k8s.io/client-go/pkg/api/v1"
	ext_v1beta1 "k8s.io/client-go/pkg/apis/extensions/v1beta1"

	// install api
	_ "k8s.io/client-go/pkg/api/install"
	_ "k8s.io/client-go/pkg/apis/extensions/install"
)

type Volume struct {
	api_v1.Volume `yaml:",inline"`
	Size          string
}

type App struct {
	Name              string            `yaml:"name"`
	Replicas          *int32            `yaml:"replicas,omitempty"`
	Expose            bool              `yaml:"expose,omitempty"`
	Labels            map[string]string `yaml:"labels,omitempty"`
	PersistentVolumes []Volume          `yaml:"persistentVolumes,omitempty"`
	api_v1.PodSpec    `yaml:",inline"`
}

func ReadFile(f string) ([]byte, error) {
	data, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, errors.Wrap(err, "file reading failed")
	}
	return data, nil
}

func Convert(v *viper.Viper, cmd *cobra.Command) error {

	for _, file := range strings.Split(v.GetStringSlice("files")[0], ",") {
		d, err := ReadFile(file)
		if err != nil {
			return errors.New(err.Error())
		}

		var app App
		err = yaml.Unmarshal(d, &app)
		if err != nil {
			return errors.Wrap(err, "could not unmarshal into internal struct")
		}
		log.Debugf("file: %s, object unmrashalled: %#v", file, app)

		runtimeObjects, err := CreateK8sObjects(&app)

		for _, runtimeObject := range runtimeObjects {
			gvk, isUnversioned, err := api.Scheme.ObjectKind(runtimeObject)
			if err != nil {
				return errors.Wrap(err, "ConvertToVersion failed")
			}
			if isUnversioned {
				return errors.New(fmt.Sprintf("ConvertToVersion failed: can't output unversioned type: %T", runtimeObject))
			}

			runtimeObject.GetObjectKind().SetGroupVersionKind(gvk)

			data, err := yaml.Marshal(runtimeObject)
			if err != nil {
				return errors.Wrap(err, "failed to marshal object")
			}

			writeObject := func(o runtime.Object, data []byte) error {
				_, err := fmt.Fprintln(os.Stdout, "---")
				if err != nil {
					return errors.Wrap(err, "could not print to STDOUT")
				}

				_, err = os.Stdout.Write(data)
				return errors.Wrap(err, "could not write to STDOUT")
			}

			err = writeObject(runtimeObject, data)
			if err != nil {
				return errors.Wrap(err, "failed to write object")
			}
		}
	}

	return nil
}

func getLabels(app *App) map[string]string {
	labels := map[string]string{"app": app.Name}
	return labels
}

func initService(app *App) *api_v1.Service {
	return &api_v1.Service{
		ObjectMeta: api_v1.ObjectMeta{
			Name:   app.Name,
			Labels: app.Labels,
		},
		Spec: api_v1.ServiceSpec{
			Selector: app.Labels,
		},
	}
}

func allPorts(app *App) []api_v1.ContainerPort {
	var ports []api_v1.ContainerPort
	for _, c := range app.Containers {
		for _, p := range c.Ports {
			ports = append(ports, p)
		}
	}
	return ports
}

func createDeployment(app *App) *ext_v1beta1.Deployment {
	// bare minimum deployment
	return &ext_v1beta1.Deployment{
		ObjectMeta: api_v1.ObjectMeta{
			Name:   app.Name,
			Labels: app.Labels,
		},
		Spec: ext_v1beta1.DeploymentSpec{
			Replicas: app.Replicas,
			Template: api_v1.PodTemplateSpec{
				ObjectMeta: api_v1.ObjectMeta{
					Name:   app.Name,
					Labels: app.Labels,
				},
				// get pod spec out of the original info
				Spec: api_v1.PodSpec(app.PodSpec),
			},
		},
	}
}

func CreateK8sObjects(app *App) ([]runtime.Object, error) {

	var objects []runtime.Object
	var svc *api_v1.Service

	if app.Labels == nil {
		app.Labels = getLabels(app)
	}

	ports := allPorts(app)
	if len(ports) > 0 {
		// bare minimum service
		svc = initService(app)
		if app.Expose {
			svc.Spec.Type = api_v1.ServiceTypeLoadBalancer
			// TODO: create ingress
		}
		// update the service based on the ports given in the app
		for _, p := range ports {
			// adding the ports to the service
			svc.Spec.Ports = append(svc.Spec.Ports, api_v1.ServicePort{
				Name:       fmt.Sprintf("port-%d", p.ContainerPort),
				Port:       int32(p.ContainerPort),
				TargetPort: intstr.FromInt(int(p.ContainerPort)),
			})
		}
	}

	var vols []api_v1.VolumeMount
	for _, c := range app.Containers {
		// get all the volumes and create a pvc
		vols = append(vols, c.VolumeMounts...)
	}

	// update the app with volumes info
	var pvcs []runtime.Object
	for _, vm := range vols {
		// update the pod, volumes info
		app.Volumes = append(app.Volumes, api_v1.Volume{
			Name: vm.Name,
			VolumeSource: api_v1.VolumeSource{
				PersistentVolumeClaim: &api_v1.PersistentVolumeClaimVolumeSource{
					ClaimName: vm.Name,
				},
			},
		})

		// create pvc
		size, err := resource.ParseQuantity("100Mi")
		if err != nil {
			return nil, errors.Wrap(err, "could not read volume size")
		}

		pvc := &api_v1.PersistentVolumeClaim{
			ObjectMeta: api_v1.ObjectMeta{
				Name: vm.Name,
			},
			Spec: api_v1.PersistentVolumeClaimSpec{
				Resources: api_v1.ResourceRequirements{
					Requests: api_v1.ResourceList{
						api_v1.ResourceStorage: size,
					},
				},
				AccessModes: []api_v1.PersistentVolumeAccessMode{api_v1.ReadWriteOnce},
			},
		}
		pvcs = append(pvcs, pvc)
	}
	// if only one container set name of it as app name
	if len(app.Containers) == 1 && app.Containers[0].Name == "" {
		app.Containers[0].Name = app.Name
	}

	deployment := createDeployment(app)

	objects = append(objects, deployment)
	log.Debugf("app: %s, deployment: %s", app.Name, spew.Sprint(deployment))
	objects = append(objects, svc)
	log.Debugf("app: %s, service: %s", app.Name, spew.Sprint(svc))
	objects = append(objects, pvcs...)
	log.Debugf("app: %s, pvc: %s", app.Name, spew.Sprint(pvcs))

	return objects, nil
}

// how to specify volume size
// how to expose certain service using ingress
//
