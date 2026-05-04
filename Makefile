export

POSTGRES_USER = login
POSTGRES_PASSWORD = pass
POSTGRES_DB = db-name
DB_BASE_URL = postgres://$(POSTGRES_USER):${POSTGRES_PASSWORD}@localhost:5432
DB_MIGRATE_URL = $(DB_BASE_URL)/$(POSTGRES_DB)?sslmode=disable
MIGRATE_PATH = ./migration/postgres/apple

up:
	docker compose  up --build -d --force-recreate

down:
	docker compose down

run: mod
	go run ./cmd/app

mod:
	go mod tidy

mod-update:
	go get -u all
	go mod tidy

lint:
	golangci-lint run

test:
	go test -v -cover ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

migrate-install:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.18.1

migrate-create:
	migrate create -ext sql -dir "$(MIGRATE_PATH)" migration-name

db-create:
	@echo "🔄 Creating database $(POSTGRES_DB) if not exists..."
	@docker exec -t postgres psql -U $(POSTGRES_USER) -d postgres -tc \
		"SELECT 1 FROM pg_database WHERE datname = '$(POSTGRES_DB)'" | grep -q 1 || \
	docker exec -t postgres psql -U $(POSTGRES_USER) -d postgres -c \
		"CREATE DATABASE \"$(POSTGRES_DB)\";"

migrate-up:
	migrate -database "$(DB_MIGRATE_URL)" -path "$(MIGRATE_PATH)" up

migrate-down:
	migrate -database "$(DB_MIGRATE_URL)" -path "$(MIGRATE_PATH)" down -all

swagger:
	swag init -g cmd/server/main.go -o ./docs --parseInternal