run:
	go run ./cmd/app/main.go

up:
	docker-compose up -d

build:
	docker-compose build

down:
	docker-compose down
