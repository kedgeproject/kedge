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
	os_deploy_v1 "github.com/openshift/origin/pkg/apps/apis/apps/v1"
	build_v1 "github.com/openshift/origin/pkg/build/apis/build/v1"
	image_v1 "github.com/openshift/origin/pkg/image/apis/image/v1"
	os_route_v1 "github.com/openshift/origin/pkg/route/apis/route/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	api_v1 "k8s.io/kubernetes/pkg/api/v1"
	batch_v1 "k8s.io/kubernetes/pkg/apis/batch/v1"
	ext_v1beta1 "k8s.io/kubernetes/pkg/apis/extensions/v1beta1"
)

// VolumeClaim is used to define Persistent Volumes for app
// kedgeSpec: io.kedge.VolumeClaim
type VolumeClaim struct {
	// Data from the kubernetes persistent volume claim spec
	// k8s: io.k8s.kubernetes.pkg.api.v1.PersistentVolumeClaimSpec
	api_v1.PersistentVolumeClaimSpec `json:",inline"`
	// Size of persistent volume
	Size string `json:"size"`
	// k8s: io.k8s.apimachinery.pkg.apis.meta.v1.ObjectMeta
	meta_v1.ObjectMeta `json:",inline"`
}

// ServicePortMod is used to define Kubernetes service's port
// kedgeSpec: io.kedge.ServicePort
type ServicePortMod struct {
	// k8s: io.k8s.kubernetes.pkg.api.v1.ServicePort
	api_v1.ServicePort `json:",inline"`
	// Host to create ingress automatically. Endpoint allows specifying an
	// ingress resource in the format '<Host>/<Path>'
	// +optional
	Endpoint string `json:"endpoint"`
}

// ServiceSpecMod is used to define Kubernetes service
// kedgeSpec: io.kedge.ServiceSpec
type ServiceSpecMod struct {
	// k8s: io.k8s.kubernetes.pkg.api.v1.ServiceSpec
	api_v1.ServiceSpec `json:",inline"`
	// The list of ports that are exposed by this service. More info:
	// https://kubernetes.io/docs/concepts/services-networking/service/#virtual-ips-and-service-proxies
	// ref: io.kedge.ServicePort
	// +optional
	Ports []ServicePortMod `json:"ports,conflicting"`
	// The list of portMappings, where each portMapping allows specifying port,
	// targetPort and protocol in the format '<port>:<targetPort>/<protocol>'
	// +optional
	PortMappings []intstr.IntOrString `json:"portMappings,omitempty"`
	// k8s: io.k8s.apimachinery.pkg.apis.meta.v1.ObjectMeta
	meta_v1.ObjectMeta `json:",inline"`
}

// IngressSpecMod defines Kubernetes Ingress object
// kedgeSpec: io.kedge.IngressSpec
type IngressSpecMod struct {
	// k8s: io.k8s.kubernetes.pkg.apis.extensions.v1beta1.IngressSpec
	ext_v1beta1.IngressSpec `json:",inline"`
	// k8s: io.k8s.apimachinery.pkg.apis.meta.v1.ObjectMeta
	meta_v1.ObjectMeta `json:",inline"`
}

// kedgeSpec: io.kedge.RouteSpec
type RouteSpecMod struct {
	// k8s: v1.RouteSpec
	os_route_v1.RouteSpec `json:",inline"`
	// k8s: io.k8s.apimachinery.pkg.apis.meta.v1.ObjectMeta
	meta_v1.ObjectMeta `json:",inline"`
}

// Container defines a single application container that you want to run within a pod.
// kedgeSpec: io.kedge.ContainerSpec
type Container struct {
	// One common definitions for 'livenessProbe' and 'readinessProbe'
	// this allows to have only one place to define both probes (if they are the same)
	// Periodic probe of container liveness and readiness. Container will be restarted
	// if the probe fails. Cannot be updated. More info:
	// https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes
	// ref: io.k8s.kubernetes.pkg.api.v1.Probe
	// +optional
	Health *api_v1.Probe `json:"health,omitempty"`
	// k8s: io.k8s.kubernetes.pkg.api.v1.Container
	api_v1.Container `json:",inline"`
}

// ConfigMapMod holds configuration data for pods to consume.
// kedgeSpec: io.kedge.ConfigMap
type ConfigMapMod struct {
	// k8s: io.k8s.kubernetes.pkg.api.v1.ConfigMap
	api_v1.ConfigMap `json:",inline"`
	// k8s: io.k8s.apimachinery.pkg.apis.meta.v1.ObjectMeta
	// We need ObjectMeta here, even though it is present in api.ConfigMap
	// because the one upstream has a JSON tag "metadata" due to which it
	// cannot be merged at ConfigMap's root level. The ObjectMeta here
	// overwrites the one in upstream and lets us merge ObjectMeta at
	// ConfigMap's root YAML syntax
	meta_v1.ObjectMeta `json:",inline"`
}

// PodSpecMod is a description of a pod
type PodSpecMod struct {
	// List of containers belonging to the pod. Containers cannot currently be
	// added or removed. There must be at least one container in a Pod. Cannot be updated.
	// ref: io.kedge.ContainerSpec
	Containers []Container `json:"containers,conflicting,omitempty"`
	// List of initialization containers belonging to the pod. Init containers are
	// executed in order prior to containers being started. If any init container
	// fails, the pod is considered to have failed and is handled according to its
	// restartPolicy. The name for an init container or normal container must be
	// unique among all containers.
	// ref: io.kedge.ContainerSpec
	// +optional
	InitContainers []Container `json:"initContainers,conflicting,omitempty"`
	// k8s: io.k8s.kubernetes.pkg.api.v1.PodSpec
	api_v1.PodSpec `json:",inline"`
}

