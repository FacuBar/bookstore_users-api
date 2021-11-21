mysql:
	docker run --name users-mysql -p 9000:3306 -e MYSQL_ROOT_PASSWORD=secret -d mysql

createdb: 
	docker exec -it users-mysql mysql --user='root' --password='secret' --execute='CREATE DATABASE users_db'

dropdb:
	docker exec -it users-mysql mysql --user='root' --password='secret' --execute='DROP DATABASE users_db'

migrateup:
	migrate -path migration/ -database "mysql://root:secret@tcp(localhost:9000)/users_db" -verbose up

migratedown: 
	migrate -path migration/ -database "mysql://root:secret@tcp(localhost:9000)/users_db" -verbose down

server:
	go run cmd/main.go

.PHONY: mysql createdb dropdb	migrateup migratedown server