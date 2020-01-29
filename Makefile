start-local:
	CONFIG_PATH=local.env go run cmd/app/main.go
init-db:
	docker-compose -f scripts/env/docker-compose.yaml up -d
	ENV_FILE=local.env ./scripts/env/postgres.sh
	ENV_FILE=test.env ./scripts/env/postgres.sh
	cd postgres/migrations; umigrate migrate -c ../../local.env; umigrate migrate -c ../../test.env;
test:
	ENV_FILE=test.env go test -v ./...
test-unit:
	go test -v -short ./...
test-integration:
	ENV_FILE=test.env go test -run Integration ./...
test-coverage:
	ENV_FILE=test.env go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out