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
	docker compose run --rm backend golangci-lint run -v --fix

test: build
	docker compose run --rm backend go test -shuffle=on -count=1 ./...

test-app: 
	pnpm install --prefix frontend \
	&& pnpm --prefix frontend test

clean:
	docker compose down --remove-orphans --volumes

generate: env
	docker compose build backend
	docker compose run --rm backend sh scripts/generate.sh

create-migration: env
	docker compose build backend
	docker compose run --rm backend sh db/scripts/create_migration.sh $(name)

migrate:
	docker compose up -d postgres
	docker compose run --rm backend sh db/scripts/migrate.sh

schema-dump: env
	docker compose build backend
	docker compose run --rm backend sh -c "sh db/scripts/dump.sh > backend/db/schema.sql"