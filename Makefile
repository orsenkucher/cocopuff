.DEFAULT_GOAL := tools

.PHONY: download
download:
	@echo Download go.mod dependencies
	@go mod download
	@go mod tidy

.PHONY: tools
tools: download
	@echo Installing tools from tools.go
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

.PHONY: gen
gen: 
	@echo Go generate
	@go generate ./...

.PHONY: generate
generate: tools gen
