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
	batch_v1 "k8s.io/client-go/pkg/apis/batch/v1"
	ext_v1beta1 "k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

type VolumeClaim struct {
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
	Ports              []ServicePortMod `json:"ports,conflicting"`
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

type PodSpecMod struct {
	Containers     []Container `json:"containers,conflicting,omitempty"`
	InitContainers []Container `json:"initContainers,conflicting,omitempty"`
	api_v1.PodSpec `json:",inline"`
}

type SecretMod struct {
	Name          string `json:"name,omitempty"`
	api_v1.Secret `json:",inline"`
}

type Controller struct {
	Controller string `json:"controller,omitempty"`
}

type JobSpecMod struct {
	Name             string `json:"name,omitempty"`
	api_v1.PodSpec   `json:",inline"`
	batch_v1.JobSpec `json:",inline"`
	ExtraResources   []string `json:"extraResources,omitempty"`
}

type DeploymentSpecMod struct {
	Name                       string            `json:"name"`
	Labels                     map[string]string `json:"labels,omitempty"`
	VolumeClaims               []VolumeClaim     `json:"volumeClaims,omitempty"`
	ConfigMaps                 []ConfigMapMod    `json:"configMaps,omitempty"`
	Services                   []ServiceSpecMod  `json:"services,omitempty"`
	Ingresses                  []IngressSpecMod  `json:"ingresses,omitempty"`
	Secrets                    []SecretMod       `json:"secrets,omitempty"`
	ExtraResources             []string          `json:"extraResources,omitempty"`
	PodSpecMod                 `json:",inline"`
	ext_v1beta1.DeploymentSpec `json:",inline"`
}
