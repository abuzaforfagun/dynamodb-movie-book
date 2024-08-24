# Makefile

api:
	cd src/api && swag init -g cmd/main.go && cd cmd && go run .

eventprocessor:
	cd cmd/event_processor && go run .
