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

	"encoding/json"

	log "github.com/Sirupsen/logrus"
	os_deploy_v1 "github.com/openshift/origin/pkg/apps/apis/apps/v1"
	build_v1 "github.com/openshift/origin/pkg/build/apis/build/v1"
	image_v1 "github.com/openshift/origin/pkg/image/apis/image/v1"
	os_route_v1 "github.com/openshift/origin/pkg/route/apis/route/v1"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/api"
	api_v1 "k8s.io/kubernetes/pkg/api/v1"
	batch_v1 "k8s.io/kubernetes/pkg/apis/batch/v1"
)

// This function will search in the pod level volumes
// and see if the volume with given name is defined
func isVolumeDefined(volumes []api_v1.Volume, name string) bool {
	for _, v := range volumes {
		if v.Name == name {
			return true
		}
	}
	return false
}

// search through all the persistent volumes defined in the root level
func isPVCDefined(volumes []VolumeClaim, name string) bool {
	for _, v := range volumes {
		if v.Name == name {
			return true
		}
	}
	return false
}

// GetScheme() returns runtime.Scheme with supported Kubernetes API resource
// definitions which Kedge supports right now.
// The core v1 scheme is first initialized and then other controllers' scheme
// is added to that scheme, e.g. batch/v1 scheme is added to add support for
// Jobs controller to the v1 Scheme.
// Also, (from upstream) Scheme defines methods for serializing and deserializing API objects, a type
// registry for converting group, version, and kind information to and from Go
// schemas, and mappings between Go schemas of different versions. A scheme is the
// foundation for a versioned API and versioned configuration over time.
func GetScheme() (*runtime.Scheme, error) {
	// Initializing the scheme with the core v1 api
	scheme := api.Scheme

	// TODO: find a way where we don't have to add all the subsequent schemes
	// to the v1 scheme, instead we should be able to have different scheme for
	// different controllers

	// Adding the apps scheme to support DeploymentConfig
	if err := os_deploy_v1.AddToScheme(scheme); err != nil {
		return nil, errors.Wrap(err, "unable to add 'apps' (OpenShift) to scheme")
	}

	// Adding the batch scheme to support Jobs
	if err := batch_v1.AddToScheme(scheme); err != nil {
		return nil, errors.Wrap(err, "unable to add 'batch' to scheme")
	}

	// Adding the route scheme to support OpenShift routes
	if err := os_route_v1.AddToScheme(scheme); err != nil {
		return nil, errors.Wrap(err, "unable to add 'routes' to scheme")
	}

	// Adding the image scheme to support OpenShift ImageStreams
	if err := image_v1.AddToScheme(scheme); err != nil {
		return nil, errors.Wrap(err, "unable to add 'image' to scheme")
	}

	// Adding the build scheme to support OpenShift buildConfigs
	if err := build_v1.AddToScheme(scheme); err != nil {
		return nil, errors.Wrap(err, "unable to add 'build' to scheme")
	}

	return scheme, nil
}

// SetGVK() sets Group, Version and Kind for the generated Kubernetes resources.
// This takes in a generated Kubernetes API resource's runtime object and
// runtime scheme based on which the GVK will be set.
func SetGVK(runtimeObject runtime.Object, scheme *runtime.Scheme) error {
	gvk, isUnversioned, err := scheme.ObjectKind(runtimeObject)
	if err != nil {
		return errors.Wrap(err, "ConvertToVersion failed")
	}
	if isUnversioned {
		return fmt.Errorf("ConvertToVersion failed: can't output unversioned type: %T", runtimeObject)
	}
	runtimeObject.GetObjectKind().SetGroupVersionKind(gvk)
	return nil
}

func getInt32Addr(i int32) *int32 {
	return &i
}

func getInt64Addr(i int64) *int64 {
	return &i
}

func PrettyPrintObjects(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}

// addKeyValueToMap adds a key value pair to a given map[string]string only if
// the map does not contain the supplied key. Creates a new map if map is empty.
// We need to return the map because in case a nil map is passed to this
// function, the new map created will not be reflected in the original nil map.
func addKeyValueToMap(k string, v string, m map[string]string) map[string]string {

	if len(m) == 0 {
		m = make(map[string]string)
	}

	if _, ok := m[k]; !ok {
		m[k] = v
	} else {
		log.Debugf("not adding '%v: %v' to map since there exists a user defined label '%v: %v'", k, v, k, m[k])
	}

	return m
}
