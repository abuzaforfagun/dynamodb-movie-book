# Makefile

api:
	cd src/api && swag init -g cmd/main.go && cd cmd && go run .

eventprocessor:
	cd src/event-consumer/cmd && go run .

api-unit-test:
	cd src/api && go test -tags='!integration' ./...

api-integration-test:
	cd src/api/integration_tests && go test -tags=integration  ./...