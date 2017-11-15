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

package build

import (
	"testing"

	_ "k8s.io/kubernetes/pkg/api/install"
)

func TestGetImageTag(t *testing.T) {

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "image FQDN given with port of registry and image tag",
			input: "myregistryhost:5000/fedora/httpd:version1.0",
			want:  "version1.0",
		},
		{
			name:  "image FQDN given with port of registry and no image tag",
			input: "myregistryhost:5000/fedora/httpd",
			want:  "latest",
		},
		{
			name:  "image FQDN given with image tag",
			input: "myregistryhost/fedora/httpd:version1.0",
			want:  "version1.0",
		},
		{
			name:  "image FQDN given without image tag",
			input: "myregistryhost/fedora/httpd",
			want:  "latest",
		},
		{
			name:  "repo name/image name without image tag",
			input: "fedora/httpd",
			want:  "latest",
		},
		{
			name:  "image name without image tag",
			input: "httpd",
			want:  "latest",
		},
		{
			name:  "repo name/image name with image tag",
			input: "fedora/httpd:v1",
			want:  "v1",
		},
		{
			name:  "image name with image tag",
			input: "httpd:v1",
			want:  "v1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetImageTag(tt.input); got != tt.want {
				t.Errorf("GetImageTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetImageName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "image name with image tag",
			input: "httpd:v1",
			want:  "httpd",
		},
		{
			name:  "repo name/image name with image tag",
			input: "fedora/httpd:v1",
			want:  "httpd",
		},
		{
			name:  "image name without image tag",
			input: "httpd",
			want:  "httpd",
		},
		{
			name:  "repo name/image name without image tag",
			input: "fedora/httpd",
			want:  "httpd",
		},
		{
			name:  "image FQDN given without image tag",
			input: "myregistryhost/fedora/httpd",
			want:  "httpd",
		},
		{
			name:  "image FQDN given with image tag",
			input: "myregistryhost/fedora/httpd:version1.0",
			want:  "httpd",
		},
		{
			name:  "image FQDN given with port of registry and no image tag",
			input: "myregistryhost:5000/fedora/httpd",
			want:  "httpd",
		},
		{
			name:  "image FQDN given with port of registry and image tag",
			input: "myregistryhost:5000/fedora/httpd:version1.0",
			want:  "httpd",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetImageName(tt.input); got != tt.want {
				t.Errorf("GetImageName() = %v, want %v", got, tt.want)
			}
		})
	}
}
