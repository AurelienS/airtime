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

zip:
	@zip -r base_code.zip . -x "*.git*" -x "*.vscode*" -x "logs/*" -x "tmp/*" -x "go.sum" -x "README.md" -x ".air.toml" -x "dev.env" -x "Makefile" -x "sqlc.yaml" -x "tailwind.config.js"
