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
	"reflect"
	"testing"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	api_v1 "k8s.io/kubernetes/pkg/api/v1"

	"github.com/davecgh/go-spew/spew"
)

func TestPopulateProbes(t *testing.T) {
	t.Logf("Running failing tests")
	failingTests := []struct {
		name  string
		input Container
	}{
		{
			name: "health and livenessProbe given together",
			input: Container{
				Health: &api_v1.Probe{
					Handler: api_v1.Handler{
						Exec: &api_v1.ExecAction{Command: []string{"mysqladmin", "ping"}},
					}, InitialDelaySeconds: 30, TimeoutSeconds: 5,
				},
				Container: api_v1.Container{
					LivenessProbe: &api_v1.Probe{
						Handler: api_v1.Handler{
							Exec: &api_v1.ExecAction{Command: []string{"mysqladmin", "ping"}},
						}, InitialDelaySeconds: 30, TimeoutSeconds: 5,
					},
				},
			},
		},
		{
			name: "health and readinessProbe given together",
			input: Container{
				Health: &api_v1.Probe{
					Handler: api_v1.Handler{
						Exec: &api_v1.ExecAction{Command: []string{"mysqladmin", "ping"}},
					}, InitialDelaySeconds: 30, TimeoutSeconds: 5,
				},
				Container: api_v1.Container{
					ReadinessProbe: &api_v1.Probe{
						Handler: api_v1.Handler{
							Exec: &api_v1.ExecAction{Command: []string{"mysqladmin", "ping"}},
						}, InitialDelaySeconds: 30, TimeoutSeconds: 5,
					},
				},
			},
		},
		{
			name: "health and livenessProbe and readinessProbe given together",
			input: Container{
				Health: &api_v1.Probe{
					Handler: api_v1.Handler{
						Exec: &api_v1.ExecAction{Command: []string{"mysqladmin", "ping"}},
					}, InitialDelaySeconds: 30, TimeoutSeconds: 5,
				},
				Container: api_v1.Container{
					LivenessProbe: &api_v1.Probe{
						Handler: api_v1.Handler{
							Exec: &api_v1.ExecAction{Command: []string{"mysqladmin", "ping"}},
						}, InitialDelaySeconds: 30, TimeoutSeconds: 5,
					},
					ReadinessProbe: &api_v1.Probe{
						Handler: api_v1.Handler{
							Exec: &api_v1.ExecAction{Command: []string{"mysqladmin", "ping"}},
						}, InitialDelaySeconds: 30, TimeoutSeconds: 5,
					},
				},
			},
		},
	}

	for _, test := range failingTests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			if _, err := populateProbes(test.input); err == nil {
				t.Fatalf("expected failure but passed for input: %s",
					prettyPrintObjects(test.input))
			} else {
				t.Logf("failed with error: %v", err)
			}
		})
	}

	t.Logf("Running passing tests")
	passingTests := []struct {
		name   string
		input  Container
		output Container
	}{
		{
			name: "valid health given",
			input: Container{
				Health: &api_v1.Probe{
					Handler: api_v1.Handler{
						Exec: &api_v1.ExecAction{Command: []string{"mysqladmin", "ping"}},
					}, InitialDelaySeconds: 30, TimeoutSeconds: 5,
				},
			},
			output: Container{
				Container: api_v1.Container{
					LivenessProbe: &api_v1.Probe{
						Handler: api_v1.Handler{
							Exec: &api_v1.ExecAction{Command: []string{"mysqladmin", "ping"}},
						}, InitialDelaySeconds: 30, TimeoutSeconds: 5,
					},
					ReadinessProbe: &api_v1.Probe{
						Handler: api_v1.Handler{
							Exec: &api_v1.ExecAction{Command: []string{"mysqladmin", "ping"}},
						}, InitialDelaySeconds: 30, TimeoutSeconds: 5,
					},
				},
			},
		},
		{
			name: "only livenessProbe given",
			input: Container{
				Container: api_v1.Container{
					LivenessProbe: &api_v1.Probe{
						Handler: api_v1.Handler{
							Exec: &api_v1.ExecAction{Command: []string{"mysqladmin", "ping"}},
						}, InitialDelaySeconds: 30, TimeoutSeconds: 5,
					},
				},
			},
			output: Container{
				Container: api_v1.Container{
					LivenessProbe: &api_v1.Probe{
						Handler: api_v1.Handler{
							Exec: &api_v1.ExecAction{Command: []string{"mysqladmin", "ping"}},
						}, InitialDelaySeconds: 30, TimeoutSeconds: 5,
					},
				},
			},
		},
		{
			name: "only readinessProbe given",
			input: Container{
				Container: api_v1.Container{
					ReadinessProbe: &api_v1.Probe{
						Handler: api_v1.Handler{
							Exec: &api_v1.ExecAction{Command: []string{"mysqladmin", "ping"}},
						}, InitialDelaySeconds: 30, TimeoutSeconds: 5,
					},
				},
			},
			output: Container{
				Container: api_v1.Container{
					ReadinessProbe: &api_v1.Probe{
						Handler: api_v1.Handler{
							Exec: &api_v1.ExecAction{Command: []string{"mysqladmin", "ping"}},
						}, InitialDelaySeconds: 30, TimeoutSeconds: 5,
					},
				},
			},
		},
		{
			name:   "nothing given",
			input:  Container{},
			output: Container{},
		},
	}

	for _, test := range passingTests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got, err := populateProbes(test.input)
			if err != nil {
				t.Fatalf("got error %v, for input: %s",
					err, prettyPrintObjects(test.input))
			}

			if !reflect.DeepEqual(got, test.output) {
				t.Fatalf("expected: %s, got: %s", prettyPrintObjects(test.output),
					prettyPrintObjects(got))
			}
		})
	}
}

