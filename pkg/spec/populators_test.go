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
	"sort"
	"testing"

	"fmt"
	"github.com/davecgh/go-spew/spew"
	"k8s.io/apimachinery/pkg/util/intstr"
	api_v1 "k8s.io/client-go/pkg/api/v1"
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

func TestConvertMapToList(t *testing.T) {

	tests := []struct {
		input  map[string]string
		output []string
	}{
		{
			input:  map[string]string{"one": "", "two": "", "three": ""},
			output: []string{"one", "three", "two"},
		},
		{
			input:  map[string]string{"deployment": "", "application": "", "configMap": "", "containers": ""},
			output: []string{"application", "configMap", "containers", "deployment"},
		},
	}

	for _, test := range tests {
		got := getMapKeys(test.input)
		sort.Strings(got)
		if !reflect.DeepEqual(test.output, got) {
			t.Errorf("expected: %+v got: %+v", test.output, got)
		}
	}
}

var cms = []ConfigMapMod{
	{
		Name: "test1", Data: map[string]string{"ten": "TEN"},
	},
	{
		Name: "test2",
		Data: map[string]string{"two": "TWO", "four": "FOUR", "eight": "EIGHT"},
	},
}
var secrets = []SecretMod{
	{
		Name: "test1",
		Secret: api_v1.Secret{
			Data:       map[string][]byte{"one": []byte("ONE"), "five": []byte("FIVE")},
			StringData: map[string]string{"three": "THREE", "four": "FOUR"},
		},
	},
	{
		Name: "test2",
		Secret: api_v1.Secret{
			Data:       map[string][]byte{"one": []byte("ONE"), "two": []byte("TWO")},
			StringData: map[string]string{"three": "THREE", "four": "FOUR"},
		},
	},
	{
		Name: "test3",
	},
}

func TestSearchConfigMap(t *testing.T) {
	t.Run("running passing test", func(t *testing.T) {
		t.Parallel()
		_, err := searchConfigMap(cms, "test2")
		if err != nil {
			t.Fatalf("%v", err)
		}
	})

	t.Run("running failing test", func(t *testing.T) {
		t.Parallel()
		_, err := searchConfigMap(cms, "test5")
		if err == nil {
			t.Fatalf("should have failed but passed")
		} else {
			t.Logf("failed with error: %v", err)
		}
	})
}

func TestGetSecretDataKeys(t *testing.T) {
	t.Run("running passing test", func(t *testing.T) {
		t.Parallel()
		// TODO: also check the returned keys the problem with it
		// is that it is list which is generated from a map
		// so order will change everytime so using reflect.DeepEqual won't
		// be helpful
		_, err := getSecretDataKeys(secrets, "test2")
		if err != nil {
			t.Fatalf("%v", err)
		}
	})

	t.Run("running failing test", func(t *testing.T) {
		t.Parallel()
		_, err := getSecretDataKeys(secrets, "test5")
		if err == nil {
			t.Fatalf("should have failed but passed")
		} else {
			t.Logf("failed with error: %v", err)
		}
	})

}

