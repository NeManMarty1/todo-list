FROM golang:1.23-bullseye AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN go build -o /todo-api ./cmd/app/main.go

FROM debian:bullseye-slim
WORKDIR /root/
COPY --from=builder /todo-api .
COPY migrations /migrations
COPY config.yaml .
COPY .env .
EXPOSE 8080
CMD ["./todo-api"]