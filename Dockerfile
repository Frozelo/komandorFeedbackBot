FROM golang:1.22.1-alpine AS builder


WORKDIR /app


COPY go.mod go.sum ./


RUN go mod download


COPY . .

RUN go build -o main ./main.go


FROM alpine:latest


WORKDIR /app


COPY --from=builder /app/main /app/main
COPY --from=builder /app/internal/bot/config.yml /app/internal/bot/config.yml


EXPOSE 8080

ENTRYPOINT ["/app/main"]