func TestConvertEnvFromToEnvs(t *testing.T) {

	// pass configmap and see all the envs
	// pass secret and see all the envs
	// pass configmap and secret and see all the envs
	t.Logf("Running passing tests")
	passingTests := []struct {
		name   string
		input  []api_v1.EnvFromSource
		output []api_v1.EnvVar
	}{
		{
			name: "all envs from configMap",
			input: []api_v1.EnvFromSource{
				{
					ConfigMapRef: &api_v1.ConfigMapEnvSource{
						LocalObjectReference: api_v1.LocalObjectReference{Name: "test2"},
					},
				},
			},
			output: []api_v1.EnvVar{
				{
					Name: "eight",
					ValueFrom: &api_v1.EnvVarSource{
						ConfigMapKeyRef: &api_v1.ConfigMapKeySelector{
							LocalObjectReference: api_v1.LocalObjectReference{Name: "test2"},
							Key:                  "eight",
						},
					},
				},
				{
					Name: "four",
					ValueFrom: &api_v1.EnvVarSource{
						ConfigMapKeyRef: &api_v1.ConfigMapKeySelector{
							LocalObjectReference: api_v1.LocalObjectReference{Name: "test2"},
							Key:                  "four",
						},
					},
				},
				{
					Name: "two",
					ValueFrom: &api_v1.EnvVarSource{
						ConfigMapKeyRef: &api_v1.ConfigMapKeySelector{
							LocalObjectReference: api_v1.LocalObjectReference{Name: "test2"},
							Key:                  "two",
						},
					},
				},
			},
		},
		{
			name: "all envs from secret",
			input: []api_v1.EnvFromSource{
				{
					SecretRef: &api_v1.SecretEnvSource{
						LocalObjectReference: api_v1.LocalObjectReference{Name: "test2"},
					},
				},
			},
			output: []api_v1.EnvVar{
				{
					Name: "four",
					ValueFrom: &api_v1.EnvVarSource{
						SecretKeyRef: &api_v1.SecretKeySelector{
							LocalObjectReference: api_v1.LocalObjectReference{Name: "test2"},
							Key:                  "four",
						},
					},
				},
				{
					Name: "one",
					ValueFrom: &api_v1.EnvVarSource{
						SecretKeyRef: &api_v1.SecretKeySelector{
							LocalObjectReference: api_v1.LocalObjectReference{Name: "test2"},
							Key:                  "one",
						},
					},
				},
				{
					Name: "three",
					ValueFrom: &api_v1.EnvVarSource{
						SecretKeyRef: &api_v1.SecretKeySelector{
							LocalObjectReference: api_v1.LocalObjectReference{Name: "test2"},
							Key:                  "three",
						},
					},
				},
				{
					Name: "two",
					ValueFrom: &api_v1.EnvVarSource{
						SecretKeyRef: &api_v1.SecretKeySelector{
							LocalObjectReference: api_v1.LocalObjectReference{Name: "test2"},
							Key:                  "two",
						},
					},
				},
			},
		},
		{
			name: "envs from both secret and configmap",
			input: []api_v1.EnvFromSource{
				{
					SecretRef: &api_v1.SecretEnvSource{
						LocalObjectReference: api_v1.LocalObjectReference{Name: "test1"},
					},
					ConfigMapRef: &api_v1.ConfigMapEnvSource{
						LocalObjectReference: api_v1.LocalObjectReference{Name: "test1"},
					},
				},
			},
			output: []api_v1.EnvVar{
				{
					Name: "ten",
					ValueFrom: &api_v1.EnvVarSource{
						ConfigMapKeyRef: &api_v1.ConfigMapKeySelector{
							LocalObjectReference: api_v1.LocalObjectReference{Name: "test1"},
							Key:                  "ten",
						},
					},
				},
				{
					Name: "five",
					ValueFrom: &api_v1.EnvVarSource{
						SecretKeyRef: &api_v1.SecretKeySelector{
							LocalObjectReference: api_v1.LocalObjectReference{Name: "test1"},
							Key:                  "five",
						},
					},
				},
				{
					Name: "four",
					ValueFrom: &api_v1.EnvVarSource{
						SecretKeyRef: &api_v1.SecretKeySelector{
							LocalObjectReference: api_v1.LocalObjectReference{Name: "test1"},
							Key:                  "four",
						},
					},
				},
				{
					Name: "one",
					ValueFrom: &api_v1.EnvVarSource{
						SecretKeyRef: &api_v1.SecretKeySelector{
							LocalObjectReference: api_v1.LocalObjectReference{Name: "test1"},
							Key:                  "one",
						},
					},
				},
				{
					Name: "three",
					ValueFrom: &api_v1.EnvVarSource{
						SecretKeyRef: &api_v1.SecretKeySelector{
							LocalObjectReference: api_v1.LocalObjectReference{Name: "test1"},
							Key:                  "three",
						},
					},
				},
			},
		},
	}

	for _, test := range passingTests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got, err := convertEnvFromToEnvs(test.input, cms, secrets)
			if err != nil {
				t.Fatalf("errored out: %v", err)
			}
			if !reflect.DeepEqual(test.output, got) {
				t.Fatalf("expected: %s, got: %s",
					prettyPrintObjects(test.output), prettyPrintObjects(got))
			}
		})
	}

	t.Logf("Running failing tests")
	failingTests := []struct {
		name  string
		input []api_v1.EnvFromSource
	}{
		{
			name: "configMap that does not exists",
			input: []api_v1.EnvFromSource{
				{
					ConfigMapRef: &api_v1.ConfigMapEnvSource{
						LocalObjectReference: api_v1.LocalObjectReference{Name: "test9"},
					},
				},
			},
		},
		{
			name: "secret that does not exists",
			input: []api_v1.EnvFromSource{
				{
					SecretRef: &api_v1.SecretEnvSource{
						LocalObjectReference: api_v1.LocalObjectReference{Name: "test5"},
					},
				},
			},
		},
	}

	for _, test := range failingTests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			_, err := convertEnvFromToEnvs(test.input, cms, secrets)
			if err == nil {
				t.Fatalf("expected failure but passed for input: %s", prettyPrintObjects(test.input))
			} else {
				t.Logf("errored out with: %v", err)
			}
		})
	}
}

