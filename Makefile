.PHONY: run build test clean migrate-up migrate-down docker-build docker-run

run:
	go run cmd/api/main.go

build:
	go build -o bin/api cmd/api/main.go

test:
	go test -v ./...

clean:
	rm -rf bin/

docker-build:
	docker build -t hr-backend:latest .

docker-run:
	docker run -p 8080:8080 --env-file .env hr-backend:latest

deps:
	go mod download
	go mod tidy

lint:
	golangci-lint run

# Database commands
db-create:
	docker run --name hrms-postgres -e POSTGRES_USER=hrms_user -e POSTGRES_PASSWORD=password -e POSTGRES_DB=hrms_db -p 5432:5432 -d postgres:15

db-start:
	docker start hrms-postgres

db-stop:
	docker stop hrms-postgres

redis-create:
	docker run --name hrms-redis -p 6379:6379 -d redis:7-alpine

redis-start:
	docker start hrms-redis

redis-stop:
	docker stop hrms-redis
