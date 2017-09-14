package e2e

import (
	"bytes"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/clientcmd"
)

// Hardcoding the location of kedge, which is in root of project directory
var ProjectPath = "$GOPATH/src/github.com/kedgeproject/kedge/"
var KedgeLoc = ProjectPath + "kedge"

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func createClient() (*kubernetes.Clientset, error) {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	return kubernetes.NewForConfig(config)
}

func createNS(clientset *kubernetes.Clientset, name string) (*v1.Namespace, error) {
	ns := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	return clientset.CoreV1().Namespaces().Create(ns)
}

func RunKedge(files []string, namespace string) ([]byte, error) {
	args := []string{"create", "-n", namespace}
	for _, file := range files {
		args = append(args, "-f")
		args = append(args, os.ExpandEnv(file))
	}
	cmd := exec.Command(os.ExpandEnv(KedgeLoc), args...)

	var out, stdErr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stdErr

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("error running %q\n%s %s",
			fmt.Sprintf("kedge %s", strings.Join(args, " ")),
			stdErr.String(), err)
	}
	return out.Bytes(), nil
}

func mapkeys(m map[string]int) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func PodsStarted(t *testing.T, clientset *kubernetes.Clientset, namespace string, podNames []string) error {
	// convert podNames to map
	podUp := make(map[string]int)
	for _, p := range podNames {
		podUp[p] = 0
	}

	// Timeouts after 9 minutes if the Pod has not yet started
	// 9 minute reasoning = 1 minute before 10-minute Golang test timeout.
	timeout := time.After(9 * time.Minute)
	tick := time.Tick(time.Second)

	for {
		select {
		case <-timeout:
			return fmt.Errorf("pods did not come up in given time: 5 minutes")
		case <-tick:
			t.Logf("pods not started yet: %q", strings.Join(mapkeys(podUp), " "))

			pods, err := clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{})
			if err != nil {
				return errors.Wrap(err, "error while listing all pods")
			}
			// iterate on all pods we care about
			for k := range podUp {
				for _, p := range pods.Items {
					if strings.Contains(p.Name, k) && p.Status.Phase == v1.PodRunning {
						t.Logf("Pod %q started!", p.Name)
						delete(podUp, k)
					}
				}
			}
		}
		if len(podUp) == 0 {
			break
		}
	}
	return nil
}

func getEndPoints(t *testing.T, clientset *kubernetes.Clientset, namespace string, svcs []ServicePort) (map[string]string, error) {
	// find the minikube ip
	node, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "error while listing all nodes")
	}
	nodeIP := node.Items[0].Status.Addresses[0].Address
	t.Logf("node ip address %s", nodeIP)

	// get all running services
	runningSvcs, err := clientset.CoreV1().Services(namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "error while listing all services")
	}

	endpoint := make(map[string]string)
	for _, svc := range svcs {
		for _, s := range runningSvcs.Items {
			if s.Name == svc.Name {
				for _, p := range s.Spec.Ports {
					if p.Port == svc.Port {
						port := p.NodePort
						v := fmt.Sprintf("http://%s:%d", nodeIP, port)
						k := fmt.Sprintf("%s:%d", svc.Name, svc.Port)
						endpoint[k] = v
					}
				}
			}
		}
	}
	t.Logf("endpoints: %#v", endpoint)
	return endpoint, nil
}

func pingEndPoints(t *testing.T, ep map[string]string) error {
	timeout := time.After(5 * time.Minute)
	tick := time.Tick(time.Second)

	for {
		select {
		case <-timeout:
			return fmt.Errorf("could not curl on the service in given time: 5 minutes")
		case <-tick:
			for e, u := range ep {
				timeout := time.Duration(5 * time.Second)
				client := http.Client{
					Timeout: timeout,
				}
				respose, err := client.Get(u)
				if err != nil {
					t.Logf("error while making http request %q for service %q, err: %v", u, e, err)
					time.Sleep(1 * time.Second)
					continue
				}
				if respose.Status == "200 OK" {
					t.Logf("%q is running!", e)
					delete(ep, e)
				} else {
					return fmt.Errorf("for service %q got %q", e, respose.Status)
				}
			}
		}
		if len(ep) == 0 {
			break
		}
	}
	return nil
}

