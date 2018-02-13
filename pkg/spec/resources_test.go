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
	build_v1 "github.com/openshift/origin/pkg/build/apis/build/v1"
	image_v1 "github.com/openshift/origin/pkg/image/apis/image/v1"
	os_route_v1 "github.com/openshift/origin/pkg/route/apis/route/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	api_v1 "k8s.io/kubernetes/pkg/api/v1"
	ext_v1beta1 "k8s.io/kubernetes/pkg/apis/extensions/v1beta1"
)

func TestFixServices(t *testing.T) {
	passingTests := []struct {
		Name    string
		Input   *App
		Output  *App
		Success bool
	}{
		{
			Name: "only one service given",
			Input: &App{
				ObjectMeta: meta_v1.ObjectMeta{
					Name: "test",
				},
				Services: []ServiceSpecMod{
					{
						Ports: []ServicePortMod{{ServicePort: api_v1.ServicePort{Port: 8080}}},
					},
				},
			},
			Output: &App{
				ObjectMeta: meta_v1.ObjectMeta{
					Name: "test",
				},
				Services: []ServiceSpecMod{
					{
						ObjectMeta: meta_v1.ObjectMeta{
							Name: "test",
						},
						Ports: []ServicePortMod{{ServicePort: api_v1.ServicePort{Port: 8080}}}},
				},
			},
			Success: true,
		},
		{
			Name: "Global labels specified",
			Input: &App{
				ObjectMeta: meta_v1.ObjectMeta{
					Name: "test",
					Labels: map[string]string{
						"foo": "bar",
					},
				},
				Services: []ServiceSpecMod{
					{
						ObjectMeta: meta_v1.ObjectMeta{
							Name: "test",
						},
						Ports: []ServicePortMod{{ServicePort: api_v1.ServicePort{Port: 8080}}}},
				},
			},
			Output: &App{
				ObjectMeta: meta_v1.ObjectMeta{
					Name: "test",
					Labels: map[string]string{
						"foo": "bar",
					},
				},
				Services: []ServiceSpecMod{
					{
						ObjectMeta: meta_v1.ObjectMeta{
							Name: "test",
							Labels: map[string]string{
								"foo": "bar",
							},
						},
						Ports: []ServicePortMod{{ServicePort: api_v1.ServicePort{Port: 8080}}}},
				},
			},
			Success: true,
		},
	}

	for _, test := range passingTests {
		t.Logf("Running test: %s", test.Name)
		err := test.Input.fixServices()
		if test.Success {
			if err != nil {
				t.Errorf("expected to pass but failed with: %v", err)
			}
			if !reflect.DeepEqual(test.Input, test.Output) {
				t.Errorf("expected: \n%s\n\n got:\n %s\n\n", PrettyPrintObjects(test.Output),
					PrettyPrintObjects(test.Input))
			}
		} else {
			if err == nil {
				t.Errorf("test was expected to fail, but it passed")
			}

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
	expected := []VolumeClaim{
		{
			ObjectMeta: meta_v1.ObjectMeta{
				Name: appName,
				Labels: map[string]string{
					appLabelKey: appName,
				},
			},
		},
	}
	got, err := fixVolumeClaims(passingTest, appName)
	if err != nil {
		t.Errorf("expected to pass but failed with: %v", err)
	}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected: %s, got: %s", PrettyPrintObjects(expected),
			PrettyPrintObjects(got))
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
	expected := []ConfigMapMod{
		{
			ObjectMeta: meta_v1.ObjectMeta{
				Name: appName,
				Labels: map[string]string{
					appLabelKey: appName,
				},
			},
		}}
	got, err := fixConfigMaps(passingTest, appName)
	if err != nil {
		t.Errorf("expected to pass but failed with: %v", err)
	}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected: %s, got: %s", PrettyPrintObjects(expected),
			PrettyPrintObjects(got))
	}
}

