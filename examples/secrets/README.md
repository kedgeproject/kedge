# Using secrets

Creating secret:

Create a secret by defining it at the root level -
```yaml
secrets:
- name: wordpress
  data:
    MYSQL_ROOT_PASSWORD: YWRtaW4=
    MYSQL_PASSWORD: cGFzc3dvcmQ=
```

Now consuming it, see the snippet from [db.yaml](db.yaml):

```yaml
  envFrom:
  - secretRef:
      name: wordpress
```

Alternatively, it can also be consumed by referencing it manually in `env:`

```yaml
  env:
  - name: MYSQL_ROOT_PASSWORD
    valueFrom:
      secretKeyRef:
        name: wordpress
        key: MYSQL_ROOT_PASSWORD
...
  - name: MYSQL_PASSWORD
    valueFrom:
      secretKeyRef:
        name: wordpress
        key: DB_PASSWD
```

## Ref:

- [Secrets](https://kubernetes.io/docs/concepts/configuration/secret/)
- [Using secret](https://kubernetes.io/docs/api-reference/v1.6/#envvarsource-v1-core)
