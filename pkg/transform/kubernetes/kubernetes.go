package kubernetes

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	"github.com/surajssd/kapp/pkg/spec"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/resource"
	"k8s.io/client-go/pkg/runtime"
	"k8s.io/client-go/pkg/util/intstr"

	log "github.com/Sirupsen/logrus"
	api_v1 "k8s.io/client-go/pkg/api/v1"
	ext_v1beta1 "k8s.io/client-go/pkg/apis/extensions/v1beta1"

	// install api (register and add types to api.Schema)
	_ "k8s.io/client-go/pkg/api/install"
	_ "k8s.io/client-go/pkg/apis/extensions/install"
)

func getLabels(app *spec.App) map[string]string {
	labels := map[string]string{"app": app.Name}
	return labels
}

func createIngresses(app *spec.App) ([]runtime.Object, error) {
	var ings []runtime.Object

	for _, i := range app.Ingresses {
		ing := &ext_v1beta1.Ingress{
			ObjectMeta: api_v1.ObjectMeta{
				Name:   i.Name,
				Labels: app.Labels,
			},
			Spec: i.IngressSpec,
		}
		ings = append(ings, ing)
	}
	return ings, nil
}

func createServices(app *spec.App) ([]runtime.Object, error) {
	var svcs []runtime.Object
	for _, s := range app.Services {
		svc := &api_v1.Service{
			ObjectMeta: api_v1.ObjectMeta{
				Name:   s.Name,
				Labels: app.Labels,
			},
			Spec: s.ServiceSpec,
		}
		for _, servicePortMod := range s.Ports {
			svc.Spec.Ports = append(svc.Spec.Ports, servicePortMod.ServicePort)
		}
		if len(svc.Spec.Selector) == 0 {
			svc.Spec.Selector = app.Labels
		}
		svcs = append(svcs, svc)

		// Generate ingress if "endpoint" is mentioned in app.Services.Ports[].Endpoint
		for _, port := range s.Ports {
			if port.Endpoint != "" {
				var host string
				var path string
				endpoint := strings.SplitN(port.Endpoint, "/", 2)
				switch len(endpoint) {
				case 1:
					host = endpoint[0]
					path = "/"
				case 2:
					host = endpoint[0]
					path = "/" + endpoint[1]
				default:
					return nil, errors.New(fmt.Sprintf("Invalid syntax for endpoint: %v", port.Endpoint))
				}

				ingressName := s.Name + "-" + strconv.FormatInt(int64(port.Port), 10)
				endpointIngress := &ext_v1beta1.Ingress{
					ObjectMeta: api_v1.ObjectMeta{
						Name:   ingressName,
						Labels: app.Labels,
					},
					Spec: ext_v1beta1.IngressSpec{
						Rules: []ext_v1beta1.IngressRule{
							{
								Host: host,
								IngressRuleValue: ext_v1beta1.IngressRuleValue{
									HTTP: &ext_v1beta1.HTTPIngressRuleValue{
										Paths: []ext_v1beta1.HTTPIngressPath{
											{
												Path: path,
												Backend: ext_v1beta1.IngressBackend{
													ServiceName: s.Name,
													ServicePort: intstr.IntOrString{
														IntVal: port.Port,
													},
												},
											},
										},
									},
								},
							},
						},
					},
				}
				svcs = append(svcs, endpointIngress)
			}
		}
	}
	return svcs, nil
}

