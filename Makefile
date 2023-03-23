migrateup:
	migrate -path db/migration -database "postgresql://postgres:root@localhost:5432/StudentDB?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://postgres:root@localhost:5432/StudentDB?sslmode=disable" -verbose down
