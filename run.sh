#!/bin/sh

# Generate templates
templ generate

# Build the project for Render with the necessary flags
go build -tags netgo -ldflags '-s -w' -o main cmd/server/main.go

# Run the built binary
./main
