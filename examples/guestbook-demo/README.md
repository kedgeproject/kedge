# Guestbook Demo

This will deploy the famous "Guestbook" from numerous Kubernetes examples! [source](https://kubernetes.io/docs/tutorials/stateless-application/guestbook/)

In all, we will define:
 - services
 - deployments
 - configMap
 - persistentVolumeClaims
 - secrets

# How to deploy

1. Download the files

```sh
wget https://raw.githubusercontent.com/kedgeproject/kedge/master/examples/guestbook-demo/backend.yaml
wget https://raw.githubusercontent.com/kedgeproject/kedge/master/examples/guestbook-demo/frontend.yaml
wget https://raw.githubusercontent.com/kedgeproject/kedge/master/examples/guestbook-demo/db.yaml
```

2. Deploy using `kedge`

```sh
$ kedge apply -f backend.yaml -f frontend.yaml -f db.yaml
service "guestbook" created
deployment "guestbook" created
service "backend" created
deployment "backend" created
persistentvolumeclaim "mongodb-data" created
service "database" created
secret "mongodb-admin" created
secret "mongodb-user" created
configmap "mongodb-user" created
deployment "database" created
```

3. Access your Guestbook instance

If you are using `minikube` for local Kubernetes deployment, you can access your Guestbook instance using `minikube service`

```sh
$ minikube service guestbook 
Opening kubernetes service default/guestbook in default browser...
```
