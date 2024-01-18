.ONESHELL:

.PHONY: run serve gen watch build trigger-refresh browser-refresh zip ent

# Main commands
watch: templ build trigger-refresh

build:
	@go build -o ./tmp/main ./cmd/server/main.go

# Generate code
templ:
	@templ generate

ent:
	@go generate ./internal/storage/ent

# Browser refresh
trigger-refresh:
	@touch browser-refresh-trigger.nothing

browser-refresh:
	@browser-sync start --config bs-config.js

# minimal zip for IA use
zip:
	@zip -r base_code.zip . -x "*.git*" -x "*.vscode*" -x "logs/*" -x "tmp/*" -x "go.sum" -x "README.md" -x ".air.toml" -x "dev.env" -x "Makefile" -x "sqlc.yaml" -x "tailwind.config.js"
