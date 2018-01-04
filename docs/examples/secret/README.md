# Using secrets

Creating secret:

Create a secret by defining it at the root level -
```yaml
secrets:
- name: secret
  data:
    GET_HOSTS_FROM: RE5T
```
Make sure everything put in the field `data:` is base64 encoded.
For supplying plaintext secret data, use the field `stringData`.

Now consuming it, see the snippet from [db.yaml](db.yaml):

```yaml
  envFrom:
  - secretRef:
      name: secret
```

Alternatively, it can also be consumed by referencing it manually in `env:`

```yaml
  env:
  - name: GET_HOSTS_FROM
    valueFrom:
      secretKeyRef:
        name: secret
        key: GET_HOSTS_FROM
```

## Ref:

- [Secrets](https://kubernetes.io/docs/concepts/configuration/secret/)
- [Using secret](https://kubernetes.io/docs/api-reference/v1.6/#envvarsource-v1-core)
