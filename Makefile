# Load .env file
include .env
export $(shell sed 's/=.*//' .env)

# Command migration
MIGRATE_PATH = db/migrations

migrate-up:
	migrate -path $(MIGRATE_PATH) -database "postgres://$(PG_USER):$(PG_PASSWORD)@$(PG_HOST):$(PG_PORT)/$(PG_DB)?sslmode=$(PG_SSLMODE)" up

migrate-down:
	migrate -path $(MIGRATE_PATH) -database "postgres://$(PG_USER):$(PG_PASSWORD)@$(PG_HOST):$(PG_PORT)/$(PG_DB)?sslmode=$(PG_SSLMODE)" down

migrate-force:
	migrate -path $(MIGRATE_PATH) -database "postgres://$(PG_USER):$(PG_PASSWORD)@$(PG_HOST):$(PG_PORT)/$(PG_DB)?sslmode=$(PG_SSLMODE)" force 1
