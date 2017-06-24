# Kapp file reference

Each file defines one micro-service, which forms one `pod` controlled by it's
controller(right now the default controller is `deployment`).


A example using all the keys added in Kapp(not all keys from Kubernetes
API are included):

```yaml
name: database
containers:
- image: mariadb:10
  env:
  - name: MYSQL_ROOT_PASSWORD
    valueFrom:
      secretKeyRef:
        name: wordpress
        key: MYSQL_ROOT_PASSWORD
  - name: MYSQL_DATABASE
    valueFrom:
      configMapKeyRef:
        key: MYSQL_DATABASE
        name: database
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
persistentVolumes:
- name: database
  size: 500Mi
configData:
  MYSQL_DATABASE: wordpress
```

# Rootlevel constructs

App is made of Pod Spec and added fields.
More info: https://kubernetes.io/docs/api-reference/v1.6/#podspec-v1-core


## name

`name: mariadb`

| **Type** | **Required** |
|----------|--------------|
| string   | yes          |

The name of the app or micro-service this particular file defines.

## replicas

`replicas: 4`

| **Type** | **Required** |
|----------|--------------|
| integer  | no           |

Number of desired pods. This is a pointer to distinguish between explicit zero
and not specified. Defaults to 1. The valid value can only be a positive number.
This is an optional field.


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

This is `probe` spec. Rather than defining `livenessProbe` and `readinessProbe`, define only
`health`. And then it gets copied in both in the resultant spec.
But if `health` and `livenessProbe` or `readinessProbe` are defined simultaneously then the
tool will error out.

## persistentVolumes

```yaml
persistentVolumes:
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

## configData

```yaml
configData:
  MYSQL_DATABASE: wordpress
  app_data: /etc/app/data
```

| **Type** | **Required** |
|----------|--------------|
| object   | no           |

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

Each service gets generated into a Kubernetes service and ingresses respectively.

#### name

`name: wordpress`

| **Type** | **Required** |
|----------|--------------|
| string   | yes          |

The name of the service.

#### endpoint

`endpoint: www.mycoolapp.com/admin`

This is an added field in the Service port, which if specified a ingress
resource is created. The ingress resource name will be the same as the name
of service.

The complete `endpoint` the way it is defined is actually can be divided into
two parts the `URL` and `Path`, it is delimited by a forward slash.


## ingress

```yaml
ingress:
- <ingressObject>
- <ingressObject>
```

| **Type**                                  | **Required** |
|-------------------------------------------|--------------|
| array of [ingress object](#ingressObject) | no           |


### ingressObject

```yaml
ingress:
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


Each ingress is Kubernetes Ingress spec and name field.
More info: https://kubernetes.io/docs/api-reference/v1.6/#ingressspec-v1beta1-extensions

If there is only one port and user wants to expose the service then user should define one
`ingress` with `host` atleast then the rest of the `ingress` spec(things like `http`, etc.)
will be populated for the user.

More info about Probe: https://kubernetes.io/docs/api-reference/v1.6/#probe-v1-core

#### name

`name: wordpress`

| **Type** | **Required** |
|----------|--------------|
| string   | yes          |

The name of the Ingress.

## Complete example

```yaml
name: database
containers:
- image: mariadb:10
  env:
  - name: MYSQL_ROOT_PASSWORD
    valueFrom:
      secretKeyRef:
        name: wordpress
        key: MYSQL_ROOT_PASSWORD
  - name: MYSQL_DATABASE
    valueFrom:
      configMapKeyRef:
        key: MYSQL_DATABASE
        name: database
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
ingress:
- name: pseudo-wordpress
  rules:
  - host: minikube.local
    http:
      paths:
      - backend:
          serviceName: wordpress
          servicePort: 8080
        path: /
persistentVolumes:
- name: database
  size: 500Mi
  accessModes:
  - ReadWriteOnce
configData:
  MYSQL_DATABASE: wordpress
```