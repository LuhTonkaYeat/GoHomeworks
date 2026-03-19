FROM golang:1.26.1-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g cmd/gateway/main.go --output ./docs

RUN go build -o gateway ./cmd/gateway

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/gateway .
COPY --from=builder /app/docs ./docs

EXPOSE 8080

CMD ["./gateway"]