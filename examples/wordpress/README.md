# Wordpress

Deploy Wordpress with MariaDB using Kedge.


`mariadb.yaml` defines two Secrets (`database-root-password` and `database-user-password`) and one ConfigMap (`mariadb`).
`wordpress.yaml` gets information on how to connect to database from `database-user-password` Secret and `mariadb` ConfigMap.