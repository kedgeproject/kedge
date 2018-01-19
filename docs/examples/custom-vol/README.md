# Volumes

To define a volume you have to do two things

- Define a `volumeMount` in `containers.volumeMounts`

Check out the following snippet from [redis-master.yaml](./redis-master.yaml)
```yaml
  volumeMounts:
  - name: persistent
    mountPath: /data
```

Here you mention what is the name of the volume from the root level in `name` field and then in `mountPath` define the path where you wanna mount the volume inside the container.

- Secondly define root level `volumeClaims`

Check out the following snippet from [db.yaml](./db.yaml)
```yaml
volumeClaims:
- name: persistent
  size: 500Mi
  accessModes:
  - ReadWriteOnce
```

The `name` here should match with the `name` field in `containers.volumeMounts`. This is where you define the `size` of the volume as well.

Field `accessModes` is optional and defaults to `ReadWriteOnce`.

## Ref:

- [Container level Volume Mounts](https://kubernetes.io/docs/api-reference/v1.6/#volumemount-v1-core)
- [volumeClaims](https://kubernetes.io/docs/api-reference/v1.6/#volume-v1-core)
