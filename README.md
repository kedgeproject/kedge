# OpenComposition

Experiment to openly composing apps for Kuberentes.

## Why?

You must be wondering why this project when we have other projects to do the same
the answer to that is below:

#### Compared with OpenCompose

* No new language invented, using Kubernetes structs.
* Saves from inventing leaky abstractions.
* Someone defining apps here would find it easier in the world of Kubernetes.
* OpenCompose adds it's own abstractions which is new thing someone needs to
learn, while this approach keeps user near to Kubernetes.
* Also for every new feature a user needs to use she has to rely on OpenCompose
developers to add it to the tool and language, while this approach gives lot of
functionality out of the box.

#### Compared with Kubernetes artifacts

* No need to define everything that Kubernetes needs, define necessary things
and other things will be assumed as defaults by the tool.
* User can choose to define minimum required essentials or define
each and everything by herself.


## Install

```bash
go install
```

## Usage

```bash
$ opencomposition convert -f examples/wordpress/web.yaml -f examples/wordpress/db.yaml
---
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  creationTimestamp: null
  labels:
    app: web
  name: web
spec:
  replicas: 2
  strategy: {}
...
```


