FROM golang:1.25.1-alpine AS builder

WORKDIR /app
COPY . .

RUN go build -o app ./cmd/main.go

FROM alpine:3.19
WORKDIR /app

COPY --from=builder /app/app .

COPY config ./config

EXPOSE 80
CMD ["./app"]
