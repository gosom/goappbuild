#!make
default: help

# generate help info from comments: thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ## help information about make commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

test: ## runs the unit tests
	go test -v -race -short -timeout 1m ./...

test-full: ## runs the unit tests and the end2end tests (slow the first time since it downloads docker images)
	go test -v -race -timeout 5m ./...

test-cover: ## outputs the unittest coverage statistics
	go test -v -race -timeout 5m ./... -coverprofile coverage.out
	go tool cover -func coverage.out
	rm coverage.out

test-cover-report: ## an html report of the coverage statistics
	go test -v ./... -covermode=count -coverpkg=./... -coverprofile coverage.out
	go tool cover -html coverage.out -o coverage.html
	open coverage.html

run: ## runs the application
	go run cmd/goappbuild/main.go

lint: ## runs the linter
	golangci-lint -v run ./...

gen: ## runs go generate
	go generate ./...

dc-up: ## starts the docker-compose dev environment
	docker-compose -f build/docker-compose.dev.yml up -d

dc-down: ## stops the docker-compose dev environment
	docker-compose -f build/docker-compose.dev.yml down -v

db-enter: ## enters the database container
	docker-compose -f build/docker-compose.dev.yml exec db psql -U postgres postgres

db-migrate-up: ## runs the database migrations
	@docker run \
		--network host \
		-v ${PWD}/postgres/migrations:/migrations \
		migrate/migrate \
		-path=/migrations \
		-database=postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable \
		up 

db-migrate-down: ## down the database migrations by 1
	@docker run \
		--network host \
		-v ${PWD}/postgres/migrations:/migrations \
		migrate/migrate \
		-path=/migrations \
		-database=postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable \
		down 1

db-migrate-create: ## creates a new migration
	@docker run \
		--network host \
		-v ${PWD}/postgres/migrations:/migrations \
		migrate/migrate \
		-path=/migrations \
		create -ext sql -dir /migrations -seq $(name)

godoc: ## runs godoc
	@echo "open http://localhost:6060/pkg/github.com/gosom/goappbuild"
	godoc -http=:6060 

