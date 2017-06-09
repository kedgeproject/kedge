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

	log "github.com/Sirupsen/logrus"
	api_v1 "k8s.io/client-go/pkg/api/v1"
	ext_v1beta1 "k8s.io/client-go/pkg/apis/extensions/v1beta1"

	// install api

	_ "k8s.io/client-go/pkg/api/install"
	_ "k8s.io/client-go/pkg/apis/extensions/install"
)

type Volume struct {
	api_v1.Volume `yaml:",inline"`
	Size          string   `yaml:"size"`
	AccessModes   []string `yaml:"accessModes"`
}

type Service struct {
	Name               string `yaml:"name,omitempty"`
	api_v1.ServiceSpec `yaml:",inline"`
}

type App struct {
	Name              string            `yaml:"name"`
	Replicas          *int32            `yaml:"replicas,omitempty"`
	Expose            bool              `yaml:"expose,omitempty"`
	Labels            map[string]string `yaml:"labels,omitempty"`
	PersistentVolumes []Volume          `yaml:"persistentVolumes,omitempty"`
	ConfigData        map[string]string `yaml:"configData,omitempty"`
	Services          []Service         `yaml:"services,omitempty"`
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

func createServices(app *App) []runtime.Object {
	var svcs []runtime.Object
	for _, s := range app.Services {
		svc := &api_v1.Service{
			ObjectMeta: api_v1.ObjectMeta{
				Name:   s.Name,
				Labels: app.Labels,
			},
			Spec: s.ServiceSpec,
		}
		if len(svc.Spec.Selector) == 0 {
			svc.Spec.Selector = app.Labels
		}
		svcs = append(svcs, svc)

		if s.Type == api_v1.ServiceTypeLoadBalancer {
			// TODO: create a ingress resource
		}
	}
	return svcs
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

func isVolumeDefined(app *App, name string) bool {
	for _, v := range app.PersistentVolumes {
		if name == v.Name {
			return true
		}
	}
	return false
}

func searchVolumeIndex(app *App, name string) int {
	for i := 0; i < len(app.PersistentVolumes); i++ {
		if name == app.PersistentVolumes[i].Name {
			return i
		}
	}
	return -1
}

func createPVC(v *Volume) (*api_v1.PersistentVolumeClaim, error) {
	// create pvc
	if v.Size == "" {
		v.Size = "100Mi"
	}
	size, err := resource.ParseQuantity(v.Size)
	if err != nil {
		return nil, errors.Wrap(err, "could not read volume size")
	}

	pvc := &api_v1.PersistentVolumeClaim{
		ObjectMeta: api_v1.ObjectMeta{
			Name: v.Name,
		},
		Spec: api_v1.PersistentVolumeClaimSpec{
			Resources: api_v1.ResourceRequirements{
				Requests: api_v1.ResourceList{
					api_v1.ResourceStorage: size,
				},
			},
		},
	}
	for _, mode := range v.AccessModes {
		switch mode {
		case "ReadWriteOnce":
			pvc.Spec.AccessModes = append(pvc.Spec.AccessModes, api_v1.ReadWriteOnce)
		case "ReadOnlyMany":
			pvc.Spec.AccessModes = append(pvc.Spec.AccessModes, api_v1.ReadOnlyMany)
		case "ReadWriteMany":
			pvc.Spec.AccessModes = append(pvc.Spec.AccessModes, api_v1.ReadWriteMany)
		}
	}
	if len(v.AccessModes) == 0 {
		pvc.Spec.AccessModes = []api_v1.PersistentVolumeAccessMode{api_v1.ReadWriteOnce}
	}

	return pvc, nil
}

func add(objects *[]runtime.Object, object runtime.Object) {
	if object != nil {
		*objects = append(*objects, object)
		log.Debugf("%s\n", spew.Sprint(object))
	}
}

func CreateK8sObjects(app *App) ([]runtime.Object, error) {

	var objects []runtime.Object

	if app.Labels == nil {
		app.Labels = getLabels(app)
	}

	svcs := createServices(app)

	var pvcs []runtime.Object
	for _, c := range app.Containers {
		for _, vm := range c.VolumeMounts {

			// User won't be giving this so we have to create it
			// so that the pod spec is complete
			podVolume := api_v1.Volume{
				Name: vm.Name,
				VolumeSource: api_v1.VolumeSource{
					PersistentVolumeClaim: &api_v1.PersistentVolumeClaimVolumeSource{
						ClaimName: vm.Name,
					},
				},
			}
			app.Volumes = append(app.Volumes, podVolume)

			if isVolumeDefined(app, vm.Name) {
				i := searchVolumeIndex(app, vm.Name)
				pvc, err := createPVC(&app.PersistentVolumes[i])
				if err != nil {
					return nil, errors.Wrap(err, "cannot create pvc")
				}
				pvcs = append(pvcs, pvc)
				continue
			}

			v := Volume{podVolume, "100Mi", []string{"ReadWriteOnce"}}
			app.PersistentVolumes = append(app.PersistentVolumes, v)
			pvc, err := createPVC(&v)
			if err != nil {
				return nil, errors.Wrap(err, "cannot create pvc")
			}
			pvcs = append(pvcs, pvc)
		}
	}

	// if only one container set name of it as app name
	if len(app.Containers) == 1 && app.Containers[0].Name == "" {
		app.Containers[0].Name = app.Name
	}

	//var configMap *api_v1.ConfigMap
	//if len(app.ConfigData) > 0 {
	//	configMap = &api_v1.ConfigMap{
	//		ObjectMeta: api_v1.ObjectMeta{
	//			Name: app.Name,
	//		},
	//		Data: app.ConfigData,
	//	}
	//}

	deployment := createDeployment(app)

	add(&objects, deployment)
	log.Debugf("app: %s, deployment: %s", app.Name, spew.Sprint(deployment))
	//add(&objects, configMap)
	//log.Debugf("app: %s, configMap: %s", spew.Sprint(configMap))
	for _, svc := range svcs {
		add(&objects, svc)
	}
	log.Debugf("app: %s, service: %s", app.Name, spew.Sprint(svcs))
	for _, pvc := range pvcs {
		add(&objects, pvc)
	}
	log.Debugf("app: %s, pvc: %s", app.Name, spew.Sprint(pvcs))

	return objects, nil
}

// how to expose certain service using ingress
//