func TestFixBuildConfigs(t *testing.T) {
	appName := "testAppName"
	tests := []struct {
		name    string
		input   []BuildConfigSpecMod
		output  []BuildConfigSpecMod
		success bool
	}{
		{
			name: "passing one BuildConfig without name",
			input: []BuildConfigSpecMod{
				{},
			},
			output: []BuildConfigSpecMod{
				{
					ObjectMeta: meta_v1.ObjectMeta{
						Name: appName,
						Labels: map[string]string{
							appLabelKey: appName,
						},
					},
				},
			},
			success: true,
		},
		{
			name: "passing one BuildConfig with name",
			input: []BuildConfigSpecMod{
				{
					ObjectMeta: meta_v1.ObjectMeta{
						Name: "bcName",
					},
				},
			},
			output: []BuildConfigSpecMod{
				{
					ObjectMeta: meta_v1.ObjectMeta{
						Name: "bcName",
						Labels: map[string]string{
							appLabelKey: appName,
						},
					},
				},
			},
			success: true,
		},
		{
			name: "passing multiple BuildConfigs without names",
			input: []BuildConfigSpecMod{
				{},
				{},
			},
			output:  nil,
			success: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fixedBuildConfigs, err := fixBuildConfigs(test.input, appName)

			switch test.success {
			case true:
				if err != nil {
					t.Errorf("Expected test to pass but got an error -\n%v", err)
				}
			case false:
				if err == nil {
					t.Errorf("For the input -\n%v\nexpected test to fail, but test passed", PrettyPrintObjects(test.input))
				}
			}

			if !reflect.DeepEqual(fixedBuildConfigs, test.output) {
				t.Errorf("Expected fixed BuildConfigs to be -\n%v\nBut got -\n%v\n", PrettyPrintObjects(test.output), PrettyPrintObjects(fixedBuildConfigs))
			}
		})
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
	expected := []SecretMod{
		{
			ObjectMeta: meta_v1.ObjectMeta{
				Name: appName,
				Labels: map[string]string{
					appLabelKey: appName,
				},
			},
		},
	}
	got, err := fixSecrets(passingTest, appName)
	if err != nil {
		t.Errorf("expected to pass but failed with: %v", err)
	}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected: %s, got: %s", PrettyPrintObjects(expected),
			PrettyPrintObjects(got))
	}
}

func TestFixImageStreams(t *testing.T) {
	appName := "testAppName"
	tests := []struct {
		name    string
		input   []ImageStreamSpecMod
		output  []ImageStreamSpecMod
		success bool
	}{
		{
			name: "passing one imageStream without name",
			input: []ImageStreamSpecMod{
				{},
			},
			output: []ImageStreamSpecMod{
				{
					ObjectMeta: meta_v1.ObjectMeta{
						Name: appName,
						Labels: map[string]string{
							appLabelKey: appName,
						},
					},
				},
			},
			success: true,
		},
		{
			name: "passing one imageStream with name",
			input: []ImageStreamSpecMod{
				{
					ObjectMeta: meta_v1.ObjectMeta{
						Name: "imageStreamName",
					},
				},
			},
			output: []ImageStreamSpecMod{
				{
					ObjectMeta: meta_v1.ObjectMeta{
						Name: "imageStreamName",
						Labels: map[string]string{
							appLabelKey: appName,
						},
					},
				},
			},
			success: true,
		},
		{
			name: "passing multiple ingresses without names",
			input: []ImageStreamSpecMod{
				{},
				{},
			},
			output:  nil,
			success: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fixedImageStreams, err := fixImageStreams(test.input, appName)

			switch test.success {
			case true:
				if err != nil {
					t.Errorf("Expected test to pass but got an error -\n%v", err)
				}
			case false:
				if err == nil {
					t.Errorf("For the input -\n%v\nexpected test to fail, but test passed", PrettyPrintObjects(test.input))
				}
			}

			if !reflect.DeepEqual(fixedImageStreams, test.output) {
				t.Errorf("Expected fixed imageStreams to be -\n%v\nBut got -\n%v\n", PrettyPrintObjects(test.output), PrettyPrintObjects(fixedImageStreams))
			}
		})
	}
}

