# Kapp roadmap

The goal of Kapp is to have the __best__ possible experience when deploying Kubernetes artifacts. 

Here we outline our goals for future releases:

## Kapp 0.1.0

* [ ] Parameterize the ingress host value or endpoint value
* [ ] Rename persistentVolumes to something better like volumeClaims
* [ ] Add support for Jobs, CronJobs, HorizontalPodAutoScaler, etc.
* [ ] Comprehensive documentation of each and every feature
* [ ] The way secrets and configmaps are referred in containers.env it is too much indirection, making it easier.
* [ ] configData are confusing terms rename them.
* [ ] Add convenient shortcuts to the make configMapKeyRef and secretKeyRef usage becomes easier.
* [ ] Add such shortcuts to other part of the spec, but user can choose to use shortcuts or define everything without having to use shortcuts.
* [ ] Create intelligent defaults to the tool, to reduce what a user will write.
* [ ] Generate OpenShift artifacts? All things related to OpenShift like builds, routes, etc.
* [ ] Fields could be renamed to have easier or better name.
