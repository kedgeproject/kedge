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
	api_v1 "k8s.io/client-go/pkg/api/v1"
	ext_v1beta1 "k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

// VolumeClaim is used to define Persistent Volumes for app
// kedgeSpec: io.kedge.VolumeClaim
type VolumeClaim struct {
	// Data from the kubernetes persistent volume claim spec
	// k8s: io.k8s.kubernetes.pkg.api.v1.PersistentVolumeClaimSpec
	api_v1.PersistentVolumeClaimSpec `json:",inline"`
	// Name of the persistent Volume Claim
	// +optional
	Name string `json:"name"`
	// Size of persistent volume
	Size string `json:"size"`
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
	// Name of the service
	// +optional
	Name string `json:"name,omitempty"`
	// The list of ports that are exposed by this service. More info:
	// https://kubernetes.io/docs/concepts/services-networking/service/#virtual-ips-and-service-proxies
	// ref: io.kedge.ServicePort
	Ports []ServicePortMod `json:"ports,conflicting"`
	// The list of portMappings, where each portMapping allows specifying port,
	// targetPort and protocol in the format '<port>:<targetPort>/<protocol>'
	PortMappings []string `json:"portMappings,omitempty"`
}

// IngressSpecMod defines Kubernetes Ingress object
// kedgeSpec: io.kedge.IngressSpec
type IngressSpecMod struct {
	// Name of the ingress
	// +optional
	Name string `json:"name"`
	// k8s: io.k8s.kubernetes.pkg.apis.extensions.v1beta1.IngressSpec
	ext_v1beta1.IngressSpec `json:",inline"`
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
	// Name of the configMap
	// +optional
	Name string `json:"name,omitempty"`
	// Data contains the configuration data. Each key must consist of alphanumeric characters, '-', '_' or '.'
	Data map[string]string `json:"data,omitempty"`
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
	// Name of the secret
	// +optional
	Name string `json:"name,omitempty"`
	// k8s: io.k8s.kubernetes.pkg.api.v1.Secret
	api_v1.Secret `json:",inline"`
}

// ControllerFields are the common fields in every controller Kedge supports
type ControllerFields struct {
	// Name of the micro-service
	Name string `json:"name"`
	// Specify which Kubernetes controller to generate. Defaults to deployment.
	// +optional
	Controller string `json:"controller,omitempty"`
	// Map of string keys and values that can be used to organize and categorize
	// (scope and select) objects. May match selectors of replication controllers
	// and services. More info: http://kubernetes.io/docs/user-guide/labels
	// +optional
	Labels map[string]string `json:"labels,omitempty"`
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
	// List of Kubernetes Secrets
	// ref: io.kedge.SecretSpec
	// +optional
	Secrets []SecretMod `json:"secrets,omitempty"`
	// List of Kubernetes resource files, that can be directly given to Kubernetes
	// +optional
	ExtraResources []string `json:"extraResources,omitempty"`

	PodSpecMod `json:",inline"`
}

// DeploymentSpecMod is Kedge's extension of Kubernetes DeploymentSpec and allows
// defining a complete kedge application
// kedgeSpec: io.kedge.DeploymentSpec
type DeploymentSpecMod struct {
	ControllerFields `json:",inline"`
	// k8s: io.k8s.kubernetes.pkg.apis.apps.v1beta1.DeploymentSpec
	ext_v1beta1.DeploymentSpec `json:",inline"`
}

type Controller struct {
	Controller string `json:"controller,omitempty"`
}
