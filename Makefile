include .env.development # change .development to .production while migrate to production database

build:
	go build -o .bin/main cmd/main/app.go

prod: build
	./.bin/main production

dev:
	go run cmd/main/app.go development

db_up:
	migrate -path ./schema -database "postgresql://admin:${DB_POSTGRES_PASSWORD}@localhost:5432/exchanger?sslmode=disable" up

db_down:
	migrate -path ./schema -database "postgresql://admin:${DB_POSTGRES_PASSWORD}@localhost:5432/exchanger?sslmode=disable" down