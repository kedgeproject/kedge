# OpenComposition

Alternate experiment to openly composing apps for Kuberentes.


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


