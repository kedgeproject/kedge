# envFrom

For exporting environment variable inside a container which is defined in
`configMap` referencing it can be a painful thing, so `envFrom` helps you import
all the content inside a particular `configMap` into container as env.

At container level you will need to define a field called `envFrom` as shown from
the snippet [web.yaml](web.yaml):

```yaml
  envFrom:
  - configMapRef:
      name: web
```

Here `envFrom` is a list of various references. In above example we are referring
a `configMap` called `web`. This means that populate all the data defined in
`configMap` called `web`.

This also imposes a restriction that you will need to define a `configMap` with
same name(`web` in our case here) in the same file in the root level field called
`configMaps`. See snippet from [web.yaml](web.yaml) below:

```yaml
configMaps:
- data:
    WORDPRESS_DB_NAME: wordpress
    WORDPRESS_DB_HOST: "database:3306"
```

If the referred `configMap` is not defined then the tool will throw an error.

So when this is defined the resulting output will have those `envs` populated
for you as shown below:

```yaml
$ kedge generate -f web.yaml
---
...
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
        name: web
```

In the above converted output you can see that the envs `WORDPRESS_DB_HOST` and
`WORDPRESS_DB_NAME` are auto-populated from the `configMap`.

You might ask why not use the `envFrom` field that is available in the
`container` from Kubernetes 1.6? The answer is lot of users of `kedge` are
still not on latest Kubernetes so this is a way to use this awesome feature
without having the need to use newer cluster, but with downside that you will
have to define configMap in the same file.
