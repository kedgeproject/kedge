## Guestbook Official Example

* This example is taken from [here](https://kubernetes.io/docs/tutorials/stateless-application/guestbook/).

* This will run only on  Kubernetes.

* To run this example,

```
$ kedge apply -f frontend.yaml -f redis-master.yaml -f redis-slave.yaml
service "frontend" created
deployment "frontend" created
service "redis-master" created
deployment "redis-master" created
service "redis-slave" created
deployment "redis-slave" created
```
