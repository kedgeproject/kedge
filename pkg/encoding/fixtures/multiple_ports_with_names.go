package fixtures

var MultiplePortsWithNames []byte = []byte(
	`name: test
containers:
 - image: nginx
services:
- ports:
  - port: 8080
    name: port-1
  - port: 8081
    name: port-2
  - port: 8082
    name: port-3
`)