func deleteNamespace(t *testing.T, clientset *kubernetes.Clientset, namespace string) {
	if err := clientset.CoreV1().Namespaces().Delete(namespace, &metav1.DeleteOptions{}); err != nil {
		t.Logf("error deleting namespace %q: %v", namespace, err)
	}
	t.Logf("successfully deleted namespace: %q", namespace)
}

type ServicePort struct {
	Name string
	Port int32
}

type testData struct {
	TestName         string
	Namespace        string
	InputFiles       []string
	PodStarted       []string
	NodePortServices []ServicePort
}

func Test_Integration(t *testing.T) {
	clientset, err := createClient()
	if err != nil {
		t.Fatalf("error getting kube client: %v", err)
	}

	tests := []testData{
		{
			TestName:  "Testing configMap",
			Namespace: "configmap",
			InputFiles: []string{
				ProjectPath + "docs/examples/configmap/db.yaml",
				ProjectPath + "docs/examples/configmap/web.yaml",
			},
			PodStarted: []string{"web"},
			NodePortServices: []ServicePort{
				{Name: "wordpress", Port: 8080},
			},
		},
		{
			TestName:  "Testing customVol",
			Namespace: "customvol",
			InputFiles: []string{
				ProjectPath + "docs/examples/customVol/db.yaml",
				ProjectPath + "docs/examples/customVol/web.yaml",
			},
			PodStarted: []string{"web"},
			NodePortServices: []ServicePort{
				{Name: "wordpress", Port: 8080},
			},
		},
		{
			TestName:  "Testing envFrom",
			Namespace: "envfrom",
			InputFiles: []string{
				ProjectPath + "docs/examples/envFrom/db.yaml",
				ProjectPath + "docs/examples/envFrom/web.yaml",
			},
			PodStarted: []string{"web"},
			NodePortServices: []ServicePort{
				{Name: "wordpress", Port: 8080},
			},
		},
		{
			TestName:  "Testing extraResources",
			Namespace: "extraresources",
			InputFiles: []string{
				ProjectPath + "docs/examples/extraResources/app.yaml",
			},
			PodStarted: []string{"web"},
			NodePortServices: []ServicePort{
				{Name: "web", Port: 80},
			},
		},
		{
			TestName:  "Testing health",
			Namespace: "health",
			InputFiles: []string{
				ProjectPath + "docs/examples/health/db.yaml",
				ProjectPath + "docs/examples/health/web.yaml",
			},
			PodStarted: []string{"web"},
			NodePortServices: []ServicePort{
				{Name: "wordpress", Port: 8080},
			},
		},
		{
			TestName:  "Testing healthChecks",
			Namespace: "healthchecks",
			InputFiles: []string{
				ProjectPath + "docs/examples/healthchecks/db.yaml",
				ProjectPath + "docs/examples/healthchecks/web.yaml",
			},
			PodStarted: []string{"web"},
			NodePortServices: []ServicePort{
				{Name: "wordpress", Port: 8080},
			},
		},
		{
			TestName:  "Testing secret",
			Namespace: "secrets",
			InputFiles: []string{
				ProjectPath + "docs/examples/secrets/db.yaml",
				ProjectPath + "docs/examples/secrets/web.yaml",
			},
			PodStarted: []string{"web"},
			NodePortServices: []ServicePort{
				{Name: "wordpress", Port: 8080},
			},
		},
		{
			TestName:  "Testing single file",
			Namespace: "singlefile",
			InputFiles: []string{
				ProjectPath + "docs/examples/single_file/wordpress.yml",
			},
			PodStarted: []string{"wordpress"},
			NodePortServices: []ServicePort{
				{Name: "wordpress", Port: 8080},
			},
		},
		{
			TestName:  "Normal Wordpress test",
			Namespace: "wordpress",
			InputFiles: []string{
				ProjectPath + "examples/wordpress/wordpress.yaml",
				ProjectPath + "examples/wordpress/mariadb.yaml",
			},
			PodStarted: []string{"wordpress"},
			NodePortServices: []ServicePort{
				{Name: "wordpress", Port: 80},
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
			convertedOutput, err := RunKedge(test.InputFiles, test.Namespace)
			if err != nil {
				t.Fatalf("error running kedge: %v", err)
			}
			t.Log(string(convertedOutput))

			// see if the pods are running
			if err := PodsStarted(t, clientset, test.Namespace, test.PodStarted); err != nil {
				t.Fatalf("error finding running pods: %v", err)
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
