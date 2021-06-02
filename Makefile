migrateup:
	migrate -path builds/migration -database "mysql://prasta:gladiators88@tcp(localhost:3306)/auth-basic" -verbose up

migratedown:
	migrate -path builds/migration -database "mysql://prasta:gladiators88@tcp(localhost:3306)/auth-basic" -verbose down

dbup:
	podman-compose -f builds/mysql/docker-compose.yml up -d

dbdown:
	podman-compose -f builds/mysql/docker-compose.yml down