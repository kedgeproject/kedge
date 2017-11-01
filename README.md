# Kedge : Simple, Concise & Declarative Kubernetes Applications

[![Semaphore Build Status Widget]][Semaphore Build Status] [![Travis Build Status Widget]][Travis Build Status] [![Coverage Status Widget]][Coverage Status] [![GoDoc Widget]][GoDoc] [![GoReportCard Widget]][GoReportCardResult] [![Slack Widget]][Slack] 

## What is Kedge?

![logo](/docs/images/logo.png)

Kedge is a simple, easy and declarative way to define and deploy applications to Kubernetes by writing very concise application definitions.

Key features and goals include:

  - __Declarative:__ Declarative definitions specifying developer's intent.
  - __Simplicity:__ Using a simple and concise specification that is easy to understand and define.
  - __Concise:__ Define just the necessary bits and Kedge will do the rest. Kedge will interprolate and pick the best defaults for your application to run on Kubernetes.
  - __Multi-container environments:__ Define your containers, services and applications in one simple file, or split them into multiple files.
  - __Familiar structure:__ Using a familiar YAML structure as Kubernetes, it's easy to pick-up and understand Kedge.
  - __Built on top of Kubernetes Pod definition:__ Leverages Kuberenetes Pod definition (PodSpec) and avoids leaky abstractions.


## Project status

We are a very evolving project with high velocity, we have listed a [file reference specification](docs/file-reference.md) as well as document our RFC's and changes as [GitHub issues](https://github.com/kedgeproject/kedge/issues).

Check out our [roadmap](ROADMAP.md) as we push towards a __0.4.0__ release.

## Using Kedge

### Installing

Kedge is released via GitHub on a three-week cycle, you can see all current releases on the [GitHub release page](https://github.com/kedgeproject/kedge/releases).

__Linux and macOS:__

```sh
# Linux
curl -L https://github.com/kedgeproject/kedge/releases/download/v0.4.0/kedge-linux-amd64 -o kedge

# macOS
curl -L https://github.com/kedgeproject/kedge/releases/download/v0.4.0/kedge-darwin-amd64 -o kedge

chmod +x kedge
sudo mv ./kedge /usr/local/bin/kedge
```

__Windows:__

Download from [GitHub](https://github.com/kedgeproject/kedge/releases/download/v0.4.0/kedge-windows-amd64.exe) and add the binary to your PATH.

### Installing the latest binary (master)

You can download latest binary (built on each master PR merge) for [Linux (amd64)][Bintray Latest Linux], [macOS (darwin)][Bintray Latest macOS] or [Windows (amd64)][Bintray Latest Windows] from [Bintray](https://bintray.com):

__Linux and macOS:__

```sh
# Linux 
curl -L https://dl.bintray.com/kedgeproject/kedge/latest/kedge-linux-amd64 -o kedge

# macOS
curl -L https://dl.bintray.com/kedgeproject/kedge/latest/kedge-darwin-amd64 -o kedge

chmod +x kedge
sudo mv ./kedge /usr/local/bin/kedge
```

__Windows:__

Download from [Bintray](https://dl.bintray.com/kedgeproject/kedge/latest/kedge-windows-amd64.exe) and add the binary to your PATH.

__You can also download and build Kedge via Go:__

```sh
go get github.com/kedgeproject/kedge
```

### Trying it out

We have an [extensive list of examples](docs/examples) to check out, but the simplest of them all is a [standard http example](https://raw.githubusercontent.com/kedgeproject/kedge/master/docs/examples/simplest/httpd.yaml) with [minikube](https://github.com/kubernetes/minikube):

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

We can now create this example on the Kubernetes cluster:

```sh
$ kedge create -f httpd.yaml
deployment "httpd" created
service "httpd" created
```

And access it:

```sh
$ kubectl get po,deploy,svc
NAME                        READY     STATUS    RESTARTS   AGE
po/httpd-3617778768-ddlrs   1/1       Running   0          1m

NAME           DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
deploy/httpd   1         1         1            1           1m

NAME             CLUSTER-IP   EXTERNAL-IP   PORT(S)          AGE
svc/httpd        10.0.0.187   <nodes>       8080:31385/TCP   1m
svc/kubernetes   10.0.0.1     <none>        443/TCP          18h

$ minikube service httpd
Opening kubernetes service default/httpd in default browser...
```

Our examples range from [as simple as you can get](docs/examples/simplest) to [every possible key you can use](docs/examples/all). More can be found in the [docs/examples](docs/examples) directory.

## Community, Discussion, Contribution, and Support

__Contributing:__ Kedge is an evolving project and contributions are happily welcome. Feel free to open up an issue or even a PR. Read our [contributing guide](CONTRIBUTING.md) for more details. If you're interested in submitting a patch, feel free to check our [development guide](docs/development.md) as well for ease into the project.

__Chat (Slack):__ We're fairly active on [Slack](https://kedgeproject.slack.com#kedge). You can invite yourself at [slack.kedgeproject.org](http://slack.kedgeproject.org).

## License

Unless otherwise stated (ex. `/vendor` files), all code is licensed under the [Apache 2.0 License](LICENSE). Portions of the project use libraries and code from other projects, the appropriate license can be found within the code (header of the file) or root directory within the `vendor` folder.

[Semaphore Build Status]: https://semaphoreci.com/cdrage/kedge
[Semaphore Build Status Widget]: https://semaphoreci.com/api/v1/cdrage/kedge/branches/master/badge.svg
[Travis Build Status]: https://travis-ci.org/kedgeproject/kedge
[Travis Build Status Widget]: https://travis-ci.org/kedgeproject/kedge.svg?branch=master
[Coverage Status Widget]: https://coveralls.io/repos/github/kedgeproject/kedge/badge.svg?branch=master
[Coverage Status]: https://coveralls.io/github/kedgeproject/kedge?branch=master
[GoDoc]: https://godoc.org/github.com/kedgeproject/kedge
[GoDoc Widget]: https://godoc.org/github.com/kedgeproject/kedge?status.svg
[Slack]: http://slack.kedgeproject.org
[Slack Widget]: https://s3.eu-central-1.amazonaws.com/ngtuna/join-us-on-slack.png
[Bintray Latest Linux]:https://dl.bintray.com/kedgeproject/kedge/latest/kedge-linux-amd64
[Bintray Latest macOS]:https://dl.bintray.com/kedgeproject/kedge/latest/kedge-darwin-amd64
[Bintray Latest Windows]:https://dl.bintray.com/kedgeproject/kedge/latest/kedge-windows-amd64.exe
[GoReportCard Widget]: https://goreportcard.com/badge/github.com/kedgeproject/kedge
[GoReportCardResult]: https://goreportcard.com/report/github.com/kedgeproject/kedge
