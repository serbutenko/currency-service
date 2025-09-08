FROM golang:1.22 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o currency-service ./cmd/server
FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=builder /app/currency-service .
EXPOSE 50051
CMD ["./currency-service"]
