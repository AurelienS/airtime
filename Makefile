.PHONY: run serve gen tailwind watch build

run: gen serve

serve:
	@go run cmd/server/main.go

gen:
	@go generate ./..
	@templ generate

watch: templ build

templ:
	@templ generate 

build:
	@go build -o ./tmp/main ./cmd/server/main.go 

tailwind:
	@npx tailwind -i web/template/style.css -o web/static/style.css --watch

zip:
	@zip -r base_code.zip . -x "*.git*" -x "*.vscode*" -x "logs/*" -x "tmp/*" -x "go.sum" -x "README.md" -x ".air.toml" -x "dev.env" -x "Makefile" -x "sqlc.yaml" -x "tailwind.config.js"
