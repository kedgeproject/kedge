package fixtures

var MultiplePortsNoNames []byte = []byte(
	`name: test
containers:
 - image: nginx
services:
- name: nginx
  ports:
  - port: 8080
  - port: 8081
  - port: 8082
`)
