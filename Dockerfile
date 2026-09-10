FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o posts-service ./cmd/main

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/posts-service .
COPY config/config.yaml ./config/config.yaml

EXPOSE 8000

CMD ["./posts-service"]