package spec

import (
	api_v1 "k8s.io/client-go/pkg/api/v1"
	ext_v1beta1 "k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

type PersistentVolume struct {
	api_v1.PersistentVolumeClaimSpec `json:",inline"`
	Name                             string `json:"name"`
	Size                             string `json:"size"`
}

type ServicePortMod struct {
	api_v1.ServicePort `json:",inline"`
	// Endpoint allows specifying an ingress resource in the format
	// `<Host>/<Path>`
	Endpoint string `json:"endpoint"`
}

type ServiceSpecMod struct {
	api_v1.ServiceSpec `json:",inline"`
	Name               string           `json:"name,omitempty"`
	Ports              []ServicePortMod `json:"ports"`
}

type IngressSpecMod struct {
	Name                    string `json:"name"`
	ext_v1beta1.IngressSpec `json:",inline"`
}

type Container struct {
	// one common definitions for livenessProbe and readinessProbe
	// this allows to have only one place to define both probes (if they are the same)
	Health           *api_v1.Probe `json:"health,omitempty"`
	api_v1.Container `json:",inline"`
}

type ConfigMapMod struct {
	Name string            `json:"name,omitempty"`
	Data map[string]string `json:"data,omitempty"`
}

type App struct {
	Name              string             `json:"name"`
	Replicas          *int32             `json:"replicas,omitempty"`
	Labels            map[string]string  `json:"labels,omitempty"`
	PersistentVolumes []PersistentVolume `json:"persistentVolumes,omitempty"`
	ConfigMaps        []ConfigMapMod     `json:"configMaps,omitempty"`
	Services          []ServiceSpecMod   `json:"services,omitempty"`
	Ingress           []IngressSpecMod   `json:"ingress,omitempty"`

	// overwrite containers from PodSpec
	Containers []Container `json:"containers,omitempty"`

	api_v1.PodSpec             `json:",inline"`
	ext_v1beta1.DeploymentSpec `json:",inline"`
}
