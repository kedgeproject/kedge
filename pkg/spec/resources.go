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
	"strings"

	log "github.com/Sirupsen/logrus"
	build_v1 "github.com/openshift/origin/pkg/build/apis/build/v1"
	image_v1 "github.com/openshift/origin/pkg/image/apis/image/v1"
	os_route_v1 "github.com/openshift/origin/pkg/route/apis/route/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	api_v1 "k8s.io/kubernetes/pkg/api/v1"
	ext_v1beta1 "k8s.io/kubernetes/pkg/apis/extensions/v1beta1"

	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/kubernetes/pkg/api"

	// install api (register and add types to api.Schema)
	_ "k8s.io/kubernetes/pkg/api/install"
	_ "k8s.io/kubernetes/pkg/apis/extensions/install"
)

// allLabelKey is the key that Kedge injects in every Kubernetes resource that
// it generates as an ObjectMeta label
const (
	appLabelKey   = "app"
	appVersion    = "appversion"
	BuildLabelKey = "build"
)

// Fix

func (app *App) fixServices() error {

	// auto populate name only if one service is specified without any name
	if len(app.Services) == 1 && app.Services[0].ObjectMeta.Name == "" {
		app.Services[0].ObjectMeta.Name = app.Name
	}

	for i, service := range app.Services {
		if service.ObjectMeta.Name == "" {
			return fmt.Errorf("please specify name for app.services[%d]", i)
		}

		for key, value := range app.ObjectMeta.Labels {
			service.ObjectMeta.Labels = addKeyValueToMap(key, value, service.ObjectMeta.Labels)
		}

		// copy root annotations (already has "appVersion: <app.AppVersion>" annotation)
		for key, value := range app.ObjectMeta.Annotations {
			service.ObjectMeta.Annotations = addKeyValueToMap(key, value, service.ObjectMeta.Annotations)
		}

		// this should be the last statement in this for loop
		app.Services[i] = service
	}
	return nil
}

func fixVolumeClaims(volumeClaims []VolumeClaim, appName string) ([]VolumeClaim, error) {

	// auto populate name only if one volumeClaim is specified without any name
	if len(volumeClaims) == 1 && volumeClaims[0].ObjectMeta.Name == "" {
		volumeClaims[0].ObjectMeta.Name = appName
	}

	for i, pVolume := range volumeClaims {
		if pVolume.ObjectMeta.Name == "" {
			return nil, fmt.Errorf("please specify name for app.volumeClaims[%d]", i)
		}

		pVolume.ObjectMeta.Labels = addKeyValueToMap(appLabelKey, appName, pVolume.ObjectMeta.Labels)

		// this should be the last statement in this for loop
		volumeClaims[i] = pVolume
	}
	return volumeClaims, nil
}

func fixConfigMaps(configMaps []ConfigMapMod, appName string) ([]ConfigMapMod, error) {

	// auto populate name only if one configMap is specified without any name
	if len(configMaps) == 1 && configMaps[0].ObjectMeta.Name == "" {
		configMaps[0].ObjectMeta.Name = appName
	}

	for i, cm := range configMaps {
		if cm.ObjectMeta.Name == "" {
			return nil, fmt.Errorf("please specify name for app.configMaps[%d]", i)
		}

		cm.ObjectMeta.Labels = addKeyValueToMap(appLabelKey, appName, cm.ObjectMeta.Labels)

		// this should be the last statement in this for loop
		configMaps[i] = cm
	}
	return configMaps, nil
}

func fixSecrets(secrets []SecretMod, appName string) ([]SecretMod, error) {

	// auto populate name only if one secret is specified without any name
	if len(secrets) == 1 && secrets[0].ObjectMeta.Name == "" {
		secrets[0].ObjectMeta.Name = appName
	}

	for i, sec := range secrets {
		if sec.Name == "" {
			return nil, fmt.Errorf("please specify name for app.secrets[%d]", i)
		}

		sec.ObjectMeta.Labels = addKeyValueToMap(appLabelKey, appName, sec.ObjectMeta.Labels)

		// this should be the last statement in this for loop
		secrets[i] = sec
	}
	return secrets, nil
}

