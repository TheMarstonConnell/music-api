install: tidy gen
	@go install ./

build: tidy gen
	@go build -o build/backup ./

remove-doc:
	@rm -rf ./doc

run: install remove-doc
	@backup serve

gen:
	@go install github.com/a-h/templ/cmd/templ@latest
	@templ generate

tidy:
	@go mod tidy


format-tools:
	go install mvdan.cc/gofumpt@v0.6.0
	gofumpt -l -w .


lint: format-tools
	golangci-lint run

format: format-tools
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/lcd/statik/statik.go" | xargs gofumpt -w -s
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/lcd/statik/statik.go" | xargs misspell -w
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/lcd/statik/statik.go" | xargs goimports -w -local github.com/jackalLabs/canine-chain



.PHONY: install format lint format-tools gem tidy remove-doc run