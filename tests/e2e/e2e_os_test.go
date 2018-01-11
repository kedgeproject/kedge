package e2e

import (
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"strconv"
	"strings"
	"testing"
)

func waitForBuildComplete(namespace string, buildName string) error {
	return wait.Poll(waitInterval, jobTimeout, func() (bool, error) {
		var buildOut bool
		buildStatus, err := runCmd("oc get build --namespace=" + namespace +
			" --template='{{ range .items }}{{ if (eq .metadata.name \"" +
			buildName + "\") }}{{ eq .status.phase \"Complete\" }}{{end}}{{end}}'")
		if err != nil {
			return false, errors.Wrap(err, "error getting build status")
		}

		buildOut, _ = strconv.ParseBool(string(buildStatus))
		return buildOut, nil
	})
}

func runKedgeS2i(imageName string, baseImage string) error {
	s2iCmd := BinaryLocation + " build --s2i --image " + imageName + " -b " + baseImage
	_, err := runCmd(s2iCmd)
	if err != nil {
		return errors.Wrap(err, "error build s2i image")
	}
	return nil
}

// TODO: Use OpenShift client-go API instead of go-template
func Test_os_Integration(t *testing.T) {
	clientset, err := createClient()
	if err != nil {
		t.Fatalf("error getting kube client: %v", err)
	}

	tests := []testData{
		{
			TestName:  "Testing routes",
			Namespace: "testroutes",
			InputFiles: []string{
				ProjectPath + TestPath + "routes/httpd.yml",
			},
			PodStarted: []string{"httpd"},
			NodePortServices: []ServicePort{
				{Name: "httpd", Port: 8080},
			},
		},

		{
			TestName:  "Testing buildconfig",
			Namespace: "buildconfig",
			InputFiles: []string{
				ProjectPath + TestPath + "buildConfigs/ruby-hello-world.yml",
			},
			PodStarted: []string{"ruby"},
			Type:       "buildconfig",
		},

		{
			TestName:  "Testing s2i image",
			Namespace: "s2i",
			InputFiles: []string{
				ProjectPath + TestPath + "s2i/configs/",
			},
			PodStarted: []string{"redis"},
			Type:       "s2i",
			BaseImage:  "centos/python-35-centos7:3.5",
		},
	}

	_, err = clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		t.Fatalf("OC cluster is not running or not accessible: %v", err)
	}

	for _, test := range tests {
		test := test // capture range variable
		t.Run(test.TestName, func(t *testing.T) {
			t.Parallel()
			// create an OpenShift project
			_, err := runCmd("oc new-project " + test.Namespace)
			if err != nil {
				t.Fatalf("error creating project: %v", err)
			}
			t.Logf("project %q created", test.Namespace)

			defer deleteNamespace(t, clientset, test.Namespace)

			if test.Type == "s2i" {
				err := runKedgeS2i(test.Namespace, test.BaseImage)
				if err != nil {
					t.Fatalf("error running kedge s2i: %v", err)
				}
			}

			// run kedge
			convertedOutput, err := RunBinary(test.InputFiles, test.Namespace)
			if err != nil {
				t.Fatalf("error running kedge: %v", err)
			}
			t.Log(string(convertedOutput))

			// see if the pods are running
			if err := PodsStarted(t, clientset, test.Namespace, test.PodStarted); err != nil {
				t.Fatalf("error finding running pods: %v", err)
			}

			if (test.Type == "buildconfig") || (test.Type == "s2i") {
				listBuilds, err := runCmd("oc get build --namespace=" + test.Namespace +
					" --template='{{range .items}}{{ .metadata.name }}:{{end}}' | tr \":\" \"\n\"")
				if err != nil {
					t.Fatalf("error getting the Build list: %v", err)
				}

				eachBuild := strings.Split(string(listBuilds), "\n")
				for _, build := range eachBuild {
					if len(build) > 0 {
						err := waitForBuildComplete(test.Namespace, build)
						if err != nil {
							t.Fatalf("Build failed: %v", err)
						}

						t.Logf("Successfully completed the build: %s", build)
					}
				}
			}

			// get endpoints for all services
			endPoints, err := getEndPoints(t, clientset, test.Namespace, test.NodePortServices)
			if err != nil {
				t.Fatalf("error getting nodes: %v", err)
			}

			if err := pingEndPoints(t, endPoints); err != nil {
				t.Fatalf("error pinging endpoint: %v", err)
			}

			t.Logf("Successfully pinged all endpoints!")

		})
	}
}
