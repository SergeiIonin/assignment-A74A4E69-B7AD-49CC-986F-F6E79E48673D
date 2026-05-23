FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY config ./config
COPY cmd ./cmd
COPY internal ./internal

RUN GOOS=linux go build -ldflags="-w -s" -o main ./cmd/main.go

FROM alpine:3.21

RUN adduser -D appuser

WORKDIR /app

COPY --from=builder /app/main .

RUN chown appuser:appuser /app/main

USER appuser

EXPOSE 8080

CMD ["./main"]