func TestPopulateServicePortNames(t *testing.T) {
	serviceName := "batman"
	tests := []struct {
		name               string
		inputServicePorts  []api_v1.ServicePort
		outputServicePorts []api_v1.ServicePort
	}{
		{
			name: "Passing only one servicePort, no population should happen",
			inputServicePorts: []api_v1.ServicePort{
				{
					Port: 8080,
				},
			},
			outputServicePorts: []api_v1.ServicePort{
				{
					Name: fmt.Sprintf("%v-%v", serviceName, 8080),
					Port: 8080,
				},
			},
		},
		{
			name: "Passing multiple servicePorts with no names, population should happen",
			inputServicePorts: []api_v1.ServicePort{
				{
					Port: 8080,
				},
				{
					Port: 8081,
				},
			},
			outputServicePorts: []api_v1.ServicePort{
				{
					Name: fmt.Sprintf("%v-%v", serviceName, 8080),
					Port: 8080,
				},
				{
					Name: fmt.Sprintf("%v-%v", serviceName, 8081),
					Port: 8081,
				},
			},
		},
		{
			name: "Passing multiple servicePorts with names, no population should happen",
			inputServicePorts: []api_v1.ServicePort{
				{
					Name: fmt.Sprintf("%v-%v", "prepopulated", 8080),
					Port: 8080,
				},
				{
					Name: fmt.Sprintf("%v-%v", "prepopulated", 8081),
					Port: 8081,
				},
			},
			outputServicePorts: []api_v1.ServicePort{
				{
					Name: fmt.Sprintf("%v-%v", "prepopulated", 8080),
					Port: 8080,
				},
				{
					Name: fmt.Sprintf("%v-%v", "prepopulated", 8081),
					Port: 8081,
				},
			},
		},
		{
			name: "Passing multiple servicePorts, some with names, some without, selective population should happen",
			inputServicePorts: []api_v1.ServicePort{
				{
					Name: fmt.Sprintf("%v-%v", "prepopulated", 8080),
					Port: 8080,
				},
				{
					Name: fmt.Sprintf("%v-%v", "prepopulated", 8081),
					Port: 8081,
				},
				{
					Port: 8082,
				},
				{
					Port: 8083,
				},
			},
			outputServicePorts: []api_v1.ServicePort{
				{
					Name: fmt.Sprintf("%v-%v", "prepopulated", 8080),
					Port: 8080,
				},
				{
					Name: fmt.Sprintf("%v-%v", "prepopulated", 8081),
					Port: 8081,
				},
				{
					Name: fmt.Sprintf("%v-%v", serviceName, 8082),
					Port: 8082,
				},
				{
					Name: fmt.Sprintf("%v-%v", serviceName, 8083),
					Port: 8083,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			populateServicePortNames(serviceName, test.inputServicePorts)

			if !reflect.DeepEqual(test.inputServicePorts, test.outputServicePorts) {
				t.Errorf("For input\n%v\nExpected output to be\n%v", spew.Sprint(test.inputServicePorts), spew.Sprint(test.outputServicePorts))
			}
		})
	}
}

func TestPopulateVolumes(t *testing.T) {
	volumeClaims := []VolumeClaim{
		{
			ObjectMeta: meta_v1.ObjectMeta{
				Name: "foo",
			},
		},
		{
			ObjectMeta: meta_v1.ObjectMeta{
				Name: "bar",
			},
		},
		{
			ObjectMeta: meta_v1.ObjectMeta{
				Name: "barfoo",
			},
		},
	}

	volumes := []api_v1.Volume{{Name: "foo"}}

	// a volumeMount is defined but that is not there in volumeClaims
	// neither it is in pod level volumes, so this should fail
	failingContainers := []api_v1.Container{
		{VolumeMounts: []api_v1.VolumeMount{{Name: "baz"}}},
	}

	if _, err := populateVolumes(failingContainers, volumeClaims, volumes); err == nil {
		t.Errorf("should have failed but passed for volumeMount that" +
			" does not exist.")
	}

	passingContainers := []api_v1.Container{
		{VolumeMounts: []api_v1.VolumeMount{{Name: "bar"}}},
		{VolumeMounts: []api_v1.VolumeMount{{Name: "barfoo"}}},
	}
	expected := []api_v1.Volume{
		{
			Name: "bar",
			VolumeSource: api_v1.VolumeSource{
				PersistentVolumeClaim: &api_v1.PersistentVolumeClaimVolumeSource{
					ClaimName: "bar",
				},
			},
		},
		{
			Name: "barfoo",
			VolumeSource: api_v1.VolumeSource{
				PersistentVolumeClaim: &api_v1.PersistentVolumeClaimVolumeSource{
					ClaimName: "barfoo",
				},
			},
		},
	}

	newVols, err := populateVolumes(passingContainers, volumeClaims, volumes)
	if err != nil {
		t.Fatalf("test failed: %v", err)
	}
	if !reflect.DeepEqual(newVols, expected) {
		t.Fatalf("expected: %s, got: %s", prettyPrintObjects(expected), prettyPrintObjects(newVols))
	}
}
