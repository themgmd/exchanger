include .env

build:
	go build -o .bin/main cmd/main/app.go

prod: build
	./.bin/main -config stable

dev:
	go run cmd/main/app.go