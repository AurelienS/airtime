#!/bin/sh

go mod download
go install github.com/a-h/templ/cmd/templ@latest
templ generate
go build -tags netgo -ldflags '-s -w' -o main cmd/server/main.go