func fixImageStreams(imageStreams []ImageStreamSpecMod, appName string) ([]ImageStreamSpecMod, error) {

	// auto populate name only if one ImageStream is specified without any name
	if len(imageStreams) == 1 && imageStreams[0].Name == "" {
		imageStreams[0].ObjectMeta.Name = appName
	}

	for i, is := range imageStreams {
		if is.Name == "" {
			return nil, fmt.Errorf("please specify name for app.imageStreams[%d]", i)
		}

		is.ObjectMeta.Labels = addKeyValueToMap(appLabelKey, appName, is.ObjectMeta.Labels)

		// this should be the last statement in this for loop
		imageStreams[i] = is
	}
	return imageStreams, nil
}

func fixBuildConfigs(buildConfigs []BuildConfigSpecMod, appName string) ([]BuildConfigSpecMod, error) {

	// auto populate name only if one buildConfig is specified without any name
	if len(buildConfigs) == 1 && buildConfigs[0].Name == "" {
		buildConfigs[0].ObjectMeta.Name = appName
	}

	for i, bc := range buildConfigs {
		if bc.Name == "" {
			return nil, fmt.Errorf("please specify name for app.buildConfigs[%d]", i)
		}

		bc.ObjectMeta.Labels = addKeyValueToMap(appLabelKey, appName, bc.ObjectMeta.Labels)

		// this should be the last statement in this for loop
		buildConfigs[i] = bc
	}
	return buildConfigs, nil
}

func fixIngresses(ingresses []IngressSpecMod, appName string) ([]IngressSpecMod, error) {

	// auto populate name only if one ingress is specified without any name
	if len(ingresses) == 1 && ingresses[0].Name == "" {
		ingresses[0].ObjectMeta.Name = appName
	}

	for i, ing := range ingresses {
		if ing.Name == "" {
			return nil, fmt.Errorf("please specify name for app.ingresses[%d]", i)
		}

		ing.ObjectMeta.Labels = addKeyValueToMap(appLabelKey, appName, ing.ObjectMeta.Labels)

		// this should be the last statement in this for loop
		ingresses[i] = ing
	}
	return ingresses, nil
}

func fixRoutes(routes []RouteSpecMod, appName string) ([]RouteSpecMod, error) {

	// auto populate name only if one route is specified without any name
	if len(routes) == 1 && routes[0].Name == "" {
		routes[0].ObjectMeta.Name = appName
	}

	for i, route := range routes {
		if route.Name == "" {
			return nil, fmt.Errorf("please specify name for app.routes[%d]", i)
		}

		route.ObjectMeta.Labels = addKeyValueToMap(appLabelKey, appName, route.ObjectMeta.Labels)

		// this should be the last statement in this for loop
		routes[i] = route
	}
	return routes, nil
}

func fixContainers(containers []Container, appName string) ([]Container, error) {

	// auto populate name only if one ingress is specified without any name
	if len(containers) == 1 && containers[0].Name == "" {
		containers[0].Name = appName
	}

	for i, c := range containers {
		if c.Name == "" {
			return nil, fmt.Errorf("please specify name for app.containers[%d]", i)
		}
	}
	return containers, nil
}

