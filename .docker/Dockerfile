# --- BUILD STAGE ---
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o api-service ./cmd/main.go

# --- RUN STAGE ---
FROM alpine:latest AS release

COPY --from=builder /app/api-service /app/api-service/
COPY config.yaml /etc/api-service/

WORKDIR /app/api-service/

EXPOSE 8080

CMD ["./api-service"]