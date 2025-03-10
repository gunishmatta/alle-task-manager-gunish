FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache build-base

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o /app/task-manager

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache sqlite-libs && install -d /app/data

COPY --from=builder /app/task-manager .

EXPOSE 8080

CMD ["./task-manager"]
