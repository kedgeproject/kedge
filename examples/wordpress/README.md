# Wordpress

Deploy a Wordpress container with a database (MariaDB) as well as secrets using Kedge!

In all, we will define a:
 - service
 - secret
 - configmap
 - deployment

# How to deploy

1. Download the files

```sh
wget https://raw.githubusercontent.com/kedgeproject/kedge/master/examples/wordpress/wordpress.yaml
wget https://raw.githubusercontent.com/kedgeproject/kedge/master/examples/wordpress/mariadb.yaml
```

2. Deploy using `kedge`

```sh
$ kedge create -f wordpress.yaml -f mariadb.yaml
persistentvolumeclaim "database" created
service "database" created
secret "database-root-password" created
secret "database-user-password" created
configmap "database" created
deployment "database" created
service "wordpress" created
deployment "wordpress" created
```

3. Access your Wordpress instance

If you are using `minikube` for local Kubernetes deployment, you can access your Wordpress instance using `minikube service`

```sh
$ minikube service wordpress
Opening kubernetes service default/wordpress in default browser...
```
