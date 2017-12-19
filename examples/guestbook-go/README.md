## Guestbook Go Example

* This example is taken from [here](https://github.com/kubernetes/examples/blob/master/guestbook-go/README.md).

* This will run on Kubernetes as well as OpenShift.

* To run this example,

```
$ kedge apply -f guestbook.yml -f redis-master.yml -f redis-slave.yml
service "guestbook" created
deployment "guestbook" created
service "redis-master" created
deployment "redis-master" created
service "redis-slave" created
deployment "redis-slave" created
```
