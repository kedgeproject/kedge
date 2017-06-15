package kubernetes

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	"github.com/surajssd/opencomposition/pkg/spec"
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

func getLabels(app *spec.App) map[string]string {
	labels := map[string]string{"app": app.Name}
	return labels
}

func createServices(app *spec.App) []runtime.Object {
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

func createDeployment(app *spec.App) *ext_v1beta1.Deployment {
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

// search through all the persistent volumes defined in the root level
func isPVCDefined(app *spec.App, name string) bool {
	for _, v := range app.PersistentVolumes {
		if v.Name == name {
			return true
		}
	}
	return false
}

// create PVC reading the root level persistent volume field
func createPVC(v spec.PersistentVolume, labels map[string]string) (*api_v1.PersistentVolumeClaim, error) {
	// check for conditions where user has given both conflicting fields
	// or not given either fields
	if v.Size != "" && v.Resources.Requests != nil {
		return nil, errors.New(fmt.Sprintf("persistent volume %q, cannot provide size and resources at the same time", v.Name))
	}
	if v.Size == "" && v.Resources.Requests == nil {
		return nil, errors.New(fmt.Sprintf("persistent volume %q, please provide size or resources, none given", v.Name))
	}

	// if user has given size then create a "api_v1.ResourceRequirements"
	// because this can be fed to pvc directly
	if v.Size != "" {
		size, err := resource.ParseQuantity(v.Size)
		if err != nil {
			return nil, errors.Wrap(err, "could not read volume size")
		}
		// update the volume's resource so that it can be fed
		v.Resources = api_v1.ResourceRequirements{
			Requests: api_v1.ResourceList{
				api_v1.ResourceStorage: size,
			},
		}
	}
	// setting the default accessmode if none given by user
	if len(v.AccessModes) == 0 {
		v.AccessModes = []api_v1.PersistentVolumeAccessMode{api_v1.ReadWriteOnce}
	}
	pvc := &api_v1.PersistentVolumeClaim{
		ObjectMeta: api_v1.ObjectMeta{
			Name:   v.Name,
			Labels: labels,
		},
		// since we updated the pvc spec before so this can be directly fed
		// without having to do any addition extra
		Spec: api_v1.PersistentVolumeClaimSpec(v.PersistentVolumeClaimSpec),
	}
	return pvc, nil
}

// This function will search in the pod level volumes
// and see if the volume with given name is defined
func isVolumeDefined(app *spec.App, name string) bool {
	for _, v := range app.Volumes {
		if v.Name == name {
			return true
		}
	}
	return false
}

func isAnyConfigMapRef(app *spec.App) bool {
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

// Since we are automatically creating pvc from
// root level persistent volume and entry in the container
// volume mount, we alse need to update the pod's volume field
func populateVolumes(app *spec.App) error {
	for cn, c := range app.Containers {
		for vn, vm := range c.VolumeMounts {
			if isPVCDefined(app, vm.Name) {
				app.Volumes = append(app.Volumes, api_v1.Volume{
					Name: vm.Name,
					VolumeSource: api_v1.VolumeSource{
						PersistentVolumeClaim: &api_v1.PersistentVolumeClaimVolumeSource{
							ClaimName: vm.Name,
						},
					},
				})
			} else if !isVolumeDefined(app, vm.Name) {
				// pvc is not defined so we need to check if the entry is made in the pod volumes
				// since a volumeMount entry without entry in pod level volumes might cause failure
				// while deployment since that would not be a complete configuration
				return errors.New(fmt.Sprintf("neither root level Persistent Volume"+
					" nor Volume in pod spec defined for %q, "+
					"in app.containers[%d].volumeMounts[%d]", vm.Name, cn, vn))
			}
		}
	}
	return nil
}

func CreateK8sObjects(app *spec.App) ([]runtime.Object, error) {
	var objects []runtime.Object

	if app.Labels == nil {
		app.Labels = getLabels(app)
	}

	svcs := createServices(app)

	// create pvc for each root level persistent volume
	var pvcs []runtime.Object
	for _, v := range app.PersistentVolumes {
		pvc, err := createPVC(v, app.Labels)
		if err != nil {
			return nil, errors.Wrapf(err, "app %q", app.Name)
		}
		pvcs = append(pvcs, pvc)
	}
	if err := populateVolumes(app); err != nil {
		return nil, errors.Wrapf(err, "app %q", app.Name)
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
								LocalObjectReference: api_v1.LocalObjectReference{
									Name: app.Name,
								},
								Key: k,
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

func Transform(app *spec.App) ([]runtime.Object, error) {

	runtimeObjects, err := CreateK8sObjects(app)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Kubernetes objects")
	}

	for _, runtimeObject := range runtimeObjects {

		gvk, isUnversioned, err := api.Scheme.ObjectKind(runtimeObject)
		if err != nil {
			return nil, errors.Wrap(err, "ConvertToVersion failed")
		}
		if isUnversioned {
			return nil, errors.New(fmt.Sprintf("ConvertToVersion failed: can't output unversioned type: %T", runtimeObject))
		}

		runtimeObject.GetObjectKind().SetGroupVersionKind(gvk)
	}

	return runtimeObjects, nil
}
