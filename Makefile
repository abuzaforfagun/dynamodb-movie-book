# Makefile

api:
	cd src/api && swag init -g cmd/main.go && cd cmd && go run .

user-api:
	cd src/user-api && swag init -g cmd/main.go && cd cmd && go run .

movie-event-processor:
	cd src/movie-event-consumer/cmd && go run .
	
review-event-processor:
	cd src/review-event-consumer/cmd && go run .

api-unit-test:
	cd src/api && go test -tags='!integration' ./...

api-integration-test:
	cd src/api/integration_tests && go test -tags=integration  ./...