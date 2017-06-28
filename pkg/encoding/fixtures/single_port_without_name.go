package fixtures

var SinglePortWithoutName []byte = []byte(
	`name: test
containers:
 - image: nginx
services:
- ports:
  - port: 8080
`)
