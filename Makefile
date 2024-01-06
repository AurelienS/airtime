.PHONY: run serve gen

run: gen serve

serve:
	@go run cmd/server/main.go

gen: 
	@sqlc generate
	@templ generate

watch:
	@templ generate
	@go build -o ./tmp/main ./cmd/server/main.go
