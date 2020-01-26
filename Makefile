start-local:
	CONFIG_PATH=local.env go run cmd/app/main.go
init-db:
	docker-compose -f scripts/env/docker-compose.yaml up -d
	ENV_FILE=local.env ./scripts/env/postgres.sh
	ENV_FILE=test.env ./scripts/env/postgres.sh
test:
	go test -v ./...
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out