package fixtures

var SingleContainer []byte = []byte(
	`name: test
containers:
 - image: nginx
services:
  - ports:
    - port: 8080
`)
