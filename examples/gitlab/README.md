Gitlab Example
--------------

### [Example Reference](https://github.com/kubernetes/charts/tree/master/stable/gitlab-ce)

### Generating artifacts

```
$ kedge generate -f gitlab.yml -f redis.yml -f postgres.yml
---
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  creationTimestamp: null
  labels:
    app: gitlab
  name: gitlab
spec:
  strategy: {}
...
...
```

### Deploying on Kubernetes
```
$ kedge apply -f gitlab.yml -f redis.yml -f postgres.yml 
persistentvolumeclaim "gitlab-data" created
persistentvolumeclaim "gitlab-etc" created
service "gitlab" created
secret "gitlab" created
configmap "gitlab" created
deployment "gitlab" created
persistentvolumeclaim "redis-data" created
service "redis" created
secret "redis" created
deployment "redis" created
persistentvolumeclaim "data" created
service "postgresql" created
secret "postgresql" created
deployment "postgresql" created
```

```
$ kubectl get deployments
NAME         DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
gitlab       1         1         1            1           18m
postgresql   1         1         1            1           18m
redis        1         1         1            1           18m
```

```
$ kubectl get services
NAME         CLUSTER-IP   EXTERNAL-IP   PORT(S)                       AGE
gitlab       10.0.0.153   <nodes>       80:32285/TCP,1022:31717/TCP   18m
kubernetes   10.0.0.1     <none>        443/TCP                       6h
postgresql   10.0.0.203   <none>        5432/TCP                      18m
redis        10.0.0.124   <none>        6379/TCP                      18m
```

Once it's exposed to external IP, visit the IP at `http://<EXTERNAL-IP:<PORT>`, you should see a webpage with gitlab login page. 
