FROM golang:1.23-alpine AS builder

RUN apk add --no-cache gcc musl-dev libgcc libc-dev

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o platform-service ./cmd/api

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/platform-service .
COPY --from=builder /app/.env .
EXPOSE 8080

CMD ["./platform-service"]