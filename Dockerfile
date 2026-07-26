# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o weather-bot ./cmd/bot

# Runtime stage
FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/weather-bot .

ENV DATABASE_PATH=/data/weather-bot.db

VOLUME ["/data"]

CMD ["./weather-bot"]
