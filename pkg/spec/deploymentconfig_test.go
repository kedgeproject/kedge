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
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	os_deploy_v1 "github.com/openshift/origin/pkg/apps/apis/apps/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	api_v1 "k8s.io/kubernetes/pkg/api/v1"
)

func TestFixDeploymentConfig(t *testing.T) {
	tests := []struct {
		name           string
		input          *App
		expectedOutput *App
	}{
		{
			name: "No replicas passed at input, expected 1",
			input: &App{
				DeploymentConfigs: []DeploymentConfigSpecMod{
					{},
				},
			},
			expectedOutput: &App{
				DeploymentConfigs: []DeploymentConfigSpecMod{
					{
						DeploymentConfigSpec: os_deploy_v1.DeploymentConfigSpec{
							Replicas: 1,
						},
						Replicas: getInt32Addr(1),
					},
				},
			},
		},
		{
			name: "replicas set to 0 by the end user, expected 0",
			input: &App{
				DeploymentConfigs: []DeploymentConfigSpecMod{
					{
						Replicas: getInt32Addr(0),
					},
				},
			},
			expectedOutput: &App{
				DeploymentConfigs: []DeploymentConfigSpecMod{
					{
						DeploymentConfigSpec: os_deploy_v1.DeploymentConfigSpec{
							Replicas: 0,
						},
						Replicas: getInt32Addr(0),
					},
				},
			},
		},
		{
			name: "replicas set to 2 by the end user, expected 2",
			input: &App{
				DeploymentConfigs: []DeploymentConfigSpecMod{
					{
						Replicas: getInt32Addr(2),
					},
				},
			},
			expectedOutput: &App{
				DeploymentConfigs: []DeploymentConfigSpecMod{
					{
						DeploymentConfigSpec: os_deploy_v1.DeploymentConfigSpec{
							Replicas: 2,
						},
						Replicas: getInt32Addr(2),
					},
				},
			},
		},
		{
			name: "test Name, Labels, and Annotations  propagation from root level",
			input: &App{
				ObjectMeta: meta_v1.ObjectMeta{
					Name: "test",
					Labels: map[string]string{
						"foo": "bar",
					},
					Annotations: map[string]string{
						"abc": "def",
					},
				},
				DeploymentConfigs: []DeploymentConfigSpecMod{
					{},
				},
			},
			expectedOutput: &App{
				ObjectMeta: meta_v1.ObjectMeta{
					Name: "test",
					Labels: map[string]string{
						"foo": "bar",
					},
					Annotations: map[string]string{
						"abc": "def",
					},
				},
				DeploymentConfigs: []DeploymentConfigSpecMod{
					{
						ObjectMeta: meta_v1.ObjectMeta{
							Name: "test",
							Labels: map[string]string{
								"foo": "bar",
							},
							Annotations: map[string]string{
								"abc": "def",
							},
						},
						DeploymentConfigSpec: os_deploy_v1.DeploymentConfigSpec{
							Replicas: 1,
						},
						Replicas: getInt32Addr(1),
					},
				},
			},
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			test.input.fixDeploymentConfigs()
			if !reflect.DeepEqual(test.input, test.expectedOutput) {
				t.Errorf("Expected output to be:\n%v\nBut got:\n%v\n",
					prettyPrintObjects(test.expectedOutput),
					prettyPrintObjects(test.input))
			}
		})
	}
}

// &App{
// 	ObjectMeta: meta_v1.ObjectMeta{
// 		Name: "test",
// 	},
// 	DeploymentConfigs: []DeploymentConfigSpecMod{
// 		{
// 			DeploymentConfigSpec: os_deploy_v1.DeploymentConfigSpec{
// 				Replicas: 2,
// 				Template: &api_v1.PodTemplateSpec{
// 					Spec: api_v1.PodSpec{
// 						Containers: []api_v1.Container{
// 							{
// 								Image: "testImage",
// 							},
// 						},
// 					},
// 				},
// 			},
// 		},
// 	},
func TestDeploymentConfigSpecMod_CreateOpenShiftController(t *testing.T) {
	tests := []struct {
		name        string
		app         *App
		deployments []runtime.Object
		success     bool
	}{
		{
			name: "Test that it correctly converts",
			app: &App{
				ObjectMeta: meta_v1.ObjectMeta{
					Name: "testJob",
				},
				DeploymentConfigs: []DeploymentConfigSpecMod{
					{
						PodSpecMod: PodSpecMod{
							Containers: []Container{
								{
									Container: api_v1.Container{
										Name:  "testContainer",
										Image: "testImage",
									},
								},
							},
						},
						Replicas: getInt32Addr(2),
					},
				},
			},
			deployments: []runtime.Object{
				&os_deploy_v1.DeploymentConfig{
					ObjectMeta: meta_v1.ObjectMeta{
						Name: "testJob",
					},
					Spec: os_deploy_v1.DeploymentConfigSpec{
						Replicas: 2,
						Template: &api_v1.PodTemplateSpec{
							ObjectMeta: meta_v1.ObjectMeta{
								Name: "testJob",
							},
							Spec: api_v1.PodSpec{
								Containers: []api_v1.Container{
									{
										Name:  "testContainer",
										Image: "testImage",
									},
								},
							},
						},
					},
				},
			},
			success: true,
		},
		{
			name: "Test that strategy is converted correctly",
			app: &App{
				ObjectMeta: meta_v1.ObjectMeta{
					Name: "testJob",
				},
				DeploymentConfigs: []DeploymentConfigSpecMod{
					{
						PodSpecMod: PodSpecMod{
							Containers: []Container{
								{
									Container: api_v1.Container{
										Name:  "testContainer",
										Image: "testImage",
									},
								},
							},
						},
						Replicas: getInt32Addr(3),
						DeploymentConfigSpec: os_deploy_v1.DeploymentConfigSpec{
							Strategy: os_deploy_v1.DeploymentStrategy{
								Type: os_deploy_v1.DeploymentStrategyType("Rolling"),
							},
						},
					},
				},
			},
			deployments: []runtime.Object{
				&os_deploy_v1.DeploymentConfig{
					ObjectMeta: meta_v1.ObjectMeta{
						Name: "testJob",
					},
					Spec: os_deploy_v1.DeploymentConfigSpec{
						Replicas: 3,
						Strategy: os_deploy_v1.DeploymentStrategy{
							Type: os_deploy_v1.DeploymentStrategyType("Rolling"),
						},
						Template: &api_v1.PodTemplateSpec{
							ObjectMeta: meta_v1.ObjectMeta{
								Name: "testJob",
							},
							Spec: api_v1.PodSpec{
								Containers: []api_v1.Container{
									{
										Name:  "testContainer",
										Image: "testImage",
									},
								},
							},
						},
					},
				},
			},
			success: true,
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			test.app.fixDeploymentConfigs()
			dcs, err := test.app.createDeploymentConfigs()

			switch test.success {
			case true:
				if err != nil {
					t.Errorf("Expected test to pass but got an error -\n%v", err)
				}
			case false:
				if err == nil {
					t.Errorf("For the input -\n%v\nexpected test to fail, but test passed", spew.Sprint(test.app))
				}
			}

			if !reflect.DeepEqual(test.deployments, dcs) {
				t.Errorf("Expected OpenShift DeploymentConfig to be -\n%v\nBut got -\n%v", prettyPrintObjects(test.deployments), prettyPrintObjects(dcs))
			}
		})
	}
}
