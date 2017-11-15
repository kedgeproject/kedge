---
layout: default
---

# Simplifying how you define Kubernetes artifacts

## Use Kedge to deploy applications with sensible defaults

### What's Kedge?

Kedge is a simple, easy and declarative way to define and deploy applications to Kubernetes by writing very concise application definitions.

Why do people love Kedge?

  - __An extension of Kubernetes definitions:__ Use pre-existing definitions such as Pod or Container within your artifact file.
  - __Shortcuts:__ Reduce your file-size and definitions by using intuitive Kedge shortcuts.
  - __Concise:__ Define just the necessary bits, Kedge will interpolate and pick the best defaults for your application to run on Kubernetes.
  - __Define all your containers in one place:__ Define your containers, services and applications in one simple file, or split them into multiple files.

### Avoid writing long artifact files, deploy an application straight to a Kubernetes cluster

```yaml
name: httpd

containers:
- image: centos/httpd

services:
- name: httpd
  type: LoadBalancer
  portMappings: 
    - 8080:80
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
