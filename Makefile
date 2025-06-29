up:
	docker compose up -d
dev: up bash

bash: 
	docker compose exec app ash

test: 
	docker compose exec -T app go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out -o coverage.html

t: test

build: 
	docker compose build

logs: 
	docker compose logs -f app

stop: 
	docker compose stop

down: 
	docker compose down