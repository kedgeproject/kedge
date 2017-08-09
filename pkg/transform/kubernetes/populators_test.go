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

package kubernetes

import (
	"encoding/json"
	"reflect"
	"sort"
	"testing"

	"github.com/kedgeproject/kedge/pkg/spec"

	api_v1 "k8s.io/client-go/pkg/api/v1"
)

func TestPopulateProbes(t *testing.T) {
	t.Logf("Running failing tests")
	failingTests := []struct {
		name  string
		input spec.Container
	}{
		{
			name: "health and livenessProbe given together",
			input: spec.Container{
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
			input: spec.Container{
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
			input: spec.Container{
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
		input  spec.Container
		output spec.Container
	}{
		{
			name: "valid health given",
			input: spec.Container{
				Health: &api_v1.Probe{
					Handler: api_v1.Handler{
						Exec: &api_v1.ExecAction{Command: []string{"mysqladmin", "ping"}},
					}, InitialDelaySeconds: 30, TimeoutSeconds: 5,
				},
			},
			output: spec.Container{
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
			input: spec.Container{
				Container: api_v1.Container{
					LivenessProbe: &api_v1.Probe{
						Handler: api_v1.Handler{
							Exec: &api_v1.ExecAction{Command: []string{"mysqladmin", "ping"}},
						}, InitialDelaySeconds: 30, TimeoutSeconds: 5,
					},
				},
			},
			output: spec.Container{
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
			input: spec.Container{
				Container: api_v1.Container{
					ReadinessProbe: &api_v1.Probe{
						Handler: api_v1.Handler{
							Exec: &api_v1.ExecAction{Command: []string{"mysqladmin", "ping"}},
						}, InitialDelaySeconds: 30, TimeoutSeconds: 5,
					},
				},
			},
			output: spec.Container{
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
			input:  spec.Container{},
			output: spec.Container{},
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

var cms = []spec.ConfigMapMod{
	{
		Name: "test1", Data: map[string]string{"ten": "TEN"},
	},
	{
		Name: "test2",
		Data: map[string]string{"two": "TWO", "four": "FOUR", "eight": "EIGHT"},
	},
}
var secrets = []spec.SecretMod{
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
		input  spec.Container
		output spec.Container
	}{
		{
			name: "normal envFrom",
			input: spec.Container{
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
			output: spec.Container{
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
			input: spec.Container{
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
			output: spec.Container{
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

func prettyPrintObjects(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}
