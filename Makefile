.PHONY: mg-init db-init db-drop mg-up mg-down model  srv-test srv-up srv-down 

db-init:
	docker exec -it postgres createdb --username=root --owner=root simple_bank
db-drop:
	docker exec -it postgres dropdb   simple_bank
mg-init:
	migrate create -ext sql -dir db/migration -seq init_schema
mg-up:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up
mg-down:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down
model:
	sqlc generate
srv-test:
	go test -v -cover ./...
srv-up:
	docker-compose up
srv-down:
	docker-compose down

