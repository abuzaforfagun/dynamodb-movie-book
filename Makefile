# Makefile

movie-api:
	cd src/movie-api && swag init -g cmd/api/main.go && cd cmd/api && go run .
actor-api:
	cd src/actor-api && swag init -g cmd/api/main.go && cd cmd/api && go run .

actor-grpc:
	cd src/actor-api/cmd/grpc && go run .

user-api:
	cd src/user-api && swag init -g cmd/api/main.go && cd cmd/api && go run .

movie-grpc:
	cd src/movie-api/cmd/grpc && go run .

user-grpc:
	cd src/user-api/cmd/grpc && go run .

movie-event-processor:
	cd src/movie-event-consumer/cmd && go run .
	
review-event-processor:
	cd src/review-event-consumer/cmd && go run .

actor-event-processor:
	cd src/actor-event-consumer/cmd && go run .

movie-api-unit-test:
	cd src/movie-api && go test -tags='!integration' ./...
user-api-unit-test:
	cd src/user-api && go test -tags='!integration' ./...
actor-api-unit-test:
	cd src/actor-api && go test -tags='!integration' ./...

movie-api-integration-test:
	cd src/movie-api/integration_tests && go test -tags=integration  ./...
user-api-integration-test:
	cd src/user-api/integration_tests && go test -tags=integration  ./...
actor-api-integration-test:
	cd src/actor-api/integration_tests && go test -tags=integration  ./...