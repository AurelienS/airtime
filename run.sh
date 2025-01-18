#!/bin/sh

go mod download
templ generate
go build -tags netgo -ldflags '-s -w' -o main cmd/server/main.go