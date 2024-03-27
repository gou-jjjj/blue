# makefile	文件每行都是独立的shell

cmdJson="{ \
\"name\": \"\", \
\"summary\": \"\", \
\"group\": \"\", \
\"key\": \"\", \
\"arity\": , \
\"value\": \"\", \
\"arguments\": [ \
\
] \
}"


.server:
	@rm -rf ./bin
	@mkdir -p ./bin/windows/logs ./bin/linux/logs ./bin/darwin/logs
	@GOOS=windows GOARCH=amd64 go build -o ./bin/windows/blue-server.exe ./cmd/...
	@GOOS=linux GOARCH=amd64 go build -o ./bin/linux/blue-server ./cmd/...
	@GOOS=darwin GOARCH=amd64 go build -o ./bin/darwin/blue-server ./cmd/...
	@cp -r ./blue-server.json bin/windows/blue-server.json
	@cp -r ./blue-server.json bin/linux/blue-server.json
	@cp -r ./blue-server.json bin/darwin/blue-server.json
	@echo "Build server success"

bin: .server
	@cp -r ./blue-client/blue-cli.conf  bin/windows/blue-cli.conf
	@cp -r ./blue-client/blue-cli.conf  bin/linux/blue-cli.conf
	@cp -r ./blue-client/blue-cli.conf  bin/darwin/blue-cli.conf
	@GOOS=windows GOARCH=amd64 go build -o ./bin/windows/blue-client.exe ./blue-client/.
	@GOOS=linux GOARCH=amd64 go build -o ./bin/linux/blue-client ./blue-client/.
	@GOOS=darwin GOARCH=amd64 go build -o ./bin/darwin/blue-client ./blue-client/.
	@echo 'Build client success'


exec:
	@rm ./blue-client/exec.go
	@go generate ./script/gen_cli.go

commands:
	@go generate ./script/gen_cmds.go

run:
	@go run ./cmd/main.go


req:run
	@go run ./api/go/main.go


mkcmd:
	@echo "Enter command name: "; \
	read name; \
	echo "${cmdJson}" > ./commands/"$$name.json"; \
	echo "Command created successfully"
