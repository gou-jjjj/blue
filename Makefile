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


bin:
	@rm -f ./bin/*
	@go build -o bin/ ./cmd/...
	@cp -r ./config.json bin/
	@mkdir -p bin/logs

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
