# User Guide

- CLI
  - [`create`](#kedge-create)
  - [`generate`](#kedge-generate)
  - [`delete`](#kedge-delete)
  - [`version`](#kedge-version)  
  - [`init`](#kedge-init)

## `kedge create`

Deploy directly to Kubernetes without creating the artifacts. Internally, Kedge will generate the artifacts and then create it using the `kubectl` command.

### Deploy to Kubernetes

```sh
$ kedge create -f httpd.yaml
deployment "httpd" created
service "httpd" created
```

## `kedge generate`

Generate Kubernetes artifacts based upon your Kedge YAML file, see our [examples](/examples) or the [file reference](/docs/file-reference.md) on how to create said file.

In these examples, we use the [simplest of examples](/examples/simplest/httpd.yaml)

### Convert to Kubernetes artifacts

```sh
$ kedge generate -f httpd.yaml
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
$ kedge generate -f httpd.yaml | kubectl create -f -
deployment "httpd" created
service "httpd" created
```

## `kedge delete`

Deletes Kubernetes artifacts

```sh
$ kedge delete -f httpd.yaml
deployment "httpd" deleted
service "httpd" deleted
```
## `kedge version`

Outputs the current Kedge version

### Version

```sh
kedge version
```

## `kedge init`

Initialize kedge file

```sh
kedge init --file kedge.yml --name web --image centos/httpd --port 80
```

## `kedge build`

### Build

Build container image with image name `username/myapp:version`

```console
kedge build -i username/myapp:version
``` 

Here you might wanna replace the `username` with container image registry URL and then
`username`.


### Build & Push

```console
kedge build -i username/myapp:version -p
``` 

If you want to build image and also push it to the registry then use the flag `-p`.

**Note**: You should have access to push image to that container registry. Read more about
`docker login` on official [docs](https://docs.docker.com/engine/reference/commandline/login/).

### Dockerfile and context are different

If you have a file structure like this, where your `Dockerfile` resides in a directory and
your code context is different.

```console
$ tree
.
├── main.go
├── parsego.go
├── scripts
│   ├── Dockerfile
│   ├── entrypoint.sh
│   └── k8s-release
...
```

To do builds in environment like above, run following command:

```console
$ kedge build -i surajd/json-schema -c . -f scripts/Dockerfile 
INFO[0000] Building image 'surajd/json-schema' from directory 'json-schema-generator' 
INFO[0001] Image 'surajd/json-schema' from directory 'json-schema-generator' built successfully
```

### Build in minikube/minishift

If you are running Kubernetes in a environment like minikube or minishift run following
command before running this build command:

```console
eval $(minikube docker-env)
```
