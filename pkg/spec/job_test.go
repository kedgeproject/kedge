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
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	api_v1 "k8s.io/kubernetes/pkg/api/v1"
	batch_v1 "k8s.io/kubernetes/pkg/apis/batch/v1"
)

func TestJobSpecMod_CreateKubernetesController(t *testing.T) {
	tests := []struct {
		name          string
		jobSpecMod    *JobSpecMod
		kubernetesJob *batch_v1.Job
		success       bool
	}{
		{
			name: "both PodSpec and JobSpec are empty, no Job should be created",
			jobSpecMod: &JobSpecMod{
				ControllerFields: ControllerFields{
					Controller: "job",
					Secrets: []SecretMod{
						{
							Secret: api_v1.Secret{
								ObjectMeta: meta_v1.ObjectMeta{
									Name: "secret",
								},
								StringData: map[string]string{
									"testData": "testValue",
								},
							},
						},
					},
				},
			},
			kubernetesJob: nil,
			success:       true,
		},
		{
			name: "ActiveDeadlineSeconds is specified, make sure it's only populated for JobSpec and not for PodSpec",
			jobSpecMod: &JobSpecMod{
				ControllerFields: ControllerFields{
					ObjectMeta: meta_v1.ObjectMeta{
						Name: "testJob",
					},
					Controller: "job",
					PodSpecMod: PodSpecMod{
						PodSpec: api_v1.PodSpec{
							Containers: []api_v1.Container{
								{
									Name:  "testContainer",
									Image: "testImage",
								},
							},
						},
					},
				},

				ActiveDeadlineSeconds: getInt64Addr(20),
				JobSpec: batch_v1.JobSpec{
					Parallelism: getInt32Addr(2),
				},
			},
			kubernetesJob: &batch_v1.Job{
				ObjectMeta: meta_v1.ObjectMeta{
					Name: "testJob",
				},
				Spec: batch_v1.JobSpec{
					Parallelism:           getInt32Addr(2),
					ActiveDeadlineSeconds: getInt64Addr(20),
					Template: api_v1.PodTemplateSpec{
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
			success: true,
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			kJob, err := test.jobSpecMod.createKubernetesController()

			switch test.success {
			case true:
				if err != nil {
					t.Errorf("Expected test to pass but got an error -\n%v", err)
				}
			case false:
				if err == nil {
					t.Errorf("For the input -\n%v\nexpected test to fail, but test passed", spew.Sprint(test.jobSpecMod))
				}
			}

			if !reflect.DeepEqual(test.kubernetesJob, kJob) {

				t.Errorf("Expected Kubernetes Job to be -\n%v\nBut got -\n%v", prettyPrintObjects(test.kubernetesJob), prettyPrintObjects(kJob))
			}
		})
	}
}
