# Gitlab

[Example Reference](https://github.com/kubernetes/charts/tree/master/stable/gitlab-ce)

## How to Deploy

1. Download the files

```sh
$ wget https://raw.githubusercontent.com/kedgeproject/kedge/master/examples/gitlab/gitlab.yaml
$ wget https://raw.githubusercontent.com/kedgeproject/kedge/master/examples/gitlab/redis.yaml
$ wget https://raw.githubusercontent.com/kedgeproject/kedge/master/examples/gitlab/postgres.yaml
```

2. Deploy using `kedge`

```sh
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

3. Access the service

If you are using `minikube` for local Kubernetes deployment, you can access your service using `minikube service`

```sh
$ minikube service gitlab
Opening kubernetes service in default browser...
```

If you are using `minishift` for local OpenShift development, you can create Route  and access your service using it.

```sh
$ oc expose svc gitlab
route "gitlab" exposed

$ oc get route gitlab
...
```
