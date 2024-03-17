build:
	@rm -rf ./bin
	@go build -o ./bin/blue-cli ./
	@cp ./blue-cli.conf ./bin/
	@echo "Build blue-cli successfully"


cli:
	@cd ./bin && ./blue-cli


.PHONY:
run:build cli