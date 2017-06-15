package spec

import (
	api_v1 "k8s.io/client-go/pkg/api/v1"
	ext_v1beta1 "k8s.io/client-go/pkg/apis/extensions/v1beta1"
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
	Labels            map[string]string `yaml:"labels,omitempty"`
	PersistentVolumes []Volume          `yaml:"persistentVolumes,omitempty"`
	ConfigData        map[string]string `yaml:"configData,omitempty"`
	Services          []Service         `yaml:"services,omitempty"`
	api_v1.PodSpec    `yaml:",inline"`
}
