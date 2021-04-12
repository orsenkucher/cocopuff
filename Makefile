download:
	@echo Download go.mod dependencies
	@go mod download
	@go mod tidy

tools: download
	@echo Installing tools from tools.go
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

gen: 
	@echo Go generate
	@go generate ./...

generate: tools gen
