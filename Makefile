include .env

dev: 
	@task :start -w --interval=500ms

start: 
	@go run .

build:
	@go build -o bin/ .

db-create: 
	@dbmate -u $(DATABASE_URL) create

db-drop:
	@dbmate -u $(DATABASE_URL) drop

migrate: 
	@dbmate up

setup: 
	@go mod tidy
	@go mod download
	@docker-compose up -d
	@if [ ! -f .env ]; then cp .env.example .env; fi
	@make db-create
	@make migrate
	@make start