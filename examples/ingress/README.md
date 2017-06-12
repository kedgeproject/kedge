# Ingress

See the following snippet from [web.yaml](web.yaml).

```yaml
services:
- name: wordpress
  type: LoadBalancer
  ports:
  - port: 8080
    targetPort: 80
  rules:
  - host: minikube.local
```

Services here is the mix of [`service` spec](https://kubernetes.io/docs/api-reference/v1.6/#servicespec-v1-core) and [`ingress` spec](https://kubernetes.io/docs/api-reference/v1.6/#ingressspec-v1beta1-extensions). So if you want a service to be exposed via an ingress set the service `type` to be `LoadBalancer`. Then also define `rules` section.

In `rules` section is similar to the way you define the rules in actual ingress object. If you have only one port and only one rule defined with host then the rest of the rule will be populated for you. But if you have more than one ports or more than one rule then you will also have to define the `http` section in the `rule`.

So the above section gets converted to following and a service:

```yaml
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  labels:
    app: web
  name: wordpress
spec:
  rules:
  - host: minikube.local
    http:
      paths:
      - backend:
          serviceName: wordpress
          servicePort: 8080
        path: /
```

See that the http section has been populated automatically. Because this app has only one port and only one rule.
