#!/bin/sh

go mod download
go run github.com/a-h/templ/cmd/templ generate generate
go build -tags netgo -ldflags '-s -w' -o main cmd/server/main.go