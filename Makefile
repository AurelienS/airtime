.ONESHELL:

.PHONY: run serve gen watch build trigger-refresh browser-refresh zip ent

watch: templ build trigger-refresh

build:
	@go build -o ./tmp/main ./cmd/server/main.go

seed:
	@go run ./cmd/database/seed.go

templ:
	@templ generate

ent:
	@go run -mod=mod entgo.io/ent/cmd/ent generate --feature sql/upsert ./internal/storage/ent


trigger-refresh:
	@touch browser-refresh-trigger.nothing

browser-refresh:
	@browser-sync start --config bs-config.js

zip:
	@zip -r base_code.zip . -x "*.git*" -x "*.vscode*" -x "logs/*" -x "tmp/*" -x "go.sum" -x "README.md" -x ".air.toml" -x "dev.env" -x "Makefile" -x "sqlc.yaml" -x "tailwind.config.js"
