download:
	@echo Download go.mod dependencies
	@go mod download

tools: download
	@echo Installing tools from tools.go
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %
	@go mod tidy

gen: 
	@echo Go generate
	@go generate ./...

generate: tools gen
