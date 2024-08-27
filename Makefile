# Makefile

api:
	cd src/api && swag init -g cmd/main.go && cd cmd && go run .

movie-event-processor:
	cd src/movie-event-consumer/cmd && go run .
	
user-event-processor:
	cd src/user-event-consumer/cmd && go run .

api-unit-test:
	cd src/api && go test -tags='!integration' ./...

api-integration-test:
	cd src/api/integration_tests && go test -tags=integration  ./...