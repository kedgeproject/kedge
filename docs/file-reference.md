# Kedge file reference

Each file defines one micro-service, which forms one `pod` controlled by it's
controller.

A example using all the keys added in Kedge(not all keys from Kubernetes
API are included):

```yaml
name: database
controller: deployment
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
  health:
    httpGet:
      path: /
      port: 3306
services:
- name: wordpress
  ports:
  - port: 8080
    targetPort: 80
    endpoint: minikube.external/foo
volumeClaims:
- name: database
  size: 500Mi
configMaps:
- data:
    MYSQL_DATABASE: wordpress
secrets:
- name: wordpress
  data:
    MYSQL_ROOT_PASSWORD: YWRtaW4=
```

# Root level constructs

App is made of Pod Spec and added fields.
More info: https://kubernetes.io/docs/api-reference/v1.6/#podspec-v1-core


## name

`name: mariadb`

| **Type** | **Required** |
|----------|--------------|
| string   | yes          |

The name of the app or micro-service this particular file defines.

## controller

`controller: deployment`

| **Type** | **Required** |
|----------|--------------|
| string   | no           |

The Kubernetes controller of the app or micro-service this particular file
defines.

Supported controllers:
- Deployment
- Job

Default controller is **Deployment**

##### Note:
`activeDeadlineSeconds` is a conflicting field which exists in both, v1.PodSpec
and batch/v1.JobSpec, and both of these fields exist at the top level of the
Kedge spec.
So, whenever `activeDeadlineSeconds` field is set, only JobSpec is populated,
which means that `activeDeadlineSeconds` is set only for the job and not for the
pod.
To populate a pod's `activeDeadlineSeconds`, the user will have to pass this
field the long way by defining the pod exclusively under
`job.spec.template.spec.activeDeadlineSeconds`.


## labels

```yaml
labels:
  env: dev
  department: middle-tier
```

| **Type** | **Required** |
|----------|--------------|
| object   | no           |

Map of string keys and values that can be used to organize and categorize
(scope and select) objects. May match selectors of replication controllers and
services.

All the configuration created will have this label applied.
More info: http://kubernetes.io/docs/user-guide/labels

## containers

```yaml
containers:
- <containerSpec>
- <containerSpec>
```

