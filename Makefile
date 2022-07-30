DB_URL="postgresql://root:secret@localhost:5432/blogapi?sslmode=disable"

backend:
	go run server/server.go

createBlog:
	go run client/client.go -create

findBlog:
	go run client/client.go -find

postgres:
	docker run --name blogapi -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14-alpine

createdb:
	docker exec -it blogapi createdb --username=root --owner=root blogapi

dropdb:
	docker exec -it blogapi dropdb blog

evans:
	evans --host localhost --port 9090 -r repl


.PHOXY: backend, frontend, postgres, createdb, dropdb, migrateinit, migrateup, migratedown, evans