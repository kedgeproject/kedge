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
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	api_v1 "k8s.io/client-go/pkg/api/v1"
	batch_v1 "k8s.io/client-go/pkg/apis/batch/v1"

	"reflect"
)

func TestIsVolumeDefined(t *testing.T) {
	volumes := []api_v1.Volume{{Name: "foo"}, {Name: "bar"}, {Name: "baz"}}

	tests := []struct {
		Search string
		Output bool
	}{
		{Search: "bar", Output: true},
		{Search: "fooo", Output: false},
		{Search: "baz", Output: true},
	}

	t.Logf("volumes: %+v", volumes)
	for _, test := range tests {
		if test.Output != isVolumeDefined(volumes, test.Search) {
			t.Errorf("expected output to match, but did not match for"+
				" volumes %+v and search query %q", volumes, test.Search)
		} else {
			t.Logf("test passed for search query %q", test.Search)
		}
	}

}

func TestIsPVCDefined(t *testing.T) {
	volumes := []VolumeClaim{
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
				Name: "baz",
			},
		},
	}

	tests := []struct {
		Search string
		Output bool
	}{
		{Search: "bar", Output: true},
		{Search: "fooo", Output: false},
		{Search: "baz", Output: true},
	}

	t.Logf("volumes: %+v", volumes)
	for _, test := range tests {
		if test.Output != isPVCDefined(volumes, test.Search) {
			t.Errorf("expected output to match, but did not match for"+
				" volumes %+v and search query %q", volumes, test.Search)
		} else {
			t.Logf("test passed for search query %q", test.Search)
		}
	}
}

func TestSetGVK(t *testing.T) {
	jobTest := struct {
		name         string
		beforeObject *batch_v1.Job
		afterObject  *batch_v1.Job
	}{
		name: "Set GVK for a Job",
		beforeObject: &batch_v1.Job{
			ObjectMeta: v1.ObjectMeta{
				Name: "testJob",
			},
			Spec: batch_v1.JobSpec{
				Parallelism: getInt32Addr(2),
			},
		},
		afterObject: &batch_v1.Job{
			TypeMeta: v1.TypeMeta{
				Kind:       "Job",
				APIVersion: "batch/v1",
			},
			ObjectMeta: v1.ObjectMeta{
				Name: "testJob",
			},
			Spec: batch_v1.JobSpec{
				Parallelism: getInt32Addr(2),
			},
		},
	}

	t.Run(jobTest.name, func(t *testing.T) {

		scheme, err := GetScheme()
		if err != nil {
			t.Fatalf("unable to get scheme - %v", err)
		}

		if err := SetGVK(jobTest.beforeObject, scheme); err != nil {
			t.Fatalf("unable to set GVK - %v", err)
		}

		if !reflect.DeepEqual(jobTest.beforeObject, jobTest.afterObject) {
			t.Errorf("Expected runtime object after setting GVK to be -\n%v\nBut got -\n%v", jobTest.afterObject, jobTest.beforeObject)
		}
	})
}
