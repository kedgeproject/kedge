# Kedge - Concise Application Definition for Kubernetes

## What is Kedge?

Kedge is a simple and easy way to define and deploy applications on Kubernetes by writing very concise application definitions.

Key features and goals include:

  - _Simplicity:_ Using a simple and concise specification that is easy to understand and define.
  - _Multi-container environments:_ Define your containers, services and applications in one simple file, or abstract them into multiple files.
  - _Familiar structure:_ Using a familiar YAML structure as Kubernetes, it's easy to pick-up and understand Kedge.
  - _Built on top of Kubernetes PodSpec:_ Leverages Kuberenetes PodSpec and avoids leaky abstractions.
  - _No need to define everything:_ Define just the necessary bits and Kedge will do the rest. Kedge will interprolate and pick the best defaults for your application to run on Kubernetes.

## Project status

We are a very evolving project with high velocity, we have listed a [file reference specification](docs/file-reference.md) as well as document our RFC's and changes as [GitHub issues](https://github.com/kedgeproject/kedge/issues).

Check out our [roadmap](ROADMAP.md) as we push towards a __0.1.0__ release.

## Using Kedge

### Installation

The _best_ way to try Kedge is to download the most up-to-date binary from the master GitHub branch:

```sh
go get github.com/kedgeproject/kedge
```

### Trying it out

We have an [extensive list of examples](examples) to check out, but the simplest of them all is a [standard http example](https://raw.githubusercontent.com/kedgeproject/kedge/master/examples/simplest/httpd.yaml) with [minikube](https://github.com/kubernetes/minikube):

```yaml
name: httpd
containers:
- image: centos/httpd
services:
- name: httpd
  type: NodePort
  ports:
  - port: 8080
    targetPort: 80
```

We can now generate and deploy this example to Kubernetes:

```sh
kedge generate -f httpd.yaml | kubectl create -f -
deployment "httpd" created
service "httpd" created
```

And access it:

```sh
kubectl get po,deploy,svc
NAME                        READY     STATUS    RESTARTS   AGE
po/httpd-3617778768-ddlrs   1/1       Running   0          1m

NAME           DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
deploy/httpd   1         1         1            1           1m

NAME             CLUSTER-IP   EXTERNAL-IP   PORT(S)          AGE
svc/httpd        10.0.0.187   <nodes>       8080:31385/TCP   1m
svc/kubernetes   10.0.0.1     <none>        443/TCP          18h

minikube service httpd                           
Opening kubernetes service default/httpd in default browser...
```

Our examples range from [as simple as you can get](examples/simplest) to [every possible key you can use](examples/all). More can be found in the [/examples](examples) directory.

## Contributing 

Kedge is an evolving project and contributions are happily welcome. Feel free to open up an issue or even a PR. Read our [contributing guide](CONTRIBUTING.md) for more details. If you're interested in submitting a patch, feel free to check our [development guide](docs/development.md) as well for ease into the project. 

## License

Unless otherwise stated (ex. `/vendor` files), all code is licensed under the [Apache 2.0 License](LICENSE). Portions of the project use libraries and code from other projects, the appropriate license can be found within the code (header of the file) or root directory within the `vendor` folder.