func createDeployment(app *spec.App) (*ext_v1beta1.Deployment, error) {

	// We are merging whole DeploymentSpec with PodSpec.
	// This means that someone could specify containers in template.spec and also in top level PodSpec.
	// This stupid check is supposed to make sure that only one of them set.
	// TODO: merge DeploymentSpec.Template.Spec and top level PodSpec
	if !(reflect.DeepEqual(app.DeploymentSpec.Template.Spec, api_v1.PodSpec{}) || reflect.DeepEqual(app.PodSpec, api_v1.PodSpec{})) {
		return nil, fmt.Errorf("Pod can't be specfied in two places. Use top level PodSpec or template.spec (DeploymentSpec.Template.Spec) not both")
	}

	deploymentSpec := app.DeploymentSpec

	// top level PodSpec is not empty, use it for deployment template
	// we already know that if app.PodSpec is not empty app.DeploymentSpec.Template.Spec is empty
	if !reflect.DeepEqual(app.PodSpec, api_v1.PodSpec{}) {
		deploymentSpec.Template.Spec = app.PodSpec
	}

	// TODO: check if this wasn't set by user, in that case we shouldn't ovewrite it
	deploymentSpec.Template.ObjectMeta.Name = app.Name

	// TODO: merge with already existing labels and avoid duplication
	deploymentSpec.Template.ObjectMeta.Labels = app.Labels

	deployment := ext_v1beta1.Deployment{
		ObjectMeta: api_v1.ObjectMeta{
			Name:   app.Name,
			Labels: app.Labels,
		},
		Spec: deploymentSpec,
	}

	return &deployment, nil
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
	for _, c := range app.PodSpec.Containers {
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
	for cn, c := range app.PodSpec.Containers {
		for vn, vm := range c.VolumeMounts {
			if isPVCDefined(app, vm.Name) && !isVolumeDefined(app, vm.Name) {
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

func populateContainerHealth(app *spec.App) error {
	for cn, c := range app.Containers {
		// check if health and liveness given together
		if c.Health != nil && (c.ReadinessProbe != nil || c.LivenessProbe != nil) {
			return errors.New(fmt.Sprintf("cannot define field health and livnessProbe"+
				" or readinessProbe together in app.containers[%d]", cn))
		}
		if c.Health != nil {
			c.LivenessProbe = c.Health
			c.ReadinessProbe = c.Health
		}
		app.PodSpec.Containers = append(app.PodSpec.Containers, c.Container)
	}
	return nil
}

func CreateK8sObjects(app *spec.App) ([]runtime.Object, error) {
	var objects []runtime.Object

	if app.Labels == nil {
		app.Labels = getLabels(app)
	}

	svcs, err := createServices(app)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to create Kubernetes Service")
	}

	ings, err := createIngresses(app)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to create Kubernetes Ingresses")
	}

	// withdraw the health and populate actual pod spec
	if err := populateContainerHealth(app); err != nil {
		return nil, errors.Wrapf(err, "app %q", app.Name)
	}

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
	if len(app.PodSpec.Containers) == 1 && app.PodSpec.Containers[0].Name == "" {
		app.PodSpec.Containers[0].Name = app.Name
	} else if len(app.PodSpec.Containers) > 1 {
		// check if all the containers have a name
		// if not fail giving error
		for cn, c := range app.PodSpec.Containers {
			if c.Name == "" {
				return nil, fmt.Errorf("app %q: container name not defined for app.containers[%d]", app.Name, cn)
			}
		}
	}

	var configMap []runtime.Object
	for _, cd := range app.ConfigMaps {
		cm := &api_v1.ConfigMap{
			ObjectMeta: api_v1.ObjectMeta{
				Name:   cd.Name,
				Labels: app.Labels,
			},
			Data: cd.Data,
		}

		configMap = append(configMap, cm)
	}

	deployment, err := createDeployment(app)
	if err != nil {
		return nil, errors.Wrapf(err, "app %q", app.Name)
	}
	objects = append(objects, deployment)
	log.Debugf("app: %s, deployment: %s\n", app.Name, spew.Sprint(deployment))

	objects = append(objects, configMap...)
	log.Debugf("app: %s, configMap: %s\n", app.Name, spew.Sprint(configMap))

	objects = append(objects, svcs...)
	log.Debugf("app: %s, service: %s\n", app.Name, spew.Sprint(svcs))

	objects = append(objects, ings...)
	log.Debugf("app: %s, ingress: %s\n", app.Name, spew.Sprint(ings))

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
