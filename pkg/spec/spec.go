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

type Service struct {
	Name                    string `json:"name,omitempty"`
	api_v1.ServiceSpec      `json:",inline"`
	ext_v1beta1.IngressSpec `json:",inline"`
}

type Container struct {
	Health           *api_v1.Probe `json:"health,omitempty"`
	api_v1.Container `json:",inline"`
}

type App struct {
	Name              string             `json:"name"`
	Replicas          *int32             `json:"replicas,omitempty"`
	Labels            map[string]string  `json:"labels,omitempty"`
	PersistentVolumes []PersistentVolume `json:"persistentVolumes,omitempty"`
	ConfigData        map[string]string  `json:"configData,omitempty"`
	Services          []Service          `json:"services,omitempty"`
	Containers        []Container        `json:"containers,omitempty"`
	api_v1.PodSpec    `json:",inline"`
}
