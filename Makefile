# Makefile

api:
	swag init -g cmd/api/main.go && cd cmd/api && go run .

eventprocessor:
	cd cmd/event_processor && go run .
