.PHONY: default
default: pyra_build pyra

# PYRA

.PHONY: build
build: pyra_build

.PHONY: pyra_build
pyra_build:
	@templ generate
	@go build -o ./bin/pyra ./cmd/pyra

.PHONY: pyra_build_dev
pyra_build_dev:
	@go build -tags dev -o ./tmp/bin/pyra ./cmd/pyra

.PHONY: pyra
pyra: pyra_build
	@./bin/pyra

# END PYRA

# MIGRATE

.PHONY: migrate_build
migrate_build:
	@go build -o ./bin/migrate ./cmd/migrate

.PHONY: migrate
migrate: migrate_build
	@./bin/migrate

.PHONY: rollback
rollback: migrate_build
	@./bin/migrate rollback

.PHONY: migrate-status
migrate-status: migrate_build
	@./bin/migrate status

.PHONY: migrate-version
migrate-version: migrate_build
	@./bin/migrate version

# END MIGRATE

# AIR

.PHONY: air_build
air_build:
	@go build -o ./tmp/bin/air ./deps/air

# END AIR

.PHONY: seed
seed:
	@go run ./database/seeds

.PHONY: clean
clean:
	rm -rf bin/**/*
	rm -rf public/assets/**/*
	rm -rf temp/bin/**/*

.PHONY: dev
dev:
	@air

.PHONY: scratch
scratch:
	@go run ./cmd/scratch

.PHONY: psql
psql:
	@psql-18 -h localhost -p 5432 -d pyra_dev -U pyra

.PHONY: psql-test
psql-test:
	@psql-18 -h localhost -p 5433 -d pyra_test -U pyra
