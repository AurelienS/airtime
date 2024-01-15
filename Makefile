.ONESHELL:

.PHONY: run serve gen tailwind watch build trigger-refresh browser-refresh zip

run: gen serve

serve:
	@go run cmd/server/main.go

gen:
	@go generate ./..
	@templ generate

watch: tailwind templ build trigger-refresh

templ:
	@templ generate &

build:
	@go build -o ./tmp/main ./cmd/server/main.go

tailwind:
	@npm run build-css &

trigger-refresh:
	@touch browser-refresh-trigger.nothing

browser-refresh:
	@browser-sync start --config bs-config.js

zip:
	@zip -r base_code.zip . -x "*.git*" -x "*.vscode*" -x "logs/*" -x "tmp/*" -x "go.sum" -x "README.md" -x ".air.toml" -x "dev.env" -x "Makefile" -x "sqlc.yaml" -x "tailwind.config.js"