func TestFixIngresses(t *testing.T) {
	appName := "testAppName"
	tests := []struct {
		name    string
		input   []IngressSpecMod
		output  []IngressSpecMod
		success bool
	}{
		{
			name: "passing one ingress without name",
			input: []IngressSpecMod{
				{
					IngressSpec: ext_v1beta1.IngressSpec{
						Rules: []ext_v1beta1.IngressRule{
							{
								Host: "testHost",
							},
						},
					},
				},
			},
			output: []IngressSpecMod{
				{
					ObjectMeta: meta_v1.ObjectMeta{
						Name: appName,
						Labels: map[string]string{
							appLabelKey: appName,
						},
					},
					IngressSpec: ext_v1beta1.IngressSpec{
						Rules: []ext_v1beta1.IngressRule{
							{
								Host: "testHost",
							},
						},
					},
				},
			},
			success: true,
		},
		{
			name: "passing one ingress with name",
			input: []IngressSpecMod{
				{
					ObjectMeta: meta_v1.ObjectMeta{
						Name: "ingressName",
					},
					IngressSpec: ext_v1beta1.IngressSpec{
						Rules: []ext_v1beta1.IngressRule{
							{
								Host: "testHost",
							},
						},
					},
				},
			},
			output: []IngressSpecMod{
				{
					ObjectMeta: meta_v1.ObjectMeta{
						Name: "ingressName",
						Labels: map[string]string{
							appLabelKey: appName,
						},
					},
					IngressSpec: ext_v1beta1.IngressSpec{
						Rules: []ext_v1beta1.IngressRule{
							{
								Host: "testHost",
							},
						},
					},
				},
			},
			success: true,
		},
		{
			name: "passing multiple ingresses without names",
			input: []IngressSpecMod{
				{
					IngressSpec: ext_v1beta1.IngressSpec{
						Rules: []ext_v1beta1.IngressRule{
							{
								Host: "testHost1",
							},
						},
					},
				},
				{
					IngressSpec: ext_v1beta1.IngressSpec{
						Rules: []ext_v1beta1.IngressRule{
							{
								Host: "testHost2",
							},
						},
					},
				},
			},
			output:  nil,
			success: false,
		},
		{
			name: "passing multiple ingresses",
			input: []IngressSpecMod{
				{
					ObjectMeta: meta_v1.ObjectMeta{
						Name: "ingress1",
					},
					IngressSpec: ext_v1beta1.IngressSpec{
						Rules: []ext_v1beta1.IngressRule{
							{
								Host: "testHost1",
							},
						},
					},
				},
				{
					ObjectMeta: meta_v1.ObjectMeta{
						Name: "ingress2",
					},
					IngressSpec: ext_v1beta1.IngressSpec{
						Rules: []ext_v1beta1.IngressRule{
							{
								Host: "testHost2",
							},
						},
					},
				},
			},
			output: []IngressSpecMod{
				{
					ObjectMeta: meta_v1.ObjectMeta{
						Name: "ingress1",
						Labels: map[string]string{
							appLabelKey: appName,
						},
					},
					IngressSpec: ext_v1beta1.IngressSpec{
						Rules: []ext_v1beta1.IngressRule{
							{
								Host: "testHost1",
							},
						},
					},
				},
				{
					ObjectMeta: meta_v1.ObjectMeta{
						Name: "ingress2",
						Labels: map[string]string{
							appLabelKey: appName,
						},
					},
					IngressSpec: ext_v1beta1.IngressSpec{
						Rules: []ext_v1beta1.IngressRule{
							{
								Host: "testHost2",
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
			fixedIngresses, err := fixIngresses(test.input, appName)

			switch test.success {
			case true:
				if err != nil {
					t.Errorf("Expected test to pass but got an error -\n%v", err)
				}
			case false:
				if err == nil {
					t.Errorf("For the input -\n%v\nexpected test to fail, but test passed", PrettyPrintObjects(test.input))
				}
			}

			if !reflect.DeepEqual(fixedIngresses, test.output) {
				t.Errorf("Expected fixed ingresses to be -\n%v\nBut got -\n%v\n", PrettyPrintObjects(test.output), PrettyPrintObjects(fixedIngresses))
			}
		})
	}

}

func TestFixRoutes(t *testing.T) {
	appName := "testAppName"
	tests := []struct {
		name    string
		input   []RouteSpecMod
		output  []RouteSpecMod
		success bool
	}{
		{
			name: "passing one route without name",
			input: []RouteSpecMod{
				{},
			},
			output: []RouteSpecMod{
				{
					ObjectMeta: meta_v1.ObjectMeta{
						Name: appName,
						Labels: map[string]string{
							appLabelKey: appName,
						},
					},
				},
			},
			success: true,
		},
		{
			name: "passing one route with name",
			input: []RouteSpecMod{
				{
					ObjectMeta: meta_v1.ObjectMeta{
						Name: "routeName",
					},
				},
			},
			output: []RouteSpecMod{
				{
					ObjectMeta: meta_v1.ObjectMeta{
						Name: "routeName",
						Labels: map[string]string{
							appLabelKey: appName,
						},
					},
				},
			},
			success: true,
		},
		{
			name: "passing multiple ingresses without names",
			input: []RouteSpecMod{
				{},
				{},
			},
			output:  nil,
			success: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fixedRoutes, err := fixRoutes(test.input, appName)

			switch test.success {
			case true:
				if err != nil {
					t.Errorf("Expected test to pass but got an error -\n%v", err)
				}
			case false:
				if err == nil {
					t.Errorf("For the input -\n%v\nexpected test to fail, but test passed", PrettyPrintObjects(test.input))
				}
			}

			if !reflect.DeepEqual(fixedRoutes, test.output) {
				t.Errorf("Expected fixed routes to be -\n%v\nBut got -\n%v\n", PrettyPrintObjects(test.output), PrettyPrintObjects(fixedRoutes))
			}
		})
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
		t.Errorf("expected: %s, got: %s", PrettyPrintObjects(expected),
			PrettyPrintObjects(got))
	}
}

func TestValidateVolumeClaims(t *testing.T) {

	failingTest := []VolumeClaim{
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
				Name: "foo",
			},
		},
	}

	err := validateVolumeClaims(failingTest)
	if err == nil {
		t.Errorf("should have failed but passed for input: %+v", failingTest)
	} else {
		t.Logf("failed with error: %v", err)
	}

}

func TestCreateRoutes(t *testing.T) {
	tests := []struct {
		name   string
		input  *App
		output []runtime.Object
	}{
		{
			name:   "no routes passed",
			input:  &App{},
			output: nil,
		},
		{
			name: "passing 1 route definition",
			input: &App{
				Routes: []RouteSpecMod{
					{
						ObjectMeta: meta_v1.ObjectMeta{
							Name: "testRoute",
						},
						RouteSpec: os_route_v1.RouteSpec{
							Host: "testHost",
						},
					},
				},
			},
			output: []runtime.Object{
				&os_route_v1.Route{
					ObjectMeta: meta_v1.ObjectMeta{
						Name: "testRoute",
					},
					Spec: os_route_v1.RouteSpec{
						Host: "testHost",
					},
				},
			},
		},
		{
			name: "passing 2 route definitions",
			input: &App{
				Routes: []RouteSpecMod{
					{
						ObjectMeta: meta_v1.ObjectMeta{
							Name: "testRoute",
						},
						RouteSpec: os_route_v1.RouteSpec{
							Host: "testHost",
						},
					},
					{
						ObjectMeta: meta_v1.ObjectMeta{
							Name: "testRoute2",
						},
						RouteSpec: os_route_v1.RouteSpec{
							Host: "testHost2",
						},
					},
				},
			},
			output: []runtime.Object{
				&os_route_v1.Route{
					ObjectMeta: meta_v1.ObjectMeta{
						Name: "testRoute",
					},
					Spec: os_route_v1.RouteSpec{
						Host: "testHost",
					},
				},
				&os_route_v1.Route{
					ObjectMeta: meta_v1.ObjectMeta{
						Name: "testRoute2",
					},
					Spec: os_route_v1.RouteSpec{
						Host: "testHost2",
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			objects, err := test.input.createRoutes()
			if err != nil {
				t.Errorf("Creating routes failed: %v", err)
			}
			if !reflect.DeepEqual(test.output, objects) {
				t.Fatalf("Expected:\n%v\nGot:\n%v", PrettyPrintObjects(test.output), PrettyPrintObjects(objects))
			}
		})
	}

}

func TestCreateServices(t *testing.T) {
	tests := []struct {
		Name    string
		App     *App
		Objects []runtime.Object
	}{
		{
			"Single container specified",
			&App{
				ObjectMeta: meta_v1.ObjectMeta{
					Name: "test",
				},
				Deployments: []DeploymentSpecMod{
					{
						PodSpecMod: PodSpecMod{
							Containers: []Container{{Container: api_v1.Container{Image: "nginx"}}},
						},
					},
				},
				Services: []ServiceSpecMod{
					{
						ObjectMeta: meta_v1.ObjectMeta{
							Name: "test",
						},
						Ports: []ServicePortMod{{ServicePort: api_v1.ServicePort{Port: 8080}}}},
				},
			},
			append(make([]runtime.Object, 0), &api_v1.Service{
				ObjectMeta: meta_v1.ObjectMeta{Name: "test"},
				Spec:       api_v1.ServiceSpec{Ports: []api_v1.ServicePort{{Port: 8080, Name: "test-8080"}}},
			}),
		},
		{
			"Single container specified along with shortcut for route(routeEndpoint)",
			&App{
				ObjectMeta: meta_v1.ObjectMeta{
					Name: "test",
				},
				Deployments: []DeploymentSpecMod{
					{
						PodSpecMod: PodSpecMod{
							Containers: []Container{{Container: api_v1.Container{Image: "nginx"}}},
						},
					},
				},
				Services: []ServiceSpecMod{
					{
						ObjectMeta: meta_v1.ObjectMeta{
							Name: "test",
						},
						Ports: []ServicePortMod{{ServicePort: api_v1.ServicePort{Port: 8080}, RouteEndpoint: "xyz.com"}}},
				},
			},
			append(make([]runtime.Object, 0), &api_v1.Service{
				ObjectMeta: meta_v1.ObjectMeta{Name: "test"},
				Spec:       api_v1.ServiceSpec{Ports: []api_v1.ServicePort{{Port: 8080, Name: "test-8080"}}},
			}, &os_route_v1.Route{
				ObjectMeta: meta_v1.ObjectMeta{Name: "test"},
				Spec:       os_route_v1.RouteSpec{Host: "xyz.com", To: os_route_v1.RouteTargetReference{Kind: "Service", Name: "test"}},
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
				t.Fatalf("Expected:\n%v\nGot:\n%v", PrettyPrintObjects(test.Objects), PrettyPrintObjects(object))
			}
		})
	}
}

func TestCreateImageStreams(t *testing.T) {
	tests := []struct {
		name   string
		input  *App
		output []runtime.Object
	}{
		{
			name:   "no imageStreams passed",
			input:  &App{},
			output: nil,
		},
		{
			name: "passing 1 imageStream definition",
			input: &App{
				ImageStreams: []ImageStreamSpecMod{
					{
						ObjectMeta: meta_v1.ObjectMeta{
							Name: "testIS",
						},
						ImageStreamSpec: image_v1.ImageStreamSpec{
							DockerImageRepository: "testRepo",
						},
					},
				},
			},
			output: []runtime.Object{
				&image_v1.ImageStream{
					ObjectMeta: meta_v1.ObjectMeta{
						Name: "testIS",
					},
					Spec: image_v1.ImageStreamSpec{
						DockerImageRepository: "testRepo",
					}},
			},
		},
		{
			name: "passing 2 imageStream definitions",
			input: &App{
				ImageStreams: []ImageStreamSpecMod{
					{
						ObjectMeta: meta_v1.ObjectMeta{
							Name: "testIS1",
						},
						ImageStreamSpec: image_v1.ImageStreamSpec{
							DockerImageRepository: "testRepo1",
						},
					},
					{
						ObjectMeta: meta_v1.ObjectMeta{
							Name: "testIS2",
						},
						ImageStreamSpec: image_v1.ImageStreamSpec{
							DockerImageRepository: "testRepo2",
						},
					},
				},
			},
			output: []runtime.Object{
				&image_v1.ImageStream{
					ObjectMeta: meta_v1.ObjectMeta{
						Name: "testIS1",
					},
					Spec: image_v1.ImageStreamSpec{
						DockerImageRepository: "testRepo1",
					}},
				&image_v1.ImageStream{
					ObjectMeta: meta_v1.ObjectMeta{
						Name: "testIS2",
					},
					Spec: image_v1.ImageStreamSpec{
						DockerImageRepository: "testRepo2",
					}},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			objects, err := test.input.createImageStreams()
			if err != nil {
				t.Errorf("Creating imageStreams failed: %v", err)
			}
			if !reflect.DeepEqual(test.output, objects) {
				t.Fatalf("Expected:\n%v\nGot:\n%v", PrettyPrintObjects(test.output), PrettyPrintObjects(objects))
			}
		})
	}
}

func TestCreateBuildConfigs(t *testing.T) {
	tests := []struct {
		name   string
		input  *App
		output []runtime.Object
	}{
		{
			name:   "no buildConfig passed",
			input:  &App{},
			output: nil,
		},
		{
			name: "passing 1 buildConfig definition",
			input: &App{
				BuildConfigs: []BuildConfigSpecMod{
					{
						ObjectMeta: meta_v1.ObjectMeta{
							Name: "testBC",
						},
						BuildConfigSpec: build_v1.BuildConfigSpec{
							RunPolicy: "Serial",
						},
					},
				},
			},
			output: []runtime.Object{
				&build_v1.BuildConfig{
					ObjectMeta: meta_v1.ObjectMeta{
						Name: "testBC",
					},
					Spec: build_v1.BuildConfigSpec{
						RunPolicy: "Serial",
					},
				},
			},
		},
		{
			name: "passing 2 buildConfig definitions",
			input: &App{
				BuildConfigs: []BuildConfigSpecMod{
					{
						ObjectMeta: meta_v1.ObjectMeta{
							Name: "testBC1",
						},
						BuildConfigSpec: build_v1.BuildConfigSpec{
							RunPolicy: "Serial",
						},
					},
					{
						ObjectMeta: meta_v1.ObjectMeta{
							Name: "testBC2",
						},
						BuildConfigSpec: build_v1.BuildConfigSpec{
							RunPolicy: "Serial",
						},
					},
				},
			},
			output: []runtime.Object{
				&build_v1.BuildConfig{
					ObjectMeta: meta_v1.ObjectMeta{
						Name: "testBC1",
					},
					Spec: build_v1.BuildConfigSpec{
						RunPolicy: "Serial",
					},
				},
				&build_v1.BuildConfig{
					ObjectMeta: meta_v1.ObjectMeta{
						Name: "testBC2",
					},
					Spec: build_v1.BuildConfigSpec{
						RunPolicy: "Serial",
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			objects, err := test.input.createBuildConfigs()
			if err != nil {
				t.Errorf("Creating buildConfigs failed: %v", err)
			}
			if !reflect.DeepEqual(test.output, objects) {
				t.Fatalf("Expected:\n%v\nGot:\n%v", PrettyPrintObjects(test.output), PrettyPrintObjects(objects))
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
			name:        "port/protocol(lowercase) is passed",
			portMapping: "1337/udp",
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
			name:        "port/protocol(tCp) is passed",
			portMapping: "1337/tCp",
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

// TODO: add test for auto naming of single persistent volume
