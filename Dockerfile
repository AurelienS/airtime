FROM golang:1.21.5-bullseye as builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -v -o main ./cmd/server/main.go

FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y ca-certificates bash postgresql-client && rm -rf /var/lib/apt/lists/*
COPY --from=builder /app/main /app/main
COPY web/static /app/static

CMD ["/app/main"]
