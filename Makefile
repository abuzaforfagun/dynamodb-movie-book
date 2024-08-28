# Makefile

api:
	cd src/api && swag init -g cmd/main.go && cd cmd && go run .
actor-api:
	cd src/actor-api && swag init -g cmd/main.go && cd cmd && go run .

user-api:
	cd src/user-api && swag init -g cmd/main.go && cd cmd && go run .

movie-event-processor:
	cd src/movie-event-consumer/cmd && go run .
	
review-event-processor:
	cd src/review-event-consumer/cmd && go run .

actor-event-processor:
	cd src/actor-event-consumer/cmd && go run .

api-unit-test:
	cd src/api && go test -tags='!integration' ./...
user-api-unit-test:
	cd src/user-api && go test -tags='!integration' ./...
actor-api-unit-test:
	cd src/user-api && go test -tags='!integration' ./...

api-integration-test:
	cd src/api/integration_tests && go test -tags=integration  ./...
user-api-integration-test:
	cd src/user-api/integration_tests && go test -tags=integration  ./...
actor-integration-test:
	cd src/actor-api/integration_tests && go test -tags=integration  ./...