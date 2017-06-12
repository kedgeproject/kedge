# Using secrets

Using secret is similar to any other pod. As of now this does not provide a way to create a secret but only to consume it.

Creating secret:

```bash
oc create secret generic wordpress --from-literal='MYSQL_ROOT_PASSWORD=rootpasswd,DB_PASSWD=wordpress'
```

Now consuming it, see the snippet from [db.yaml](db.yaml):

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
