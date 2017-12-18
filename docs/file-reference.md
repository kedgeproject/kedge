# Introduction

> Using an example [httpd.yaml](https://raw.githubusercontent.com/kedgeproject/kedge/master/examples/simplest/httpd.yaml) file.

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

> Now run the apply command to deploy to Kubernetes

```sh
$ kedge apply -f httpd.yaml
deployment "httpd" created
service "httpd" created
```

> View the deployed service

```sh
$ minikube service httpd
Opening kubernetes service default/httpd in default browser...

$ kubectl describe svc httpd
Name:                   httpd
...
Endpoints:              172.17.0.4:80
...
```

> Extending Kedge as well as using a shortcut

```yaml
name: web

containers:
- image: nginx
  # Extending the Kedge file using:
  # https://v1-6.docs.kubernetes.io/docs/api-reference/v1.6/#container-v1-core
  livenessProbe:
    httpGet:
      path: /
      port: 80
    initialDelaySeconds: 20
    timeoutSeconds: 5

services:
- name: nginx
  type: NodePort
  # Using a Kedge-specific 'shortcut'
  portMappings:
  - 8080:80
```

__Note:__ This markdown file is best viewed at [kedgeproject.org/file-reference/](http://kedgeproject.org/file-reference/).

Kedge is a simple, easy and declarative way to define and deploy applications to Kubernetes by writing very concise application definitions.

It's an **extension** of Kubernetes constructs and extends many concepts of Kubernetes you're familiar with, such as PodSpec.

**Installation and Quick Start**

Installing Kedge can be found at [kedgeproject.org](http://kedgeproject.org) or alternatively, the [GitHub release page](https://github.com/kedgeproject/kedge/releases).

If you haven't used Kedge yet, we recommend using the [Quick Start](http://kedgeproject.org/quickstart/) guide, or follow the instructions within the side-bar.

**Extending Kubernetes**

Kedge introduces a simplification of Kubernetes constructs in order to make application development simple and easy to modify/deploy.

However, in many parts of Kedge, you're able to use the standard Kubernetes constructs you may already know.

For example, Kedge simplifies deployment by introducing the `health` key. However, you can still use constructs such as `readinessProbe` or `livenessProbe`.

# Kedge Keys

> All defineable Kedge keys

```yaml
name: <string>
appversion: <string>
controller: <string>
labels: <object>
containers:
  - <containerObject>
volumeClaims:
  - <persistentVolumeObject>
configMaps:
  - <configMapObject>
services:
  - <serviceObject>
ingresses:
  - <ingressObject>
routes:
  - <routeObject>
secrets:
  - <secretObject>
imageStreams:
  - <imageStreamObject>
buildConfigs:
  - <buildConfigObject>
includeResources:
  - <includeResources>
```

<aside class="notice">
Depending on the controller key selected. Each "app" (Kedge file) is an extension of that controller. Anything defineable within Kubernetes may also be defined within Kedge. See <a href="http://kedgeproject.org/file-reference/#controller">the controller table</a> for more information.
</aside>


#### Kedge specific

| Field    | Type     | Required     | Description  |
|----------|----------|--------------|--------------|
| name | string   | yes          | The name of the app or microservice this particular file defines. |
| appversion | string | no           | The version of the app or microservice this particular file defines. |
| controller | string   | no           | The Kubernetes controller of the app or microservice this particular file (default: "deployment") |
| labels | object   | no           | Map of string keys and values that can be used to organize and categorize (scope and select) objects. May match selectors of replication controllers and services. |
| containers | array of [containerObject](#containerobject) | yes          | [containerObject](#containerobject) |
| volumeClaims | array of [persistentVolumeObject](#persistentvolumeobject) | no           | [persistentVolumeObject](#persistentvolumeobject) |
| configMap | array of [configMapObject](#configmapobject) | no           | [configMapObject](#configmapobject) |
| services | array of [serviceObject](#serviceobject) | no           | [serviceObject](#serviceobject) |
| ingresses | array of [ingressObject](#ingressobject) | no           | [ingressObject](#ingressobject) |
| routes | array of [routeObject](#routeobject) | no           | [routeObject](#routeobject) |
| secrets | array of [secretObject](#secretobject) | no           | [secretObject](#secretobject) |
| imageStreams | array of [imageStreamObject](#imagestreamobject) | no           | [imageStreamObject](#imagestreamobject) |
| buildConfigs | array of [buildConfigObject](#buildconfigobject) | no           | [buildConfigObject](#buildconfigobject) |
| includeResources | array of [includeResourceObject](#includeresourceobject) | no           | [includeResourceObject](#includeresourceobject) |


## name

```yaml
name: mariadb
```

| Type     | Required     | Description  |
|----------|--------------|--------------|
| string   | yes          | The name of the app or microservice this particular file defines. |


## appversion

```yaml
appversion: 5.0.0-alpha
```

| Type     | Required     | Description  |
|----------|--------------|--------------|
| string   | no           | The version of the app or microservice this particular file defines. |


## controller

```yaml
controller: deployment
```

> Example using DeploymentConfig for OpenShift

```yaml
controller: deploymentconfig
name: httpd

containers:
- image: bitnami/nginx

services:
- name: httpd
  type: NodePort
  ports:
  - port: 8080
    targetPort: 8080

# DeploymentConfig related definitions
replicas: 2
```

> Example using Job

```yaml
controller: job
name: pival
containers:
- image: perl
  command: ["perl",  "-Mbignum=bpi", "-wle", "print bpi(2000)"]

# Job-related definitions
restartPolicy: Never
parallelism: 3
```

| Type     | Required     | Description |
|----------|--------------|-------------|
| string   | no           | The Kubernetes controller of the app or microservice this particular file defines (default: "deployment") |

Supported controllers:

| Name                 | Container Orchestrator | Extension                                                                                                                          |
|----------------------|------------------------|------------------------------------------------------------------------------------------------------------------------------------|
| Deployment (default) | Kubernetes             | An extension of [DeploymentSpec](https://v1-6.docs.kubernetes.io/docs/api-reference/v1.6/#deployment-v1beta1-apps)                 |
| Job                  | Kubernetes             | An extension of [JobSpec](https://v1-6.docs.kubernetes.io/docs/api-reference/v1.6/#job-v1-batch)                                   |
| DeploymentConfig     | OpenShift              | An extension of [DeploymentConfigSpec](https://docs.openshift.org/latest/rest_api/apis-apps/v1beta1.Deployment.html#object-schema) |


__Note on conflicting fields:__

A conflicting field may exist both within the "root" of the file as well as within the specification.

For example, `activeDeadlineSeconds` is a conflicting field which exists in both, PodSpec and as well as JobSpec, and both of these fields may be defined at the top level of the Kedge file specification

So, whenever `activeDeadlineSeconds` field is set, only JobSpec is populated, which means that `activeDeadlineSeconds` is set only for the Job and not for the Pod.

To populate a Pod's `activeDeadlineSeconds`, the user will have to pass this field the long way by defining the pod exclusively under `job.spec.template.spec.activeDeadlineSeconds`.


## labels

```yaml
labels:
  env: dev
  department: middle-tier
```

| Type | Required | Description |
|----------|--------------|-----|
| object   | no           | Map of string keys and values that can be used to organize and categorize (scope and select) objects. May match selectors of replication controllers and services. |


All the configuration created will have this label applied.
More info: [http://kubernetes.io/docs/user-guide/labels](http://kubernetes.io/docs/user-guide/labels)

## containers

```yaml
containers:
- <containerObject>
```

| Type                                 | Required | Description |
|------------------------------------------|--------------|-------|
| array of [containerObject](#containerobject) | yes          | [containerObject](#containerobject) | 


## volumeClaims

```yaml
volumeClaims:
- <volume>
```

| Type                                       | Required | Description |
|------------------------------------------------|--------------|---|
| array of [persistentVolumeObject](#persistentvolumeobject) | no           | [persistentVolumeObject](#persistentvolumeobject)|


## configMaps

```yaml
configMaps:
- <configMapObject>
```

| Type                         | Required | Description |
|----------------------------------|--------------|----|
| array of [configMapObject](#configmapobject) | no           | [configMapObject](#configmapobject) |


## services

```yaml
services:
- <service>
```

| Type                     | Required | Description |
|------------------------------|--------------|------|
| array of [serviceObject](#serviceobject) | no           | [serviceObject](#serviceobject) |


## ingresses

```yaml
ingresses:
- <ingressObject>
```

| Type                                  | Required | Description |
|-------------------------------------------|--------------|------|
| array of [ingressObject](#ingressobject) | no           | [ingressObject](#ingressobject) |



## routes

```yaml
routes:
- <routeObject>
```

| Type                                  | Required | Description |
|-------------------------------------------|--------------|------|
| array of [routeObject](#routeobject) | no           | [routeObject](#routeobject) |



## secrets

```yaml
secrets:
- <secret>
```

| Type                         | Required | Description |
|----------------------------------|--------------|------|
| array of [secretObject](#secretobject) | no           | [secretObject](#secretobject) |

## buildConfigs

```yaml
buildConfigs:
- <buildConfigObject>
```

| Type                         | Required | Description |
|----------------------------------|--------------|------|
| array of [buildConfig Object](#buildconfigobject) | no           | [buildConfig Object](#buildconfigobject) object |

## imageStreams

```yaml
imageStreams:
- <imageStreamObject>
```

| Type                         | Required | Description |
|----------------------------------|--------------|------|
| array of [imageStream Object](#imagestreamobject) | no           | [imageStream Object](#imagestreamobject) object |

## includeResources

```yaml
includeResources:
- <string>
```

> Example

```yaml
includeResources:
- ./kubernetes/cron-job.yaml
- secrets.yaml
```

This is list of files that are Kubernetes specific that can be passed to Kubernetes directly. Of these files, Kedge will not do any processing, but simply pass it to the container orchestrator.

| Type                         | Required | Description |
|----------------------------------|--------------|------|
| array of [includeResourceObject](#includeresourceobject) | no           | [includeResourceObject](#includeresourceobject) |


The file path are relative to the kedge application file.

This is one of the mechanisms to extend kedge beyond its capabilites to support
anything in the Kubernetes land.

# Objects

## containerObject

```yaml
containers:
  - <containerObject>
```

<aside class="notice">
Each "container" is a Kubernetes <a target="_blank" href="https://v1-6.docs.kubernetes.io/docs/api-reference/v1.6/#container-v1-core">Container Spec</a> with additional Kedge-specific keys.
</aside>

List of containers

| Field | Type     | Required     | Description  |
|-------|----------|--------------|--------------|
| health | string   | yes          | The name of the app or microservice this particular file defines. |

### health

```yaml
containers:
  - image: foobar
    health: <probe>
```

| Type     | Required     | Description  |
|----------|--------------|--------------|
| string   | yes          | The name of the app or microservice this particular file defines. |

This is `probe` spec. Rather than defining `livenessProbe` and `readinessProbe`,
define only `health`. And then it gets copied in both in the resultant spec.
But if `health` and `livenessProbe` or `readinessProbe` are defined
simultaneously then the tool will error out.

### Kubernetes extension

> Example extending `containers` with Kubernetes Container Spec

```yaml
name: web
containers:
- image: nginx
  # https://v1-6.docs.kubernetes.io/docs/api-reference/v1.6/#container-v1-core
  env:
  - name: WORDPRESS_DB_PASSWORD
    value: wordpress
  - name: WORDPRESS_DB_USER
    value: wordpress
  envFrom:
  - configMapRef:
      name: web
services:
- name: nginx
  type: NodePort
  ports:
  - port: 8080
    targetPort: 80
```

Anything [Container Spec](https://v1-6.docs.kubernetes.io/docs/api-reference/v1.6/#container-v1-core) from Kubernetes can be included within the Kedge file.

For example, keys such as `env` and `envFrom` are commonly used.

## persistentVolumeObject

```yaml
volumeClaims:
  - <persistentVolumeObject>
```

> An example of deploying a volume

```yaml
volumeClaims:
  - name: database
    size: 500Mi
```

> Or further specifically defining it

```yaml
volumeClaims:
  - name: database
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 500Mi
```

<aside class="notice">
Each "persistentVolumeObject" is a Kubernetes <a target="_blank" href="https://v1-6.docs.kubernetes.io/docs/api-reference/v1.6/#claim-v1-core">PersistentVolumeClaim</a> with additional Kedge-specific keys.
</aside>


| Field | Type     | Required     | Description  |
|-------|----------|--------------|--------------|
| name | string   | yes          | The name of the volume. This should match with the `volumeMount` defined in the `container`. |
| size | string   | yes          | Size of persistent volume claim to be created. Conflicts with [resources](#resources) field so define either of those. |
| resources | ResourceRequirements   | yes          | Resources represents the minimum resources the volume should have. Conflicts with [size](#size) field so define either of those. |
| accessModes | array of string | no           | AccessModes contains the desired access modes the volume should have. Defaults to `ReadWriteOnce`. |

A user needs to define this list of volumes and then use it in the `volumeMounts` field in `containers`. In the resultant output the `volumes` in `podSpec` will be populated automatically by the tool.

### name

```yaml
name: database
```

| Type | Required | Description |
|----------|--------------|-----|
| string   | yes          | The name of the volume. This should match with the `volumeMount` defined in the `container`. |


### size

```yaml
size: 700Mi
```

| Type | Required | Description |
|----------|--------------|-----|
| string   | yes          | Size of persistent volume claim to be created. Conflicts with [resources](#resources) field so define either of those. |

### resources

```yaml
resources:
  requests:
    storage: 500Mi
```

| Type               | Required | Description |
|------------------------|--------------|-----|
| ResourceRequirements   | yes          | Resources represents the minimum resources the volume should have. Conflicts with [size](#size) field so define either of those. |

More info: http://kubernetes.io/docs/user-guide/persistent-volumes#resources

### accessModes

```yaml
accessModes:
- ReadWriteOnce
```

| Type        | Required | Description
|-----------------|--------------|------|
| array of string | no           | AccessModes contains the desired access modes the volume should have. Defaults to `ReadWriteOnce`. |

The access modes are:

- `ReadWriteOnce` – the volume can be mounted as read-write by a single node
- `ReadOnlyMany` – the volume can be mounted read-only by many nodes
- `ReadWriteMany` – the volume can be mounted as read-write by many nodes

More info: http://kubernetes.io/docs/user-guide/persistent-volumes#access-modes-1

### Kubernetes extension

> Example extending `volumesClaims` with Kubernetes PersistentVolumeClaim Spec

```yaml
name: database
containers:
- image: mariadb:10
  env:
  - name: MYSQL_ROOT_PASSWORD
    value: rootpasswd
  - name: MYSQL_DATABASE
    value: wordpress
  - name: MYSQL_USER
    value: wordpress
  - name: MYSQL_PASSWORD
    value: wordpress
  volumeMounts:
  - name: database
    mountPath: /var/lib/mysql
services:
- name: database
  ports:
  - port: 3306
volumeClaims:
- name: database
  size: 500Mi
  # https://v1-6.docs.kubernetes.io/docs/api-reference/v1.6/#persistentvolumeclaim-v1-core
  persistentVolumeReclaimPolicy: Recycle
  storageClassName: slow
  mountOptions:
    - hard
    - nfsvers=4.1
```

Anything [PersistentVolumeClaim Spec](https://v1-6.docs.kubernetes.io/docs/api-reference/v1.6/#persistentvolumeclaim-v1-core) from Kubernetes can be included within the Kedge file.

## configMapObject

```yaml
configMaps:
  - <configMapObject>
```

> Example

```yaml
configMaps:
- name: database
  data:
    MYSQL_DATABASE: wordpress
    app_data: /etc/app/data
```

<aside class="notice">
Each "configMapObject" is a Kubernetes <a target="_blank" href="https://v1-6.docs.kubernetes.io/docs/api-reference/v1.6/#configmap-v1-core">ConfigMap Spec</a> with additional Kedge-specific keys.
</aside>

| Field | Type     | Required     | Description  |
|-------|----------|--------------|--------------|
| name | string   | yes          | The name of the configMapObject. This is optional field if only one configMapObject is defined, the default name will be the app name. |
| data | object   | yes          | Data contains the configuration data. Each key must be a valid DNS_SUBDOMAIN with an optional leading dot. |

### Name

```yaml
name: database
```

| Type | Required | Description |
|----------|--------------|-----|
| string   | yes          | The name of the configMapObject. This is optional field if only one configMapObject is defined, the default name will be the app name. |

### Data

```yaml
data:
  key: value
```

| Type | Required | Description |
|----------|--------------|--------|
| object   | yes          | Data contains the configuration data. Each key must be a valid DNS_SUBDOMAIN with an optional leading dot. |

A `configMapObject` is created out of this configuration.

### Kubernetes extension

> Example extending `configMapObjects` with Kubernetes ConfigMap Spec

```yaml
name: database
containers:
- image: mariadb:10
  env:
  - name: MYSQL_ROOT_PASSWORD
    value: rootpasswd
  - name: MYSQL_DATABASE
    valueFrom:
      configMapKeyRef:
        key: MYSQL_DATABASE
        name: database
  - name: MYSQL_USER
    value: wordpress
  - name: MYSQL_PASSWORD
    value: wordpress
services:
- name: database
  ports:
  - port: 3306
configMaps:
- data:
  # https://v1-6.docs.kubernetes.io/docs/api-reference/v1.6/#configmap-v1-core
    MYSQL_DATABASE: wordpress
```

Anything [ConfigMap Spec](https://v1-6.docs.kubernetes.io/docs/api-reference/v1.6/#configmap-v1-core) from Kubernetes can be included. **Note:** Since Kedge already implents "data" in ConfigMaps no other keys are available to be added.

## serviceObject

```yaml
services:
  - <serviceObject>
```

> Example

```yaml
services:
- name: wordpress
  ports:
  - port: 8080
    targetPort: 80
  portMappings:
  - 90:9090/tcp
```

<aside class="notice">
Each "service" is a Kubernetes <a target="_blank" href="https://v1-6.docs.kubernetes.io/docs/api-reference/v1.6/#service-v1-core">Service Spec</a> with additional Kedge-specific keys.
</aside>

| Field | Type     | Required     | Description  |
|-------|----------|--------------|--------------|
| name  | string   | yes          | The name of the service. |
| endpoint | string   | no | The endpoint of the service. |
| portMappings | array of "port" | no |  Array of ports. Ex. `80:8080/tcp` |

More info: [https://v1-6.docs.kubernetes.io/docs/api-reference/v1.6/#servicespec-v1-core](https://v1-6.docs.kubernetes.io/docs/api-reference/v1.6/#servicespec-v1-core)

Each service gets converted into a Kubernetes `service` and `ingress` respectively.

### name

```yaml
name: wordpress
```

| Type | Required | Description |
|----------|--------------|------|
| string   | yes          | The name of the service. |

### endpoint

```yaml
endpoint: www.mycoolapp.com/admin
```

| Type | Required | Description |
|----------|--------------|------|
| string   | yes          | The endpoint of the service. |

This is an added field in the Service port, which if specified an `ingress`
resource is created. The `ingress` resource name will be the same as the name
of `service`.

`endpoint` the way it is defined is can actually can be divided into
two parts the `URL` and `Path`, it is delimited by a forward slash.

### portMappings

```yaml
portMappings:
- 8081:81/UDP
```

| Type | Required | Description |
|----------|--------------|------|
| array of "port" | yes          |  Array of ports. Ex. `80:8080/tcp` |


`portMappings` is an added field to ServiceSpec.
This lets us set the port, targetPort and the protocol for a service in a single line. This is parsed and converted to a Kubernetes ServicePort object.

`portMappings` is an array of `port:targetPort/protocol` definitions, so the syntax looks like -

```yaml
portMappings:
- <port:targetPort/protocol>
- <port:targetPort/protocol>
```

The only mandatory part to specify in a portMapping is "port".
There are 4 possible cases here

- When only `port` is specified - `targetPort` is set to `port` and protocol is set to `TCP`
- When `port:targetPort` is specified - protocol is set to `TCP`
- When `port/protocol` is specified - `targetPort` is set to `port`
- When `port:targetPort/protocol` is specified - no auto population is done since all values are provided

Find a working example using `portMappings` field [here](https://github.com/kedgeproject/kedge/tree/master/docs/examples/portMappings/httpd.yaml)

### Kubernetes extension

> Example extending `service` with Kubernetes Service Spec

```yaml
name: httpd
containers:
- image: centos/httpd
services:
- name: httpd
  # https://v1-6.docs.kubernetes.io/docs/api-reference/v1.6/#servicespec-v1-core
  ports:
  - port: 8080
    targetPort: 80
  type: NodePort
```

Anything [Service Spec](https://v1-6.docs.kubernetes.io/docs/api-reference/v1.6/#servicespec-v1-core) from Kubernetes can be included within the Kedge file.

For example, keys such as `image` and `ports` are commonly used.

## ingressObject

```yaml
ingresses:
- <ingress>
```

<aside class="notice">
Each "ingress" is a Kubernetes <a target="_blank" href="https://v1-6.docs.kubernetes.io/docs/api-reference/v1.6/#ingressspec-v1beta1-extensions">Ingress Spec</a> with additional Kedge-specific keys.
</aside>


| Type                         | Required | Description |
|----------------------------------|--------------|------|
| name | string   | yes          | The name of the Ingress |

If there is only one port and user wants to expose the service then user should define one `ingress` with `host` atleast then the rest of the `ingress` spec (things like `http`, etc.) will be populated for the user.


### name

```yaml
name: wordpress
```

| Type | Required | Description |
|----------|--------------|-------|
| string   | yes          | The name of the Ingress. |

### Kubernetes extension

> Example extending `ingresses` with Kubernetes Ingress Spec

```yaml
ingresses:
- name: wordpress
  # https://v1-6.docs.kubernetes.io/docs/api-reference/v1.6/#ingressspec-v1beta1-extensions
  rules:
  - host: minikube.local
    http:
      paths:
      - backend:
          serviceName: wordpress
          servicePort: 8080
        path: /
```

Anything [Ingress Spec](https://v1-6.docs.kubernetes.io/docs/api-reference/v1.6/#ingressspec-v1beta1-extensions) from Kubernetes can be included within the Kedge file.

## routeObject

```yaml
routes:
- <route>
```

> Example

```yaml
name: webroute
to:
  kind: Service
  name: httpd
```

<aside class="notice">
Each "route" is an OpenShift <a target="_blank" href="https://docs.openshift.org/latest/rest_api/apis-route.openshift.io/v1.Route.html#object-schema">Route Spec</a> with additional Kedge-specific keys.
</aside>

| Type                         | Required | Description |
|----------------------------------|--------------|------|
| name | string   | yes          | The name of the Route |

### name

```yaml
name: wordpress
```

| Type | Required | Description |
|----------|--------------|-------|
| string   | yes          | The name of the Route |

### OpenShift extension

> Example extending `routes` with OpenShift Route Spec

```yaml
routes:
- name: httpd
  # https://docs.openshift.org/latest/rest_api/apis-route.openshift.io/v1.Route.html#object-schema
  to:
    kind: Service
    name: httpd
```

Anything [Route Spec](https://docs.openshift.org/latest/rest_api/apis-route.openshift.io/v1.Route.html#object-schema) from OpenShift can be included within the Kedge file.

## secretObject

```yaml
secrets:
  - <secret>
```

<aside class="notice">
Each "secret" is a Kubernetes <a target="_blank" href="https://v1-6.docs.kubernetes.io/docs/api-reference/v1.6/#envvarsource-v1-core">EnvVarSource Spec</a> with additional Kedge-specific keys.
</aside>

### Name

```yaml
name: wordpress
```

| Type | Required | Description |
|----------|--------------|-----|
| string   | no           | The name of the secret |

### Kubernetes extension

> Example extending `service` with Kubernetes Service Spec

```yaml
secrets:
- name: wordpress
  data:
    # https://v1-6.docs.kubernetes.io/docs/api-reference/v1.6/#envvarsource-v1-core
    # Encoded in base64
    MYSQL_ROOT_PASSWORD: YWRtaW4=
    MYSQL_PASSWORD: cGFzc3dvcmQ=
```

Anything [EnvVarSource Spec](https://v1-6.docs.kubernetes.io/docs/api-reference/v1.6/#envvarsource-v1-core) from Kubernetes can be included within the Kedge file.

## buildConfigObject

```yaml
buildConfigs:
- <buildConfigObject>
```

> Example

```yaml
name: rubybc
triggers:
- type: "ImageChange"
source:
  type: "Git"
  git:
    uri: "https://github.com/openshift/ruby-hello-world"
strategy:
  type: "Source"
  sourceStrategy:
    from:
      kind: "ImageStreamTag"
      name: "ruby-22-centos7:latest"
output:
  to:
    kind: "ImageStreamTag"
    name: "origin-ruby-sample:latest"
```

<aside class="notice">
Each "buildConfigObject" is an OpenShift <a target="_blank" href="https://docs.openshift.org/latest/rest_api/apis-build.openshift.io/v1.BuildConfig.html#object-schema">BuildConfig Spec</a> with additional Kedge-specific keys.
</aside>


| Type                         | Required | Description |
|----------------------------------|--------------|------|
| name | string   | yes          | The name of the BuildConfig |


### name

```yaml
name: wordpress
```

| Type | Required | Description |
|----------|--------------|-------|
| string   | yes          | The name of the BuildConfig |

### OpenShift extension

> Example extending `buildConfigs` with OpenShift BuildConfig Spec

```yaml
name: ruby
buildConfigs:
- triggers:
  - type: "ImageChange"
  source:
    type: "Git"
    git:
      uri: "https://github.com/openshift/ruby-hello-world"
  strategy:
    type: "Source"
    sourceStrategy:
      from:
        kind: "ImageStreamTag"
        name: "ruby-22-centos7:latest"
  output:
    to:
      kind: "ImageStreamTag"
      name: "origin-ruby-sample:latest"
```

## imageStreamObject

```yaml
imageStreams:
- <imageStreamObject>
```

> Example

```yaml
name: rubyapp
imageStreams:
- name: rubystream
  dockerImageRepository: "docker.io/openshift/ruby-20-centos7"
```

<aside class="notice">
Each "imageStreamObject" is an OpenShift <a target="_blank" href="https://docs.openshift.org/latest/rest_api/apis-image.openshift.io/v1.ImageStream.html#object-schema">ImageStream Spec</a> with additional Kedge-specific keys.
</aside>


| Type                         | Required | Description |
|----------------------------------|--------------|------|
| name | string   | yes          | The name of the ImageStream |


### name

```yaml
name: wordpress
```

| Type | Required | Description |
|----------|--------------|-------|
| string   | yes          | The name of the ImageStream |

### OpenShift extension

> Example extending `imageStreams` with OpenShift ImageStream Spec

```yaml
name: webapp
imageStreams:
- tags:
  - from:
      kind: DockerImage
      name: centos/httpd-24-centos7:latest
    name: "2.4"
```

## includeResourceObject

```yaml
includeResources:
- <string>
```

> Example

```yaml
includeResources:
- ./kubernetes/cron-job.yaml
- secrets.yaml
```

Including external resources.

| Type | Required | Description |
|----------|--------------|-----|
| string   | no           | File location of the Kubernetes resource |

# Variables

> Example using local environment variables

```yaml
name: nginx
containers:
- image: nginx:[[ NGINX_VERSION ]]
services:
- name: nginx
  ports:
  - port: 8080
    targetPort: 80
```

> Using the variables on the command line

```sh
NGINX_VERSION=1.13 kedge apply -f nginx.yaml
```

You can use variables anywhere in the Kedge file. Variable names are enclosed in double square brackets (`[[ variable_name ]]`). For example `[[ IMAGE_NAME ]]` will be replaced with value of environment variable `$IMAGE_NAME`.

> Specify default value of variable,

```yaml
name: foo
image: foo/bar:[[ TAG:latest ]]
```

You can specify a default value of a variable if it is not set. 
For example: `[[ TAG:latest ]]`
If `TAG` variable is not set, `latest` will be used.


# Controllers

There are three defineable controllers within Kedge:

- Deployment (Kubernetes) (Default)
- Job (Kubernetes)
- DeploymentConfig (OpenShift)

Some controllers such as DeploymentConfig are only usable with OpenShift.

## Deployment

```yaml
name: database
containers:
- image: mariadb:10
  envFrom:
  - configMapRef:
      name: database
  - secretRef:
      name: wordpress
  volumeMounts:
  - name: database
    mountPath: /var/lib/mysql
  livenessProbe:
    httpGet:
      path: /
      port: 3306
  readinessProbe:
    exec:
      command:
      - mysqladmin
      - ping
    initialDelaySeconds: 5
    timeoutSeconds: 1
services:
- name: wordpress
  expose: true
  ports:
  - port: 8080
    targetPort: 80
    endpoint: minikube.external/foo
ingresses:
- name: pseudo-wordpress
  rules:
  - host: minikube.local
    http:
      paths:
      - backend:
          serviceName: wordpress
          servicePort: 8080
        path: /
volumeClaims:
- name: database
  size: 500Mi
  accessModes:
  - ReadWriteOnce
configMaps:
- data:
    MYSQL_DATABASE: wordpress
secrets:
- name: wordpress
  data:
    MYSQL_ROOT_PASSWORD: YWRtaW4=
```

## Job

```yaml
controller: job
name: pival
containers:
- image: perl
  command: ["perl",  "-Mbignum=bpi", "-wle", "print bpi(2000)"]
restartPolicy: Never
parallelism: 3
```

**Note**: If no `restartPolicy` is provided it defaults to `OnFailure`.

## Deployment Config

```yaml
controller: deploymentconfig
name: httpd
replicas: 2
containers:
- image: centos/httpd
services:
- name: httpd
  type: NodePort
  ports:
  - port: 8080
    targetPort: 80
```
