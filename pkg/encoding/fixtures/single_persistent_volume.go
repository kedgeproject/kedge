package fixtures

var SinglePersistentVolume []byte = []byte(
	`name: test
containers:
 - image: nginx
services:
  - ports:
    - port: 8080
persistentVolumes:
- size: 500Mi
`)
