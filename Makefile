build:
	go install github.com/rubenv/sql-migrate/...@latest

migrate:
	sql-migrate up

db:
	docker-compose up -d

run: migrate
	go run ./cmd/api/main.go

test:
	go test  ./...