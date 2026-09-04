.PHONY: build run dev test clean docker-up docker-down migrate

build:
	go build -o bin/server cmd/api/main.go

run: build
	./bin/server

dev:
	go run cmd/api/main.go

test:
	go test ./... -v

clean:
	rm -rf bin/

docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-reset:
	docker compose down -v && docker compose up -d

migrate-up:
	psql -h localhost -U postgres -d leaderboard_system -f migrations/001_init_schema.up.sql

migrate-down:
	psql -h localhost -U postgres -d leaderboard_system -f migrations/001_init_schema.down.sql

tidy:
	go mod tidy

lint:
	golangci-lint run ./...
