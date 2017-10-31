---
layout: default
---

# Simplifying how you define Kubernetes artifacts

## Use Kedge to deploy applications with sensible defaults

### What's Kedge?

Kedge is a simple, easy and declarative way to define and deploy applications to Kubernetes by writing very concise application definitions.

Why do people love Kedge?

  - __Declarative:__ Declarative definitions specifying developer's intent.
  - __Simplicity:__ Using a simple and concise specification that is easy to understand and define.
  - __Concise:__ Define just the necessary bits and Kedge will do the rest. Kedge will interprolate and pick the best defaults for your application to run on Kubernetes.
  - __Multi-container environments:__ Define your containers, services and applications in one simple file, or split them into multiple files.
  - __Familiar structure:__ Using a familiar YAML structure as Kubernetes, it's easy to pick-up and understand Kedge.
  - __Built on top of Kubernetes Pod definition:__ Leverages Kuberenetes Pod definition (PodSpec) and avoids leaky abstractions.

### Avoid writing long artifact files, deploy an application straight to a Kubernetes cluster

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
{: .demo-code }

![Demo Gif](/img/demo.gif)
{: .demo-gif }

View our [file reference](/file-reference) for a complete overview on what Kedge can do.

### Install and deploy on Linux, macOS or Windows

Install Kedge with our simple binary!

```sh
# Linux
curl -L https://github.com/kedgeproject/kedge/releases/download/v0.4.0/kedge-linux-amd64 -o kedge

# macOS
curl -L https://github.com/kedgeproject/kedge/releases/download/v0.4.0/kedge-darwin-amd64 -o kedge

chmod +x kedge
sudo mv ./kedge /usr/local/bin/kedge
```

For Windows users, download from the [GitHub release](https://github.com/kedgeproject/kedge/releases/download/v0.4.0/kedge-windows-amd64.exe) and add the binary to your PATH.

### Pick from an example and see what Kedge is all about

Choose from multiple Kedge examples to deploy:

- [Wordpress](https://github.com/kedgeproject/kedge/tree/master/examples/wordpress)
- [GitLab](https://github.com/kedgeproject/kedge/tree/master/examples/gitlab)
- [Kubernetes Guestbook](https://github.com/kedgeproject/kedge/tree/master/examples/guestbook-demo)
