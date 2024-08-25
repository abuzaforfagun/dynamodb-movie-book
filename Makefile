# Makefile

api:
	cd src/api && swag init -g cmd/main.go && cd cmd && go run .

eventprocessor:
	cd src/event-consumer/cmd && go run .