func TestPopulateEnvFrom(t *testing.T) {
	tests := []struct {
		name   string
		input  Container
		output Container
	}{
		{
			name: "normal envFrom",
			input: Container{
				Container: api_v1.Container{
					EnvFrom: []api_v1.EnvFromSource{
						{
							ConfigMapRef: &api_v1.ConfigMapEnvSource{
								LocalObjectReference: api_v1.LocalObjectReference{Name: "test1"},
							},
						},
					},
				},
			},
			output: Container{
				Container: api_v1.Container{
					Env: []api_v1.EnvVar{
						{
							Name: "ten",
							ValueFrom: &api_v1.EnvVarSource{
								ConfigMapKeyRef: &api_v1.ConfigMapKeySelector{
									LocalObjectReference: api_v1.LocalObjectReference{Name: "test1"},
									Key:                  "ten",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "container that has envs already",
			input: Container{
				Container: api_v1.Container{
					EnvFrom: []api_v1.EnvFromSource{
						{
							ConfigMapRef: &api_v1.ConfigMapEnvSource{
								LocalObjectReference: api_v1.LocalObjectReference{Name: "test1"},
							},
						},
					},
					Env: []api_v1.EnvVar{
						{Name: "data", Value: "data"},
						{Name: "ten", Value: "ten"},
					},
				},
			},
			output: Container{
				Container: api_v1.Container{
					Env: []api_v1.EnvVar{
						{
							Name: "ten",
							ValueFrom: &api_v1.EnvVarSource{
								ConfigMapKeyRef: &api_v1.ConfigMapKeySelector{
									LocalObjectReference: api_v1.LocalObjectReference{Name: "test1"},
									Key:                  "ten",
								},
							},
						},
						{Name: "data", Value: "data"},
						{Name: "ten", Value: "ten"},
					},
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got, err := populateEnvFrom(test.input, cms, secrets)
			if err != nil {
				t.Fatalf("errored: %v", err)
			}
			if !reflect.DeepEqual(test.output, got) {
				t.Fatalf("expected: %s, got %s",
					prettyPrintObjects(test.output), prettyPrintObjects(got))
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

func TestParsePortMapping(t *testing.T) {
	tests := []struct {
		name        string
		portMapping string
		servicePort *api_v1.ServicePort
		success     bool
	}{
		{
			name:        "Nothing is passed, not even port",
			portMapping: "",
			servicePort: nil,
			success:     false,
		},
		{
			name:        "Only 'port' is passed",
			portMapping: "1337",
			servicePort: &api_v1.ServicePort{
				Port: 1337,
				TargetPort: intstr.IntOrString{
					IntVal: 1337,
				},
				Protocol: api_v1.ProtocolTCP,
			},
			success: true,
		},
		{
			name:        "port:targetPort is passed",
			portMapping: "1337:1338",
			servicePort: &api_v1.ServicePort{
				Port: 1337,
				TargetPort: intstr.IntOrString{
					IntVal: 1338,
				},
				Protocol: api_v1.ProtocolTCP,
			},
			success: true,
		},
		{
			name:        "port/protocol is passed",
			portMapping: "1337/UDP",
			servicePort: &api_v1.ServicePort{
				Port: 1337,
				TargetPort: intstr.IntOrString{
					IntVal: 1337,
				},
				Protocol: api_v1.ProtocolUDP,
			},
			success: true,
		},
		{
			name:        "port:targetPort/protocol is passed",
			portMapping: "1337:1338/UDP",
			servicePort: &api_v1.ServicePort{
				Port: 1337,
				TargetPort: intstr.IntOrString{
					IntVal: 1338,
				},
				Protocol: api_v1.ProtocolUDP,
			},
			success: true,
		},
		{
			name:        "Invalid protocol (neither TCP nor UDP) is passed",
			portMapping: "1337:1338/INVALID",
			servicePort: nil,
			success:     false,
		},
		{
			name:        "Multiple protocols passed, multiple '/' test",
			portMapping: "1337/TCP:1338/TCP",
			servicePort: nil,
			success:     false,
		},
		{
			name:        "Non int port is passed",
			portMapping: "batman:1338/TCP",
			servicePort: nil,
			success:     false,
		},
		{
			name:        "Non int targetPort is passed",
			portMapping: "1337:batman/TCP",
			servicePort: nil,
			success:     false,
		},
		{
			name:        "More than 2 ports passed, multiple ':' test",
			portMapping: "1337:1338:1339/TCP",
			servicePort: nil,
			success:     false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			sp, err := parsePortMapping(test.portMapping)

			switch test.success {
			case true:
				if err != nil {
					t.Errorf("Expected test to pass but got an error -\n%v", err)
				}
			case false:
				if err == nil {
					t.Errorf("Expected %v to fail, but test passed!", test.portMapping)
				}
			}

			if !reflect.DeepEqual(sp, test.servicePort) {
				t.Errorf("Expected ServicePort to be -\n%v\nBut got -\n%v", spew.Sprint(test.servicePort), spew.Sprint(sp))
			}
		})
	}
}

func TestPopulateVolumes(t *testing.T) {
	volumeClaims := []VolumeClaim{{Name: "foo"}, {Name: "bar"}, {Name: "barfoo"}}
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
