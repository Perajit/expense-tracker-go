API_PATH=cmd/api/main.go
MIGRATE_PATH=cmd/migrate/main.go
SEED_PATH=cmd/seed/main.go

api:
	@go run ${API_PATH}

migrate:
	@go run ${MIGRATE_PATH}

seed-dev:
	@go run ${SEED_PATH} -env=dev

seed-prod:
	@go run ${SEED_PATH} -env=prod

migrate-seed-dev:
	@go run ${MIGRATE_PATH}
	@go run ${SEED_PATH} -env=dev
