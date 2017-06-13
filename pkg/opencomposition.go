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

var (
	DefaultVolumeSize string = "100Mi"
	DefaultVolumeType string = "ReadWriteOnce"
)

type Volume struct {
	api_v1.Volume `yaml:",inline"`
	Size          string   `yaml:"size"`
	AccessModes   []string `yaml:"accessModes"`
}

type Service struct {
	Name                    string `yaml:"name,omitempty"`
	api_v1.ServiceSpec      `yaml:",inline"`
	ext_v1beta1.IngressSpec `yaml:",inline"`
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
		log.Debugf("file: %s, object unmrashalled: %#v\n", file, app)

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
			// if more than one port given then we enforce user to specify in the http

			// autogenerate
			if len(s.Rules) == 1 && len(s.Ports) == 1 {
				http := s.Rules[0].HTTP
				if http == nil {
					http = &ext_v1beta1.HTTPIngressRuleValue{
						Paths: []ext_v1beta1.HTTPIngressPath{
							{
								Path: "/",
								Backend: ext_v1beta1.IngressBackend{
									ServiceName: s.Name,
									ServicePort: intstr.FromInt(int(s.Ports[0].Port)),
								},
							},
						},
					}
				}
				ing := &ext_v1beta1.Ingress{
					ObjectMeta: api_v1.ObjectMeta{
						Name:   s.Name,
						Labels: app.Labels,
					},
					Spec: ext_v1beta1.IngressSpec{
						Rules: []ext_v1beta1.IngressRule{
							{
								IngressRuleValue: ext_v1beta1.IngressRuleValue{
									HTTP: http,
								},
								Host: s.Rules[0].Host,
							},
						},
					},
				}
				svcs = append(svcs, ing)
			} else if len(s.Rules) == 1 && len(s.Ports) > 1 {
				if s.Rules[0].HTTP == nil {
					log.Warnf("No HTTP given for multiple ports")
				}
			} else if len(s.Rules) > 1 {
				ing := &ext_v1beta1.Ingress{
					ObjectMeta: api_v1.ObjectMeta{
						Name:   s.Name,
						Labels: app.Labels,
					},
					Spec: s.IngressSpec,
				}
				svcs = append(svcs, ing)
			}
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
	if i := searchVolumeIndex(app, name); i != -1 {
		return true
	}
	return false
}

func searchVolumeIndex(app *App, name string) int {
	for i, v := range app.PersistentVolumes {
		if name == v.Name {
			return i
		}
	}
	return -1
}

func createPVC(v *Volume) (*api_v1.PersistentVolumeClaim, error) {
	if v.Size == "" {
		v.Size = DefaultVolumeSize
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

func isAnyConfigMapRef(app *App) bool {
	for _, c := range app.Containers {
		for _, env := range c.Env {
			if env.ValueFrom != nil && env.ValueFrom.ConfigMapKeyRef != nil && env.ValueFrom.ConfigMapKeyRef.Name == app.Name {
				return true
			}
		}
	}
	for _, v := range app.Volumes {
		if v.ConfigMap != nil && v.ConfigMap.Name == app.Name {
			return true
		}
	}

	return false
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

			// Retrieve a default configuration
			v := Volume{podVolume, DefaultVolumeSize, []string{DefaultVolumeType}}

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

	var configMap *api_v1.ConfigMap
	if len(app.ConfigData) > 0 {
		configMap = &api_v1.ConfigMap{
			ObjectMeta: api_v1.ObjectMeta{
				Name: app.Name,
			},
			Data: app.ConfigData,
		}

		// add it to the envs if there is no configMapRef
		// we cannot re-create the entries for configMap
		// because there is no way we will know which container wants to use it
		if len(app.Containers) == 1 && !isAnyConfigMapRef(app) {
			// iterate over the data in the configMap
			for k, _ := range app.ConfigData {
				app.Containers[0].Env = append(app.Containers[0].Env,
					api_v1.EnvVar{
						Name: k,
						ValueFrom: &api_v1.EnvVarSource{
							ConfigMapKeyRef: &api_v1.ConfigMapKeySelector{
								api_v1.LocalObjectReference{
									Name: app.Name,
								},
								k,
							},
						},
					})
			}
		} else if len(app.Containers) > 1 && !isAnyConfigMapRef(app) {
			log.Warnf("You have defined a configMap but you have not mentioned where you gonna consume it!")
		}

	}

	deployment := createDeployment(app)
	objects = append(objects, deployment)
	log.Debugf("app: %s, deployment: %s\n", app.Name, spew.Sprint(deployment))

	if configMap != nil {
		objects = append(objects, configMap)
	}
	log.Debugf("app: %s, configMap: %s\n", app.Name, spew.Sprint(configMap))

	objects = append(objects, svcs...)
	log.Debugf("app: %s, service: %s\n", app.Name, spew.Sprint(svcs))

	objects = append(objects, pvcs...)
	log.Debugf("app: %s, pvc: %s\n", app.Name, spew.Sprint(pvcs))

	return objects, nil
}
