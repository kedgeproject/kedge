# ConfigMaps

## Defining a configMap

See the following snippet from [web.yaml](./web.yaml)

```yaml
configData:
  WORDPRESS_DB_NAME: wordpress
  WORDPRESS_DB_HOST: "database:3306"
```

Define a root level field called `configData`. It is just a key value pair. If this is define a configMap with the `name` of app is created.

## Automatic population

If just `configData` is defined and in no place in the app it is referred and if there is only one container in the pod, then data in `configData` is populated as env in the pod.

In [web.yaml](./web.yaml) only `configData` is defined but it is not referred anywhere. So the entire configData is populated as env in the container.

The converted output snippet of the `deployment` generated for the `web` app looks as follows:

```yaml
$ opencomposition convert -f web.yaml
---
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  creationTimestamp: null
  labels:
    app: web
  name: web
spec:
  replicas: 2
  strategy: {}
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: web
      name: web
    spec:
      containers:
      - env:
        - name: WORDPRESS_DB_PASSWORD
          value: wordpress
        - name: WORDPRESS_DB_USER
          value: wordpress
        - name: WORDPRESS_DB_HOST
          valueFrom:
            configMapKeyRef:
              key: WORDPRESS_DB_HOST
              name: web
        - name: WORDPRESS_DB_NAME
          valueFrom:
            configMapKeyRef:
              key: WORDPRESS_DB_NAME
              name: web
        image: wordpress:4
        livenessProbe:
...
```

See that `WORDPRESS_DB_HOST` and `WORDPRESS_DB_NAME` are automatically populated, even though it is not defined in `containers.env` section in the [web.yaml](./web.yaml).

## Consuming the configMap

See the following code snippet from [db.yaml](./db.yaml)

```yaml
  - name: MYSQL_DATABASE
    valueFrom:
      configMapKeyRef:
        key: MYSQL_DATABASE
        name: database
```

This is similar to the way configMap is referred in Kubernetes.


## Ref

- [Referring a Config Map](https://kubernetes.io/docs/api-reference/v1.6/#envvarsource-v1-core)
- [ConfigMap APIs](https://kubernetes.io/docs/api-reference/v1.6/#configmap-v1-core)
- [Configure Containers Using a ConfigMap](https://kubernetes.io/docs/tasks/configure-pod-container/configmap/)
- [Use ConfigMap Data in Pods](https://kubernetes.io/docs/tasks/configure-pod-container/configure-pod-configmap/)
