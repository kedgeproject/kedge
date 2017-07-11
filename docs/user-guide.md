# User Guide

- CLI
  - [`deploy`](#kedge-deploy)
  - [`generate`](#kedge-generate)
  - [`version`](#kedge-version)

## `kedge deploy`

Deploy directly to Kubernetes without creating the artifacts. Internally, Kedge will generate the artifacts and then deploy it using the `kubectl` command.

### Deploy to Kubernetes

```sh
kedge deploy -f httpd.yaml
deployment "httpd" created
service "httpd" created
```

## `kedge generate`

Generate Kubernetes artifacts based upon your Kedge YAML file, see our [examples](/examples) or the [file reference](/docs/file-reference.md) on how to create said file.

In these examples, we use the [simplest of examples](/examples/simplest/httpd.yaml)

### Convert to Kubernetes artifacts

```sh
kedge generate -f httpd.yaml
---
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  creationTimestamp: null
  labels:
    app: httpd
  name: httpd
spec:
  strategy: {}
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: httpd
      name: httpd
    spec:
      containers:
      - image: centos/httpd
        name: httpd
        resources: {}
status: {}
---
apiVersion: v1
kind: Service
metadata:
  creationTimestamp: null
  labels:
    app: httpd
  name: httpd
spec:
  ports:
  - port: 8080
    targetPort: 80
  selector:
    app: httpd
  type: NodePort
status:
  loadBalancer: {}
```

### Generate and deploy directly to Kubernetes

Generation commands can also be "piped" to Kubernetes

```sh
kedge generate -f httpd.yaml | kubectl create -f -
deployment "httpd" created
service "httpd" created
```

## `kedge version`

Outputs the current Kedge version

### Version

```sh
kedge version
```
