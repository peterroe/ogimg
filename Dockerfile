ARG CACHEBUST=1

FROM golang:1.23.3 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ogimg ./cmd/server/main.go

FROM alpine:latest

WORKDIR /root/
COPY config/ ./config/

COPY --from=builder /app/ogimg .

EXPOSE 8888

CMD ["./ogimg"]
