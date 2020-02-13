start-local:
	CONFIG_PATH=local.env go run cmd/app/main.go
init-db:
	docker-compose -f scripts/env/docker-compose.yaml up -d
	ENV_FILE=local.env ./scripts/env/postgres.sh
	ENV_FILE=test.env ./scripts/env/postgres.sh
	cd postgres/migrations; umigrate migrate -c ../../local.env; umigrate migrate -c ../../test.env;
init-ci-env:
	docker-compose -f scripts/env/docker-compose.yaml up -d
	ENV_FILE=test.env ./scripts/env/postgres.sh
test:
	ENV_FILE=test.env go test -v ./...
test-unit:
	go test -v -short ./...
test-integration:
	ENV_FILE=test.env go test -run Integration ./...
test-coverage:
	ENV_FILE=test.env go test -coverprofile=coverage.out ./...
test-html-coverage:
	ENV_FILE=test.env go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
test-unit-coverage:
	go test -coverprofile=coverage-unit.out -short ./...
	go tool cover -html=coverage-unit.out
test-integration-coverage:
	ENV_FILE=test.env go test -coverprofile=coverage-integration.out -run Integration ./...
	go tool cover -html=coverage-integration.out