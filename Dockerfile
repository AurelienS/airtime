FROM golang:1.21.5-bullseye as builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -v -o main ./cmd/server/main.go

FROM debian:bullseye-slim
COPY --from=builder /app/main /app/main
COPY web/static /app/static
CMD ["/app/main"]
