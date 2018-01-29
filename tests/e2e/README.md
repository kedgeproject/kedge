# Running e2e tests on Kedge

The e2e tests leverages on [go testing package](https://golang.org/pkg/testing/) to run tests. In these tests, 
we execute `kedge create` on Kubernetes and OpenShift cluster with examples provided under [docs/examples](/docs/examples).

### Pre-requisites

#### kubectl 
Both kubernetes and OpenShift tests rely on [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) 
for validation. You need to have kubectl installed before running the tests.


### Running Kubernetes tests

1. For running the Kubernetes tests, you need a kubernetes cluster up and running. 
   [Minikube](https://github.com/kubernetes/minikube) is a simplest way to do that.

1. Run `make test-e2e` to execute the test. 

1. You can also provide the following options while running the tests:
    * `PARALLEL=<value>`:  no of tests to run in parallel
           
    * `TIMEOUT=<value>` : the maximum time for which the tests could run after which the tests timeout</p>
    
    * `VERBOSE=yes` : verbose mode

Example usage: `PARALLEL=3 TIMEOUT=10m VERBOSE=yes make test-e2e`

### Running OpenShift tests

1. Install OpenShift client tools [OpenShift client tools](https://github.com/openshift/origin/releases/tag/v3.7.1)

1. For running the OpenShift tests, you need an OpenShift cluster up and running with
   [Minishift](https://github.com/minishift/minishift) or [`oc cluster up`](https://github.com/openshift/origin/blob/master/docs/cluster_up_down.md) 

1. Run `make test-e2e-os` to execute the tests. You can also provide values for PARALLEL, 
   TIMEOUT, VERBOSE option as mentioned under Kubernetes tests. 
  
