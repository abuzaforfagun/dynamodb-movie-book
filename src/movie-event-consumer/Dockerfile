FROM golang:1.22.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o movie-event-consumer ./cmd

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/movie-event-consumer .

# Expose any ports that might be used (if needed)
# EXPOSE 8080

# Command to run the event consumer
CMD ["./review-event-consumer"]
