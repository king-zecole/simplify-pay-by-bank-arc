env:
	@if not exist .env copy .env.example .env

build: env
	docker compose build

server: build
	docker compose up -d

run-all: build \
	migrate \
	server

lint:
	golangci-lint run -v --fix

test: build
	go test -shuffle=on -count=1 ./backend/...

clean:
	docker compose down --remove-orphans --volumes

generate: build
	docker compose run --rm backend sh scripts/generate.sh

create-migration: build
	docker compose run --rm backend sh db/scripts/create_migration.sh $(name)

migrate:
	docker compose up -d postgres
	docker compose run --rm backend sh db/scripts/migrate.sh

schema-dump: build
	docker compose run --rm backend sh -c "sh db/scripts/dump.sh > backend/db/schema.sql"