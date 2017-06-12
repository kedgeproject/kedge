# Simple example

This directory has simplest of the examples.

## Defining Kubernetes Services

See following snippet from [web.yaml](./web.yaml) the way services info can be given in root level `services` field.

```yaml
services:
- name: wordpress
  type: NodePort
  ports:
  - port: 8080
    targetPort: 80
```

It is list of [service spec](https://kubernetes.io/docs/api-reference/v1.6/#servicespec-v1-core), which means that each app can have multiple services defined. Also see that ports are defined in the `services` field and not in the containers. You can choose to declare ports in the `containers.ports` as well but it is not required.

## Automatic Volumes

See the following code snippet from [db.yaml](./db.yaml):

```yaml
  volumeMounts:
  - name: database
    mountPath: /var/lib/mysql
```

To create a Persistent Volume Claim all you need to do is define above snippet in the `container.volumeMounts`.

This will create a PVC of size `100Mi` which is default. If you want to have a PVC of your own size see how to do it under [custom volume](../customVol) example section.
