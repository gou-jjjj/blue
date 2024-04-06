# makefile 文件用于定义自动化构建和部署规则

# 定义一个cmdJson变量，用于存储命令的JSON配置
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

# 定义测试环境的服务器地址
testAddr="root@39.101.195.49:/root/blue/linux"

# .server目标：构建服务器端二进制文件和相关配置
.server:
    # 清理旧的bin目录并创建新结构
	@rm -rf ./bin
	@mkdir -p ./bin/windows/logs ./bin/linux/logs ./bin/darwin/logs
    # 为不同操作系统构建服务器端二进制文件
	@GOOS=windows GOARCH=amd64 go build -o ./bin/windows/blue-server.exe ./cmd/...
	@GOOS=linux GOARCH=amd64 go build -o ./bin/linux/blue-server ./cmd/...
	@GOOS=darwin GOARCH=amd64 go build -o ./bin/darwin/blue-server ./cmd/...
    # 复制配置文件到每个二进制文件目录
	@cp -r ./blue-server.json bin/windows/blue-server.json
	@cp -r ./blue-server.json bin/linux/blue-server.json
	@cp -r ./blue-server.json bin/darwin/blue-server.json
	@echo "Build server success"

# bin目标：构建客户端二进制文件和相关配置
bin: .server
    # 复制客户端配置文件到各个二进制文件目录
	@cp -r ./blue-client/blue-cli.conf  bin/windows/blue-cli.conf
	@cp -r ./blue-client/blue-cli.conf  bin/linux/blue-cli.conf
	@cp -r ./blue-client/blue-cli.conf  bin/darwin/blue-cli.conf
    # 为不同操作系统构建客户端二进制文件
	@GOOS=windows GOARCH=amd64 go build -o ./bin/windows/blue-client.exe ./blue-client/.
	@GOOS=linux GOARCH=amd64 go build -o ./bin/linux/blue-client ./blue-client/.
	@GOOS=darwin GOARCH=amd64 go build -o ./bin/darwin/blue-client ./blue-client/.
	@echo 'Build client success'

# exec目标：执行生成客户端执行文件的脚本
exec:
	@rm ./blue-client/exec.go
	@go generate ./script/gen_cli.go

# commands目标：执行生成命令的脚本
commands:
	@go generate ./script/gen_cmds.go

# run目标：直接运行go代码（用于开发阶段快速测试）
run:
	@go run ./cmd/*


# req目标：运行API服务的go代码
req:run
	@go run ./api/go/main.go

# mkcmd目标：交互式创建新的命令配置文件
mkcmd:
	@echo "Enter command name: "; \
	read name; \
	echo "${cmdJson}" > ./commands/"$$name.json"; \
	echo "Command created successfully"

# scp目标：将构建的linux二进制文件上传到测试服务器
scp:bin
	scp  ./bin/linux/blue-server $(testAddr)

# testClu1和testClu2目标：用于运行集群测试的不同配置
testClu1:
	@go run ./cmd/* -c ./benchmark/testclu/blue-server1.json

testClu2:
	@go run ./cmd/* -c ./benchmark/testclu/blue-server2.json
