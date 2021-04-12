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

.PHONY: docker
docker: build up

.PHONY: build-deps
build-deps:
	@echo Building dependencies image
	@docker build -t deps -f ./deps.dockerfile .

.PHONY: build
build: build-deps
	@echo Building images
	@docker-compose build

.PHONY: up
up:
	@echo Starting up
	@docker-compose up

.PHONY: down
down:
	@echo Shutting down
	@docker-compose down
