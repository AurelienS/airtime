.ONESHELL:

.PHONY: watch build trigger-refresh browser-refresh zip ent dev

watch: templ build trigger-refresh

render: templ build-render
	./main

dev:
	./dev.sh

build-render:
	@go build -tags netgo -ldflags '-s -w' -o main cmd/server/main.go

build: templ css
	@go build -o ./tmp/main ./cmd/server/main.go

seed:
	@go run ./cmd/database/seed.go

templ:
	@templ generate

css:
	@npx tailwindcss -i ./web/view/styles.css -o ./web/static/styles.css

ent:
	@go run -mod=mod entgo.io/ent/cmd/ent generate --feature sql/upsert ./internal/storage/ent/schema


trigger-refresh:
	@touch browser-refresh-trigger.nothing

browser-refresh:
	@browser-sync start --config bs-config.js

zip:
	@zip -r base_code.zip . -x "*.git*" -x "*.vscode*" -x "logs/*" -x "tmp/*" -x "go.sum" -x "README.md" -x ".air.toml" -x "dev.env" -x "Makefile" -x "sqlc.yaml" -x "tailwind.config.js"
