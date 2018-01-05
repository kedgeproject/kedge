# HealthChecks

This adds support of `livenessProbe` and `readinessProbe`. It is defined at each container level.

Health can be specified similar to the way it is done in Kubernetes's Pod. See the [Kubernetes docs](https://kubernetes.io/docs/api-reference/v1.6/#probe-v1-core) for defining health.

See following code snippet from [guestbook.yaml](./guestbook.yaml):

```yaml
  livenessProbe:
    httpGet:
      path: /
      port: 3000
    initialDelaySeconds: 120
    timeoutSeconds: 5
  readinessProbe:
    httpGet:
      path: /
      port: 3000
    initialDelaySeconds: 5
    timeoutSeconds: 2
```

## Ref:

- [Health API reference](https://kubernetes.io/docs/api-reference/v1.6/#probe-v1)
- [Configure Liveness and Readiness Probes](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/)
- [Container probes](https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/#container-probes)
