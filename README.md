# Kedge : Simple, Concise & Declarative Kubernetes Applications

[![Semaphore Build Status Widget]][Semaphore Build Status] [![Travis Build Status Widget]][Travis Build Status] [![Coverage Status Widget]][Coverage Status] [![GoDoc Widget]][GoDoc] [![GoReportCard Widget]][GoReportCardResult]

![logo](/docs/images/logo.png)

----

Kedge is a simple, easy and declarative way to define and deploy applications to Kubernetes by writing very concise application definitions.

Why do people love Kedge?

  - __An extension of Kubernetes definitions:__ Use pre-existing definitions such as Pod or Container within your artifact file.
  - __Shortcuts:__ Reduce your file-size and definitions by using intuitive Kedge shortcuts.
  - __Concise:__ Define just the necessary bits, Kedge will interpolate and pick the best defaults for your application to run on Kubernetes.
  - __Define all your containers in one place:__ Define your containers, services and applications in one simple file, or split them into multiple files.

----

## To start using Kedge

### Installing

Kedge is released via GitHub on a three-week cycle, you can see all current releases on the [GitHub release page](https://github.com/kedgeproject/kedge/releases).

__Linux and macOS:__

```sh
# Linux
curl -L https://github.com/kedgeproject/kedge/releases/download/v0.9.0/kedge-linux-amd64 -o kedge

# macOS
curl -L https://github.com/kedgeproject/kedge/releases/download/v0.9.0/kedge-darwin-amd64 -o kedge

chmod +x kedge
sudo mv ./kedge /usr/local/bin/kedge
```

__Windows:__

Download from [GitHub](https://github.com/kedgeproject/kedge/releases/download/v0.9.0/kedge-windows-amd64.exe) and add the binary to your PATH.

A more thorough installation guide is [also available](http://kedgeproject.org/installation).

### Using Kedge

Try our [quick start guide](http://kedgeproject.org/quickstart/).

Go through our [file reference](http://kedgeproject.org/file-reference).

Then go further with our pre-existing [examples](https://github.com/kedgeproject/kedge/tree/master/examples).

## Community, Discussion, Contribution, and Support

__Contributing:__ Kedge is an evolving project and contributions are happily welcome. Feel free to open up an issue or even a PR. Read our [contributing guide](CONTRIBUTING.md) for more details. A thorough [development guide](http://kedgeproject.org/development/) is available if you're interested in contributing to Kedge.

__Chat (Slack):__ We're fairly active on [Slack](https://kedgeproject.slack.com#kedge). You can invite yourself at [slack.kedgeproject.org](http://slack.kedgeproject.org).

### License

Unless otherwise stated (ex. `/vendor` files), all code is licensed under the [Apache 2.0 License](LICENSE). Portions of the project use libraries and code from other projects, the appropriate license can be found within the code (header of the file) or root directory within the `vendor` folder.

[Semaphore Build Status]: https://semaphoreci.com/cdrage/kedge
[Semaphore Build Status Widget]: https://semaphoreci.com/api/v1/cdrage/kedge/branches/master/badge.svg
[Travis Build Status]: https://travis-ci.org/kedgeproject/kedge
[Travis Build Status Widget]: https://travis-ci.org/kedgeproject/kedge.svg?branch=master
[Coverage Status Widget]: https://coveralls.io/repos/github/kedgeproject/kedge/badge.svg?branch=master
[Coverage Status]: https://coveralls.io/github/kedgeproject/kedge?branch=master
[GoDoc]: https://godoc.org/github.com/kedgeproject/kedge
[GoDoc Widget]: https://godoc.org/github.com/kedgeproject/kedge?status.svg
[GoReportCard Widget]: https://goreportcard.com/badge/github.com/kedgeproject/kedge
[GoReportCardResult]: https://goreportcard.com/report/github.com/kedgeproject/kedge