func (app *App) Fix() error {
	var err error

	app.ObjectMeta.Labels = addKeyValueToMap(appLabelKey, app.Name, app.ObjectMeta.Labels)

	if app.Appversion != "" {
		app.ObjectMeta.Annotations = addKeyValueToMap(appVersion, app.Appversion, app.ObjectMeta.Annotations)
	}

	PrettyPrintObjects(&app.ObjectMeta)

	// fix Services
	err = app.fixServices()
	if err != nil {
		return errors.Wrap(err, "Unable to fix services")
	}

	// fix VolumeClaims
	app.VolumeClaims, err = fixVolumeClaims(app.VolumeClaims, app.Name)
	if err != nil {
		return errors.Wrap(err, "Unable to fix persistentVolume")
	}

	// fix configMaps
	app.ConfigMaps, err = fixConfigMaps(app.ConfigMaps, app.Name)
	if err != nil {
		return errors.Wrap(err, "unable to fix configMaps")
	}

	err = app.fixDeployments()
	if err != nil {
		return errors.Wrap(err, "unable to fix deployments")
	}

	err = app.fixDeploymentConfigs()
	if err != nil {
		return errors.Wrap(err, "unable to fix deploymentConfigs")
	}

	err = app.fixJobs()
	if err != nil {
		return errors.Wrap(err, "unable to fix jobs")
	}

	app.Secrets, err = fixSecrets(app.Secrets, app.Name)
	if err != nil {
		return errors.Wrap(err, "unable to fix secrets")
	}

	// fix imageStreams
	app.ImageStreams, err = fixImageStreams(app.ImageStreams, app.Name)
	if err != nil {
		return errors.Wrap(err, "unable to fix imageStreams")
	}

	// fix buildConfigs
	app.BuildConfigs, err = fixBuildConfigs(app.BuildConfigs, app.Name)
	if err != nil {
		return errors.Wrap(err, "unable to fix buildConfigs")
	}

	// fix ingresses
	app.Ingresses, err = fixIngresses(app.Ingresses, app.Name)
	if err != nil {
		return errors.Wrap(err, "unable to fix ingresses")
	}

	// fix routes
	app.Routes, err = fixRoutes(app.Routes, app.Name)
	if err != nil {
		return errors.Wrap(err, "unable to fix routes")
	}

	return nil
}

// Transform

func (app *App) getLabels() map[string]string {
	labels := map[string]string{appLabelKey: app.Name}
	return labels
}

func (app *App) createIngresses() ([]runtime.Object, error) {
	var ings []runtime.Object

	for _, i := range app.Ingresses {
		ing := &ext_v1beta1.Ingress{
			ObjectMeta: i.ObjectMeta,
			Spec:       i.IngressSpec,
		}
		ings = append(ings, ing)
	}
	return ings, nil
}

func (app *App) createRoutes() ([]runtime.Object, error) {
	var routes []runtime.Object

	for _, r := range app.Routes {
		route := &os_route_v1.Route{
			ObjectMeta: r.ObjectMeta,
			Spec:       r.RouteSpec,
		}
		routes = append(routes, route)
	}
	return routes, nil
}

