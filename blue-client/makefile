build:
	@rm -rf ./bin
	@go build -o ./bin/blue-cli ./
	@cp ./blue-cli.conf ./bin/
	@echo "Build blue-cli successfully"


cli:
	@cd ./bin ; ./blue-cli

run1:
	@go run main.go

.PHONY:
run:build cli