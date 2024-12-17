FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

# Install dependencies
RUN go mod tidy

COPY . .

COPY .env /app/.env

WORKDIR /app/cmd/server
RUN go build -o /todo-app .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /todo-app /usr/local/bin/todo-app

COPY --from=builder /app/.env /app/.env

CMD ["todo-app"]