func (app *App) createServices() ([]runtime.Object, error) {
	var svcs []runtime.Object
	for _, s := range app.Services {
		svc := &api_v1.Service{
			ObjectMeta: s.ObjectMeta,
			Spec:       s.ServiceSpec,
		}
		for _, servicePortMod := range s.Ports {
			svc.Spec.Ports = append(svc.Spec.Ports, servicePortMod.ServicePort)
		}

		for _, portMapping := range s.PortMappings {
			servicePort, err := parsePortMapping(portMapping.String())
			if err != nil {
				return nil, errors.Wrap(err, "unable to parse port mapping")
			}
			svc.Spec.Ports = append(svc.Spec.Ports, *servicePort)
		}

		populateServicePortNames(svc.Name, svc.Spec.Ports)

		if len(svc.Spec.Selector) == 0 {
			svc.Spec.Selector = app.Labels
		}
		svcs = append(svcs, svc)

		// Generate ingress if "endpoint" is mentioned in app.Services.Ports[].Endpoint
		for _, port := range s.Ports {
			if port.IngressEndpoint != "" {
				var host string
				var path string
				endpoint := strings.SplitN(port.IngressEndpoint, "/", 2)
				switch len(endpoint) {
				case 1:
					host = endpoint[0]
					path = "/"
				case 2:
					host = endpoint[0]
					path = "/" + endpoint[1]
				default:
					return nil, fmt.Errorf("Invalid syntax for endpoint: %v", port.IngressEndpoint)
				}

				ingressName := s.Name + "-" + strconv.FormatInt(int64(port.Port), 10)
				endpointIngress := &ext_v1beta1.Ingress{
					ObjectMeta: metav1.ObjectMeta{
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

			if port.RouteEndpoint != "" {

				endpointRoute := &os_route_v1.Route{
					ObjectMeta: metav1.ObjectMeta{
						Name:   svc.Name,
						Labels: app.Labels,
					},
					Spec: os_route_v1.RouteSpec{
						To: os_route_v1.RouteTargetReference{
							Kind: "Service",
							Name: svc.Name,
						},
					},
				}

				switch port.RouteEndpoint {
				case "true":
					// OpenShift will take care of route
					svcs = append(svcs, endpointRoute)
				default:
					// Route URL is given by user
					endpointRoute.Spec.Host = port.RouteEndpoint
					svcs = append(svcs, endpointRoute)

				}
			}
		}

	}
	return svcs, nil
}

// create PVC reading the root level persistent volume field
func (app *App) createPVC() ([]runtime.Object, error) {
	var pvcs []runtime.Object
	for _, v := range app.VolumeClaims {
		// check for conditions where user has given both conflicting fields
		// or not given either fields
		if v.Size != "" && v.Resources.Requests != nil {
			return nil, fmt.Errorf("persistent volume %q, cannot provide size and resources at the same time", v.Name)
		}
		if v.Size == "" && v.Resources.Requests == nil {
			return nil, fmt.Errorf("persistent volume %q, please provide size or resources, none given", v.Name)
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
		pvcs = append(pvcs, &api_v1.PersistentVolumeClaim{
			ObjectMeta: v.ObjectMeta,
			// since we updated the pvc spec before so this can be directly fed
			// without having to do any addition extra
			Spec: api_v1.PersistentVolumeClaimSpec(v.PersistentVolumeClaimSpec),
		})
	}
	return pvcs, nil
}

func (app *App) createSecrets() ([]runtime.Object, error) {
	var secrets []runtime.Object

	for _, s := range app.Secrets {
		secret := &api_v1.Secret{
			ObjectMeta: s.ObjectMeta,
			Data:       s.Data,
			StringData: s.StringData,
			Type:       s.Type,
		}
		secrets = append(secrets, secret)
	}
	return secrets, nil
}

func (app *App) createImageStreams() ([]runtime.Object, error) {
	var imageStreams []runtime.Object

	for _, is := range app.ImageStreams {
		imageStream := &image_v1.ImageStream{
			ObjectMeta: is.ObjectMeta,
			Spec:       is.ImageStreamSpec,
		}
		imageStreams = append(imageStreams, imageStream)
	}
	return imageStreams, nil
}

func (app *App) createBuildConfigs() ([]runtime.Object, error) {
	var buildConfigs []runtime.Object

	for _, b := range app.BuildConfigs {
		buildConfig := &build_v1.BuildConfig{
			ObjectMeta: b.ObjectMeta,
			Spec:       b.BuildConfigSpec,
		}
		buildConfigs = append(buildConfigs, buildConfig)
	}
	return buildConfigs, nil
}

// CreateK8sObjects, if given object DeploymentSpecMod, this function reads
// them and returns kubernetes objects as list of runtime.Object
// If the deployment is using field 'includeResources' then it will
// also return file names mentioned there as list of string
func (app *App) CreateK8sObjects() ([]runtime.Object, error) {
	var objects []runtime.Object

	if app.Labels == nil {
		app.Labels = app.getLabels()
	}

	svcs, err := app.createServices()
	if err != nil {
		return nil, errors.Wrap(err, "Unable to create Kubernetes Service")
	}

	ings, err := app.createIngresses()
	if err != nil {
		return nil, errors.Wrap(err, "Unable to create Kubernetes Ingresses")
	}

	routes, err := app.createRoutes()
	if err != nil {
		return nil, errors.Wrap(err, "unable to create OpenShift Routes")
	}

	secs, err := app.createSecrets()
	if err != nil {
		return nil, errors.Wrap(err, "Unable to create Kubernetes Secrets")
	}

	iss, err := app.createImageStreams()
	if err != nil {
		return nil, errors.Wrap(err, "unable to create OpenShift ImageStreams")
	}

	bcs, err := app.createBuildConfigs()
	if err != nil {
		return nil, errors.Wrap(err, "Unable to create OpenShift BuildConfigs")
	}

	// create pvc for each root level persistent volume
	pvcs, err := app.createPVC()
	if err != nil {
		return nil, errors.Wrap(err, "unable to create Persistent Volume Claims")
	}

	deployments, err := app.createDeployments()
	if err != nil {
		return nil, errors.Wrap(err, "Unable to create Kubernetes Deployments")
	}

	deploymentConfigs, err := app.createDeploymentConfigs()
	if err != nil {
		return nil, errors.Wrap(err, "Unable to create OpenShift DeploymentConfigs")
	}

	jobs, err := app.createJobs()
	if err != nil {
		return nil, errors.Wrap(err, "Unable to create Kubernetes Jobs")
	}

	var configMap []runtime.Object
	for _, cd := range app.ConfigMaps {
		cm := &api_v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:   cd.Name,
				Labels: app.Labels,
			},
			Data: cd.Data,
		}

		configMap = append(configMap, cm)
	}

	// please keep the order of the artifacts addition as it is

	objects = append(objects, pvcs...)
	log.Debugf("app: %s, pvc: %s\n", app.Name, spew.Sprint(pvcs))

	objects = append(objects, svcs...)
	log.Debugf("app: %s, service: %s\n", app.Name, spew.Sprint(svcs))

	objects = append(objects, ings...)
	log.Debugf("app: %s, ingress: %s\n", app.Name, spew.Sprint(ings))

	objects = append(objects, routes...)
	log.Debugf("app: %s, routes: %s\n", app.Name, spew.Sprint(routes))

	objects = append(objects, secs...)
	log.Debugf("app: %s, secret: %s\n", app.Name, spew.Sprint(secs))

	objects = append(objects, iss...)
	log.Debugf("app: %s, imageStreams: %s\n", app.Name, spew.Sprint(iss))

	objects = append(objects, bcs...)
	log.Debugf("app: %s, buildConfig: %s\n", app.Name, spew.Sprint(bcs))

	objects = append(objects, configMap...)
	log.Debugf("app: %s, configMap: %s\n", app.Name, spew.Sprint(configMap))

	objects = append(objects, deployments...)
	log.Debugf("app: %s, deployments: %s\n", app.Name, spew.Sprint(deployments))

	objects = append(objects, deploymentConfigs...)
	log.Debugf("app: %s, deploymentConfigs: %s\n", app.Name, spew.Sprint(deploymentConfigs))

	objects = append(objects, jobs...)
	log.Debugf("app: %s, jobs: %s\n", app.Name, spew.Sprint(jobs))

	//Adding annotations to all the resources
	//Objects are runtimeobjects, so accessing them using meta library
	if app.Appversion != "" {
		accessor := meta.NewAccessor()
		for _, object := range objects {
			annotations, err := accessor.Annotations(object)
			if err != nil {
				return nil, errors.Wrap(err, "cannot get annotations")
			}
			annotations = addKeyValueToMap(appVersion, app.Appversion, annotations)
			err = accessor.SetAnnotations(object, annotations)
			if err != nil {
				return nil, errors.Wrap(err, "cannot set annotations")
			}
		}
	}

	return objects, nil
}

// Validate

func validateVolumeClaims(vcs []VolumeClaim) error {
	// find the duplicate volume claim names, if found any then error out
	vcmap := make(map[string]interface{})
	for _, vc := range vcs {
		if _, ok := vcmap[vc.Name]; !ok {
			// value here does not matter
			vcmap[vc.Name] = nil
		} else {
			return fmt.Errorf("duplicate entry of volume claim %q", vc.Name)
		}
	}

	return nil
}

func (app *App) Validate() error {

	// validate volumeclaims
	if err := validateVolumeClaims(app.VolumeClaims); err != nil {
		return errors.Wrap(err, "error validating volume claims")
	}

	if err := app.validateJobs(); err != nil {
		return errors.Wrap(err, "error validating Jobs")
	}

	return nil
}

// Others

// Parse the string the get the port, targetPort and protocol
// information, and then return the resulting ServicePort object
func parsePortMapping(pm string) (*api_v1.ServicePort, error) {

	// The current syntax for portMapping is - port:targetPort/protocol
	// The only field mandatory here is "port". There are 4 possible cases here
	// which are handled in this function.

	// Case 1 - port
	// Case 2 - port:targetPort
	// Case 3 - port/protocol
	// Case 4 - port:targetPort/protocol

	var port int32
	var targetPort intstr.IntOrString
	var protocol api_v1.Protocol

	protocolSplit := strings.Split(pm, "/")
	switch len(protocolSplit) {

	// When no protocol is specified, we set the protocol to TCP
	// Case 1 - port
	// Case 2 - port:targetPort
	case 1:
		protocol = api_v1.ProtocolTCP

	// When protocol is specified
	// Case 3 - port/protocol
	// Case 4 - port:targetPort/protocol
	case 2:
		protocolUpperCase := strings.ToUpper(protocolSplit[1])
		switch api_v1.Protocol(protocolUpperCase) {
		case api_v1.ProtocolTCP, api_v1.ProtocolUDP:
			protocol = api_v1.Protocol(protocolUpperCase)
		default:
			return nil, fmt.Errorf("invalid protocol '%v' provided, the acceptable values are '%v' and '%v'", protocolSplit[1], api.ProtocolTCP, api.ProtocolUDP)
		}
	// There is no case in which splitting by "/" provides < 1 or > 2 values
	default:
		return nil, fmt.Errorf("invalid syntax for protocol '%v' provided, use 'port:targetPort/protocol'", pm)
	}

	portSplit := strings.Split(pm, ":")
	switch len(portSplit) {

	// When only port is specified
	// Case 1 - port
	// Case 3 - port/protocol
	case 1:
		// Ignoring the protocol part, if present, and converting only the port
		// part
		p, err := strconv.ParseInt(strings.Split(portSplit[0], "/")[0], 10, 32)
		if err != nil {
			return nil, errors.Wrap(err, "port is not an int")
		}

		port, targetPort.IntVal = int32(p), int32(p)

	// When port and targetPort both are specified
	// Case 2 - port:targetPort
	// Case 4 - port:targetPort/protocol
	case 2:
		p, err := strconv.ParseInt(portSplit[0], 10, 32)
		if err != nil {
			return nil, errors.Wrap(err, "port is not an int")
		}
		port = int32(p)

		// Ignoring the protocol part, if present, and converting only the
		// targetPort part
		tp, err := strconv.ParseInt(strings.Split(portSplit[1], "/")[0], 10, 32)
		if err != nil {
			return nil, errors.Wrap(err, "targetPort is not an int")
		}
		targetPort.IntVal = int32(tp)

	// There is no case in which splitting by ": provides < 1 or > 2 values
	default:
		return nil, fmt.Errorf("invalid syntax for portMapping '%v', use 'port:targetPort/protocol'", pm)
	}

	return &api_v1.ServicePort{
		Port:       port,
		TargetPort: targetPort,
		Protocol:   protocol,
	}, nil
}
