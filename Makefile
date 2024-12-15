include .env

DB_URL := postgresql://$(DB_USER):$(DB_PASSWORD)@localhost:$(DB_PORT)/$(DB_NAME)?sslmode=disable

postgres:
	docker ps -q -f name=$(DB_CONTAINER) || docker run --name $(DB_CONTAINER) -p $(DB_PORT):5432 -e POSTGRES_USER=$(DB_USER) -e POSTGRES_PASSWORD=$(DB_PASSWORD) -d postgres

createdb: 
	docker exec -it $(DB_CONTAINER) createdb --username=$(DB_USER) --owner=$(DB_USER) $(DB_NAME)

dropdb:
	docker exec -it $(DB_CONTAINER) dropdb $(DB_NAME)

migration_up:
	migrate -path db/migrations -database "$(DB_URL)" -verbose up

migration_down:
	migrate -path db/migrations -database "$(DB_URL)" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migration_up migration_down sqlc test
