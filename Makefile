# Makefile

run:
	swag init -g cmd/main.go && cd cmd && go run .