| **Type**                                 | **Required** |
|------------------------------------------|--------------|
| array of [containerSpec](#containerSpec) | yes          |

### containerSpec

#### health

```yaml
health: <probe>
```

This is `probe` spec. Rather than defining `livenessProbe` and `readinessProbe`,
define only `health`. And then it gets copied in both in the resultant spec.
But if `health` and `livenessProbe` or `readinessProbe` are defined
simultaneously then the tool will error out.

#### envFrom

```yaml
envFrom:
- configMapRef:
    name: <string>
- secretRef:
    name: <string>
```

This is similar to the envFrom field in container which is added since Kubernetes
1.6. All the data from the ConfigMaps and Secrets referred here will be populated
as `env` inside the container.

The restriction is that the ConfigMaps and Secrets also have to be defined in the
file since there is no way to get the data to be populated.

To read more about this field from the Kubernetes upstream docs see this:
https://kubernetes.io/docs/api-reference/v1.6/#envfromsource-v1-core


## volumeClaims

```yaml
volumeClaims:
- <volume>
- <volume>
```

| **Type**                                       | **Required** |
|------------------------------------------------|--------------|
| array of [persistentVolume](#persistentVolume) | no           |


List of `volume` struct.

### persistentVolume

A `pvc` is created for each `persistentVolume`. This is PersistentVolumeClaimSpec and added 
fields. More info: https://kubernetes.io/docs/api-reference/v1.6/#persistentvolumeclaimspec-v1-core

```yaml
name: database
size: 500Mi
```

**OR**

```yaml
name: database
accessModes:
- ReadWriteOnce
resources:
  requests:
    storage: 500Mi
```

A user needs to define this list of volumes and then use it in the `volumeMounts` field in
`containers`. In the resultant output the `volumes` in `podSpec` will be populated 
automatically by the tool.

#### name

`name: database`

| **Type** | **Required** |
|----------|--------------|
| string   | yes          |

The name of the volume. This should match with the `volumeMount` defined in the
`container`.


#### size

`size: 700Mi`

| **Type** | **Required** |
|----------|--------------|
| string   | yes          |

Size of persistent volume claim to be created. Conflicts with [resources](#resources) field
so define either of those.

#### resources

```yaml
resources:
  requests:
    storage: 500Mi
```

| **Type**               | **Required** |
|------------------------|--------------|
| ResourceRequirements   | yes          |

Resources represents the minimum resources the volume should have. Conflicts with
[size](#size) field so define either of those.
More info: http://kubernetes.io/docs/user-guide/persistent-volumes#resources

#### accessModes

```yaml
accessModes:
- ReadWriteOnce
```

| **Type**        | **Required** |
|-----------------|--------------|
| array of string | no           |

AccessModes contains the desired access modes the volume should have. Defaults to 
`ReadWriteOnce`.

The access modes are:
- **`ReadWriteOnce`** – the volume can be mounted as read-write by a single node
- **`ReadOnlyMany`** – the volume can be mounted read-only by many nodes
- **`ReadWriteMany`** – the volume can be mounted as read-write by many nodes

More info: http://kubernetes.io/docs/user-guide/persistent-volumes#access-modes-1

## configMaps

```yaml
configMaps:
- <configMap>
- <configMap>
```

| **Type**                         | **Required** |
|----------------------------------|--------------|
| array of [configMap](#configMap) | no           |

### configMap

```yaml
name: string
data:
  key1: value1
  key2: value2
```

example:

```yaml
name: database
data:
  MYSQL_DATABASE: wordpress
  app_data: /etc/app/data
```

#### Name

`name: database`

| **Type** | **Required** |
|----------|--------------|
| string   | yes          |

The name of the configMap. This is optional field if only one configMap is defined, the
default name will be the app name.

#### Data

| **Type** | **Required** |
|----------|--------------|
| object   | yes          |

Data contains the configuration data. Each key must be a valid
DNS_SUBDOMAIN with an optional leading dot.

A `configMap` is created out of this configuration.

## services

```yaml
services:
- <service>
- <service>
```

| **Type**                     | **Required** |
|------------------------------|--------------|
| array of [service](#service) | no           |

### service

```yaml
name: <string>
ports:
- port: <int>
  endpoint: <URL>/<Path>
portMappings:
- <port>:<targetPort>/<protocol>
<Kubernetes Service Spec>
```

Each service is Kubernetes Service spec and added fields.
More info: https://kubernetes.io/docs/api-reference/v1.6/#servicespec-v1-core

Example:
```yaml
name: wordpress
ports:
- port: 8080
  targetPort: 80
```

Each service gets converted into a Kubernetes `service` and `ingress`es
respectively.

#### name

`name: wordpress`

| **Type** | **Required** |
|----------|--------------|
| string   | yes          |

The name of the service.

#### endpoint

`endpoint: www.mycoolapp.com/admin`

This is an added field in the Service port, which if specified an `ingress`
resource is created. The `ingress` resource name will be the same as the name
of `service`.

`endpoint` the way it is defined is can actually can be divided into
two parts the `URL` and `Path`, it is delimited by a forward slash.

#### portMappings
```yaml
portMappings:
- 8081:81/UDP
```

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

## ingresses

```yaml
ingresses:
- <ingressObject>
- <ingressObject>
```

| **Type**                                  | **Required** |
|-------------------------------------------|--------------|
| array of [ingress object](#ingressObject) | no           |


### ingressObject

```yaml
ingresses:
- name: <string>
  <Ingress Spec>
```

Example:
```yaml
name: wordpress
rules:
- host: minikube.local
  http:
    paths:
    - backend:
        serviceName: wordpress
        servicePort: 8080
      path: /
```


Each `ingress` object is Kubernetes Ingress spec and `name` field.
More info: https://kubernetes.io/docs/api-reference/v1.6/#ingressspec-v1beta1-extensions

If there is only one port and user wants to expose the service then user should
define one `ingress` with `host` atleast then the rest of the `ingress`
spec(things like `http`, etc.) will be populated for the user.

More info about Probe: https://kubernetes.io/docs/api-reference/v1.6/#probe-v1-core

#### name

`name: wordpress`

| **Type** | **Required** |
|----------|--------------|
| string   | yes          |

The name of the Ingress.

## secrets

```yaml
secrets:
- <secret>
- <secret>
```

| **Type**                         | **Required** |
|----------------------------------|--------------|
| array of [secret](#secret) | no           |

###secret

```yaml
name: string
<Kubernetes Secret Definition>
```

The Kubernetes Secret resource is being reused here.
More info: https://kubernetes.io/docs/api-reference/v1.6/#secret-v1-core

So, the Kubernetes Secret resource allows specifying the secret data as base64
encoded as well as in plaintext.
This would look in kedge as:

```yaml
secrets:
- name: <name of the secret>
  data:
    <secret data key>: <base64 encoded value of the secret data>
  stringData:
    <secret data key>: <plaintext value of the secret data>
```

example:

```yaml
secrets:
- name: wordpress
  data:
    MYSQL_ROOT_PASSWORD: YWRtaW4=
    MYSQL_PASSWORD: cGFzc3dvcmQ=
```

#### Name

`name: wordpress`

| **Type** | **Required** |
|----------|--------------|
| string   | no           |

The name of the secret.

## includeResources

```yaml
includeResources:
- <string>
- <string>
```

e.g.

```yaml
includeResources:
- ./kubernetes/cron-job.yaml
- secrets.yaml
```

This is list of files that are Kubernetes resources which can be passed to
Kubernetes directly. On these list of files Kedge won't do any processing, but
pass it to Kubernetes directly.

The file path are relative to the kedge application file.

This is one of the mechanisms to extend kedge beyond its capabilites to support
anything in the Kubernetes land.

## Complete example (deployment)

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

## Example (job)

```yaml
controller: job
name: pival
containers:
- image: perl
  command: ["perl",  "-Mbignum=bpi", "-wle", "print bpi(2000)"]
restartPolicy: Never
parallelism: 3
```


# Variables
You can use variables anywhere in the Kedge file. Variable names are enclosed
in double square brackets (`[[ variable_name ]]`).
For example `[[ IMAGE_NAME ]]` will be replaced with value of environment variable `$IMAGE_NAME`

## Example
```yaml
# nginx.yaml

name: nginx
containers:
- image: nginx:[[ NGINX_VERSION ]]
services:
- name: nginx
  ports:
  - port: 8080
    targetPort: 80
```

Now you can call Kedge and define `NGINX_VERSION` variable.
```
NGINX_VERSION=1.13 kedge apply -f nginx.yaml
```
The string `[[ NGINX_VERSION ]]` will be replaced with `1.13`.
Effecting Kedge file will look as following:
```
name: nginx
containers:
- image: nginx:1.13
services:
- name: nginx
  type: ClusterIP
  ports:
  - port: 8080
    targetPort: 80
```
