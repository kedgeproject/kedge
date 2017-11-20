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
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/clientcmd"
)

/*
	This file tests the functionality of a tool that deploy to Kubernetes and checks
	to see if that tools successfully deploys said file.
*/

// Hardcoding the location of the binary, which is in root of project directory
var ProjectPath = "$GOPATH/src/github.com/kedgeproject/kedge/"
var TestPath = "docs/examples/"
var BinaryLocation = ProjectPath + "kedge"
var BinaryCommand = []string{"create", "-n"}

const (
	jobTimeout   = 10 * time.Minute
	waitInterval = 5
)

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

// Create a Kubernetes client
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

// Create the repsective namespace
func createNS(clientset *kubernetes.Clientset, name string) (*v1.Namespace, error) {
	ns := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	return clientset.CoreV1().Namespaces().Create(ns)
}

// Run the binary against the Kubernetes cluster using a specified command against a specific namespace
// requirement: your binary must have a --namespace parameter to specify a namespace location as well as
// -f to specific a file. Ex: command --namespace foobar -f foo.yaml -f bar.yaml
func RunBinary(files []string, namespace string) ([]byte, error) {
	args := append(BinaryCommand, namespace)
	for _, file := range files {
		args = append(args, "-f")
		args = append(args, os.ExpandEnv(file))
	}
	cmd := exec.Command(os.ExpandEnv(BinaryLocation), args...)

	var out, stdErr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stdErr

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("error running %q\n%s %s",
			fmt.Sprintf("command: %s", strings.Join(args, " ")),
			stdErr.String(), err)
	}
	return out.Bytes(), nil
}

// Map specific keys (utility function)
func mapkeys(m map[string]int) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Check to see if specific pods have been started
func PodsStarted(t *testing.T, clientset *kubernetes.Clientset, namespace string, podNames []string) error {
	// convert podNames to map
	podUp := make(map[string]int)
	for _, p := range podNames {
		podUp[p] = 0
	}

	// Timeouts after 10 minutes if the Pod has not yet started
	podTimeout := time.After(10 * time.Minute)
	tick := time.Tick(time.Second)

	for {
		select {
		case <-podTimeout:
			return fmt.Errorf("pods did not come up in given time: 10 minutes")
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

// Wait for the jobs to complete
func waitForJobComplete(clientset *kubernetes.Clientset, namespace string, jobName string) error {
	return wait.Poll(waitInterval, jobTimeout, func() (bool, error) {
		jobStatus, err := clientset.Batch().Jobs(namespace).Get(jobName, metav1.GetOptions{})
		if err != nil {
			return false, errors.Wrap(err, "error getting jobs")
		}

		return (jobStatus.Status.Failed == 0) && (jobStatus.Status.Active == 0), nil
	})
}

// Retrieve all required end-points via minikube (required!)
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

// Ping all endpoints
func pingEndPoints(t *testing.T, ep map[string]string) error {

	// Increase pinging timeout to more than 8 minutes to allocate
	// for replica containers across minikube
	pingTimeout := time.After(15 * time.Minute)
	tick := time.Tick(time.Second)

	for {
		select {
		case <-pingTimeout:
			return fmt.Errorf("could not ping the specific service in given time: 15 minutes")
		case <-tick:
			for e, u := range ep {
				// 10 second timeout for HTTP response
				httpTimeout := time.Duration(10 * time.Second)
				client := http.Client{
					Timeout: httpTimeout,
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

// Delete the namespace
func deleteNamespace(t *testing.T, clientset *kubernetes.Clientset, namespace string) {
	if err := clientset.CoreV1().Namespaces().Delete(namespace, &metav1.DeleteOptions{}); err != nil {
		t.Logf("error deleting namespace %q: %v", namespace, err)
	}
	t.Logf("successfully deleted namespace: %q", namespace)
}

// These structs create a specific name as well as port to ping
type ServicePort struct {
	Name string
	Port int32
}

// Here we will test all of our test data!
type testData struct {
	TestName         string
	Namespace        string
	InputFiles       []string
	PodStarted       []string
	NodePortServices []ServicePort
	Type             string
	BaseImage        string
}
