# Gitea

## How to Deploy

1. Download the files

```sh
$ wget https://raw.githubusercontent.com/kedgeproject/kedge/master/examples/gitea/gitea.yaml
$ wget https://raw.githubusercontent.com/kedgeproject/kedge/master/examples/gitea/postgres.yaml
```

2. Deploy using `kedge`

```sh
$ kedge apply -f gitea.yaml -f postgres.yaml
persistentvolumeclaim "gitea-data" created
persistentvolumeclaim "gitea-config" created
service "gitea" created
secret "gitea" created
configmap "gitea-config" created
deployment "gitea" created
persistentvolumeclaim "postgres-data" created
service "postgresql" created
secret "postgresql" created
deployment "postgresql" created
```

3. Access the service

If you are using `minikube` for local Kubernetes deployment, you can access your service using `minikube service`

```sh
$ minikube service gitea
Opening kubernetes service in default browser...
```

If you are using `minishift` for local OpenShift development, you can create Route  and access your service using it.

```sh
$ oc expose svc gitea
route "gitea" exposed

$ oc get route gitea
...
```
