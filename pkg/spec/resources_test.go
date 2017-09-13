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

	"k8s.io/apimachinery/pkg/runtime"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	api_v1 "k8s.io/client-go/pkg/api/v1"
)

func TestFixServices(t *testing.T) {
	failingTest := []ServiceSpecMod{
		{Ports: nil},
		{Ports: nil},
		{Ports: nil},
	}
	_, err := fixServices(failingTest, "")
	if err == nil {
		t.Errorf("should have failed but passed")
	} else {
		t.Logf("failed with error: %v", err)
	}

	appName := "test"
	passingTests := []struct {
		Name   string
		Input  []ServiceSpecMod
		Output []ServiceSpecMod
	}{
		{
			Name:   "only one service given",
			Input:  []ServiceSpecMod{{}},
			Output: []ServiceSpecMod{{Name: appName}},
		},
		{
			Name: "multiple ports and no port name given",
			Input: []ServiceSpecMod{
				{
					Ports: []ServicePortMod{
						{ServicePort: api_v1.ServicePort{Port: 8080}},
						{ServicePort: api_v1.ServicePort{Port: 8081}},
					},
				},
			},
			Output: []ServiceSpecMod{
				{
					Name: appName,
					Ports: []ServicePortMod{
						{
							ServicePort: api_v1.ServicePort{
								Name: appName + "-8080", Port: 8080,
							},
						},
						{
							ServicePort: api_v1.ServicePort{
								Name: appName + "-8081", Port: 8081,
							},
						},
					},
				},
			},
		},
	}

	for _, test := range passingTests {
		t.Logf("Running test: %s", test.Name)
		got, err := fixServices(test.Input, appName)
		if err != nil {
			t.Errorf("expected to pass but failed with: %v", err)
		}
		if !reflect.DeepEqual(got, test.Output) {
			t.Errorf("expected: %s, got: %s", prettyPrintObjects(test.Output),
				prettyPrintObjects(got))
		}
	}
}

func TestFixVolumeClaims(t *testing.T) {
	failingTest := []VolumeClaim{{}, {}}

	_, err := fixVolumeClaims(failingTest, "")
	if err == nil {
		t.Errorf("should have failed but passed")
	} else {
		t.Logf("failed with error: %v", err)
	}

	appName := "test"
	passingTest := []VolumeClaim{{}}
	expected := []VolumeClaim{{Name: appName}}
	got, err := fixVolumeClaims(passingTest, appName)
	if err != nil {
		t.Errorf("expected to pass but failed with: %v", err)
	}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected: %s, got: %s", prettyPrintObjects(expected),
			prettyPrintObjects(got))
	}
}

func TestFixConfigMaps(t *testing.T) {
	failingTest := []ConfigMapMod{{}, {}}
	_, err := fixConfigMaps(failingTest, "")
	if err == nil {
		t.Errorf("should have failed but passed")
	} else {
		t.Logf("failed with error: %v", err)
	}

	appName := "test"
	passingTest := []ConfigMapMod{{}}
	expected := []ConfigMapMod{{Name: appName}}
	got, err := fixConfigMaps(passingTest, appName)
	if err != nil {
		t.Errorf("expected to pass but failed with: %v", err)
	}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected: %s, got: %s", prettyPrintObjects(expected),
			prettyPrintObjects(got))
	}
}

func TestFixSecrets(t *testing.T) {
	failingTest := []SecretMod{{}, {}}
	_, err := fixSecrets(failingTest, "")
	if err == nil {
		t.Errorf("should have failed but passed")
	} else {
		t.Logf("failed with error: %v", err)
	}

	appName := "test"
	passingTest := []SecretMod{{}}
	expected := []SecretMod{{Name: appName}}
	got, err := fixSecrets(passingTest, appName)
	if err != nil {
		t.Errorf("expected to pass but failed with: %v", err)
	}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected: %s, got: %s", prettyPrintObjects(expected),
			prettyPrintObjects(got))
	}
}

func TestFixContainers(t *testing.T) {
	failingTest := []Container{{}, {}}
	_, err := fixContainers(failingTest, "")
	if err == nil {
		t.Errorf("should have failed but passed")
	} else {
		t.Logf("failed with error: %v", err)
	}

	appName := "test"
	passingTest := []Container{{}}
	expected := []Container{
		{Container: api_v1.Container{Name: appName}},
	}
	got, err := fixContainers(passingTest, appName)
	if err != nil {
		t.Errorf("expected to pass but failed with: %v", err)
	}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected: %s, got: %s", prettyPrintObjects(expected),
			prettyPrintObjects(got))
	}
}

func TestValidateVolumeClaims(t *testing.T) {

	failingTest := []VolumeClaim{{Name: "foo"}, {Name: "bar"}, {Name: "foo"}}

	err := validateVolumeClaims(failingTest)
	if err == nil {
		t.Errorf("should have failed but passed for input: %+v", failingTest)
	} else {
		t.Logf("failed with error: %v", err)
	}

}

func TestCreateServices(t *testing.T) {
	tests := []struct {
		Name    string
		App     *DeploymentSpecMod
		Objects []runtime.Object
	}{
		{
			"Single container specified",
			&DeploymentSpecMod{
				ControllerFields: ControllerFields{

					Name: "test",
					PodSpecMod: PodSpecMod{
						Containers: []Container{{Container: api_v1.Container{Image: "nginx"}}},
					},
					Services: []ServiceSpecMod{
						{Name: "test", Ports: []ServicePortMod{{ServicePort: api_v1.ServicePort{Port: 8080}}}},
					},
				},
			},
			append(make([]runtime.Object, 0), &api_v1.Service{
				ObjectMeta: metav1.ObjectMeta{Name: "test"},
				Spec:       api_v1.ServiceSpec{Ports: []api_v1.ServicePort{{Port: 8080}}},
			}),
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			object, err := test.App.createServices()
			if err != nil {
				t.Fatalf("Creating services failed: %v", err)
			}
			if !reflect.DeepEqual(test.Objects, object) {
				t.Fatalf("Expected:\n%v\nGot:\n%v", test.Objects, object)
			}
		})
	}
}

// TODO: add test for auto naming of single persistent volume
