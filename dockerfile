FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/server

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/app .

EXPOSE 8080

ENV REDIS_URL=redis://localhost:6379 \
    REAL_IP=false

ENTRYPOINT ["./app"]