// SecretMod defines secret that will be consumed by application
// kedgeSpec: io.kedge.SecretSpec
type SecretMod struct {
	// k8s: io.k8s.kubernetes.pkg.api.v1.Secret
	api_v1.Secret `json:",inline"`
	// k8s: io.k8s.apimachinery.pkg.apis.meta.v1.ObjectMeta
	// We need ObjectMeta here, even though it is present in api.Secret
	// because the one upstream has a JSON tag "metadata" due to which it
	// cannot be merged at Secret's root level. The ObjectMeta here
	// overwrites the one in upstream and lets us merge ObjectMeta at
	// Secret's root YAML syntax
	meta_v1.ObjectMeta `json:",inline"`
}

// ImageStreamSpec defines OpenShift ImageStream Object
// kedgeSpec: io.kedge.ImageStreamSpec
type ImageStreamSpecMod struct {
	// k8s: v1.ImageStreamSpec
	image_v1.ImageStreamSpec `json:",inline"`
	// k8s: io.k8s.apimachinery.pkg.apis.meta.v1.ObjectMeta
	meta_v1.ObjectMeta `json:",inline"`
}

// BuildConfigSpecMod defines OpenShift BuildConfig object
// kedgeSpec: io.kedge.BuildConfigSpec
type BuildConfigSpecMod struct {
	// k8s: v1.BuildConfigSpec
	build_v1.BuildConfigSpec `json:",inline"`
	// k8s: io.k8s.apimachinery.pkg.apis.meta.v1.ObjectMeta
	meta_v1.ObjectMeta `json:",inline"`
}

// ControllerFields are the common fields in every controller Kedge supports
type ControllerFields struct {
	// Field to specify the version of application
	// +optional
	Appversion string `json:"appversion,omitempty"`

	Controller string `json:"controller,omitempty"`
	// List of volume that should be mounted on the pod.
	// ref: io.kedge.VolumeClaim
	// +optional
	VolumeClaims []VolumeClaim `json:"volumeClaims,omitempty"`
	// List of configMaps
	// ref: io.kedge.ConfigMap
	// +optional
	ConfigMaps []ConfigMapMod `json:"configMaps,omitempty"`
	// List of Kubernetes Services
	// ref: io.kedge.ServiceSpec
	// +optional
	Services []ServiceSpecMod `json:"services,omitempty"`
	// List of Kubernetes Ingress
	// ref: io.kedge.IngressSpec
	// +optional
	Ingresses []IngressSpecMod `json:"ingresses,omitempty"`
	// List of OpenShift Routes
	// ref: io.kedge.RouteSpec
	// +optional
	Routes []RouteSpecMod `json:"routes,omitempty"`
	// List of Kubernetes Secrets
	// ref: io.kedge.SecretSpec
	// +optional
	Secrets []SecretMod `json:"secrets,omitempty"`
	// List of OpenShift ImageStreams
	// ref: io.kedge.ImageStreamSpec
	// +optional
	ImageStreams []ImageStreamSpecMod `json:"imageStreams,omitempty"`
	// List of OpenShift BuildConfigs
	// ref: io.kedge.BuildConfigSpec
	// +optional
	BuildConfigs []BuildConfigSpecMod `json:"buildConfigs,omitempty"`
	// List of Kubernetes resource files, that can be directly given to Kubernetes
	// +optional
	IncludeResources []string `json:"includeResources,omitempty"`

	// k8s: io.k8s.apimachinery.pkg.apis.meta.v1.ObjectMeta
	meta_v1.ObjectMeta `json:",inline"`
	PodSpecMod         `json:",inline"`
}

type Controller struct {
	Controller string `json:"controller,omitempty"`
}

// DeploymentSpecMod is Kedge's extension of Kubernetes DeploymentSpec and allows
// defining a complete kedge application
// kedgeSpec: io.kedge.DeploymentSpecMod
type DeploymentSpecMod struct {
	ControllerFields `json:",inline"`
	// k8s: io.k8s.kubernetes.pkg.apis.apps.v1beta1.DeploymentSpec
	ext_v1beta1.DeploymentSpec `json:",inline"`
}

// JobSpecMod is Kedge's extension of Kubernetes JobSpec and allows
// defining a complete kedge application
// kedgeSpec: io.kedge.JobSpecMod
type JobSpecMod struct {
	ControllerFields `json:",inline"`
	// k8s: io.k8s.kubernetes.pkg.apis.batch.v1.JobSpec
	batch_v1.JobSpec `json:",inline"`
	// Optional duration in seconds relative to the startTime that the job may be active
	// before the system tries to terminate it; value must be positive integer
	// This only sets ActiveDeadlineSeconds in JobSpec, not PodSpec
	// +optional
	ActiveDeadlineSeconds *int64 `json:"activeDeadlineSeconds,conflicting,omitempty"`
}

// Ochestrator: OpenShift
// DeploymentConfigSpecMod is Kedge's extension of OpenShift DeploymentConfig in order to define and allow
// a complete kedge app based on OpenShift
// kedgeSpec: io.kedge.DeploymentConfigSpecMod
type DeploymentConfigSpecMod struct {
	ControllerFields                  `json:",inline"`
	os_deploy_v1.DeploymentConfigSpec `json:",inline"`

	// Replicas is the number of desired replicas.
	// We need to add this field here despite being in v1.DeploymentConfigSpec
	// because the one in v1.DeploymentConfigSpec has the type as int32, which
	// does not let us check if the set value is 0, is it set by the user or not
	// since this field's value with default to 0. We need the default value as
	// 1. Hence, we need to check if the user has set it or not.Making the type
	// *int32 helps us do it, followed by substitution later on.
	Replicas *int32 `json:"replicas,omitempty"`
}
