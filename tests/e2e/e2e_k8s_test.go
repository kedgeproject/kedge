package e2e

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

// The "bread and butter" of the test-suite. We will iterate through
// each test that is required and make sure that not only are the pods started
// but that each test is pingable / is accessable.
func Test_k8s_Integration(t *testing.T) {
	clientset, err := createClient()
	if err != nil {
		t.Fatalf("error getting kube client: %v", err)
	}

	tests := []testData{
		{
			TestName:  "Testing configmap",
			Namespace: "configmap",
			InputFiles: []string{
				ProjectPath + TestPath + "configmap/guestbook.yaml",
				ProjectPath + TestPath + "configmap/redis-master.yaml",
				ProjectPath + TestPath + "configmap/redis-slave.yaml",
			},
			PodStarted: []string{"guestbook", "redis-master", "redis-slave"},
			NodePortServices: []ServicePort{
				{Name: "guestbook", Port: 3000},
			},
		},
		{
			TestName:  "Testing custom-vol",
			Namespace: "custom-vol",
			InputFiles: []string{
				ProjectPath + TestPath + "custom-vol/guestbook.yaml",
				ProjectPath + TestPath + "custom-vol/redis-master.yaml",
				ProjectPath + TestPath + "custom-vol/redis-slave.yaml",
			},
			PodStarted: []string{"guestbook", "redis-master", "redis-slave"},
			NodePortServices: []ServicePort{
				{Name: "guestbook", Port: 3000},
			},
		},
		{
			TestName:  "Testing health",
			Namespace: "health",
			InputFiles: []string{
				ProjectPath + TestPath + "health/guestbook.yaml",
				ProjectPath + TestPath + "health/redis-master.yaml",
				ProjectPath + TestPath + "health/redis-slave.yaml",
			},
			PodStarted: []string{"guestbook", "redis-master", "redis-slave"},
			NodePortServices: []ServicePort{
				{Name: "guestbook", Port: 3000},
			},
		},
		{
			TestName:  "Testing health-check",
			Namespace: "health-check",
			InputFiles: []string{
				ProjectPath + TestPath + "health-check/guestbook.yaml",
				ProjectPath + TestPath + "health-check/redis-master.yaml",
				ProjectPath + TestPath + "health-check/redis-slave.yaml",
			},
			PodStarted: []string{"guestbook", "redis-master", "redis-slave"},
			NodePortServices: []ServicePort{
				{Name: "guestbook", Port: 3000},
			},
		},
		{
			TestName:  "Testing secret",
			Namespace: "secret",
			InputFiles: []string{
				ProjectPath + TestPath + "secret/guestbook.yaml",
				ProjectPath + TestPath + "secret/redis-master.yaml",
				ProjectPath + TestPath + "secret/redis-slave.yaml",
			},
			PodStarted: []string{"guestbook", "redis-master", "redis-slave"},
			NodePortServices: []ServicePort{
				{Name: "guestbook", Port: 3000},
			},
		},
		{
			TestName:  "Test port-mappings",
			Namespace: "port-mappings",
			InputFiles: []string{
				ProjectPath + TestPath + "port-mappings/guestbook.yaml",
				ProjectPath + TestPath + "port-mappings/redis-master.yaml",
				ProjectPath + TestPath + "port-mappings/redis-slave.yaml",
			},
			PodStarted: []string{"guestbook", "redis-master", "redis-slave"},
			NodePortServices: []ServicePort{
				{Name: "guestbook", Port: 3000},
			},
		},
		// Special non-guestbook-go-redis tests
		{
			TestName:  "Testing single-file",
			Namespace: "single-file",
			InputFiles: []string{
				ProjectPath + TestPath + "single-file/guestbook.yaml",
			},
			PodStarted: []string{"guestbook", "redis-master", "redis-slave"},
			NodePortServices: []ServicePort{
				{Name: "guestbook", Port: 3000},
			},
		},
		{
			TestName:  "Test jobs",
			Namespace: "jobs",
			InputFiles: []string{
				ProjectPath + TestPath + "jobs/job.yaml",
			},
			PodStarted: []string{"pival"},
			Type:       "job",
		},
		{
			TestName:  "Testing include-resources",
			Namespace: "include-resources",
			InputFiles: []string{
				ProjectPath + TestPath + "include-resources/web.yaml",
			},
			PodStarted: []string{"web"},
			NodePortServices: []ServicePort{
				{Name: "web", Port: 80},
			},
		},
	}

	_, err = clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		t.Fatalf("Kubernetes cluster is not running or not accessible: %v", err)
	}

	for _, test := range tests {
		test := test // capture range variable
		t.Run(test.TestName, func(t *testing.T) {
			t.Parallel()
			// create a namespace
			_, err := createNS(clientset, test.Namespace)
			if err != nil {
				t.Fatalf("error creating namespace: %v", err)
			}
			t.Logf("namespace %q created", test.Namespace)
			defer deleteNamespace(t, clientset, test.Namespace)

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

			if test.Type == "job" {
				listJobs, err := clientset.Batch().Jobs(test.Namespace).List(metav1.ListOptions{})
				if err != nil {
					t.Fatalf("error getting the job list: %v", err)
				}

				for _, job := range listJobs.Items {
					err := waitForJobComplete(clientset, test.Namespace, job.Name)
					if err != nil {
						t.Fatalf("Job failed: %v", err)
					}

					t.Logf("Successfully completed the job: %s", job.Name)
				}
			} else {

				// get endpoints for all services
				endPoints, err := getEndPoints(t, clientset, test.Namespace, test.NodePortServices)
				if err != nil {
					t.Fatalf("error getting nodes: %v", err)
				}

				if err := pingEndPoints(t, endPoints); err != nil {
					t.Fatalf("error pinging endpoint: %v", err)
				}

				t.Logf("Successfully pinged all endpoints!")
			}
		})
	}
}
