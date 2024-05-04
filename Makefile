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

# 目录结构
binDir="./bin"
windowsBinDir="${binDir}/windows"
linuxBinDir="${binDir}/linux"
darwinBinDir="${binDir}/darwin"

# 服务器端二进制文件名
serverBinary="blue-server"
serverJSON="blue-server.json"

# 客户端二进制文件名
clientBinary="blue-client"
clientConf="blue-cli.conf"

.server:
	@rm -rf ${binDir}
	@mkdir -p ${windowsBinDir}/logs ${linuxBinDir}/logs ${darwinBinDir}/logs
	@GOOS=windows GOARCH=amd64 go build -o ${windowsBinDir}/${serverBinary}.exe ./cmd/...
	@GOOS=linux GOARCH=amd64 go build -o ${linuxBinDir}/${serverBinary} ./cmd/...
	@GOOS=darwin GOARCH=amd64 go build -o ${darwinBinDir}/${serverBinary} ./cmd/...
	@cp -r ./${serverJSON} ${windowsBinDir}/${serverJSON}
	@cp -r ./${serverJSON} ${linuxBinDir}/${serverJSON}
	@cp -r ./${serverJSON} ${darwinBinDir}/${serverJSON}
	@echo "Build server success"

bin: .server
	@cp -r ./blue-client/${clientConf}  ${windowsBinDir}/${clientConf}
	@cp -r ./blue-client/${clientConf}  ${linuxBinDir}/${clientConf}
	@cp -r ./blue-client/${clientConf}  ${darwinBinDir}/${clientConf}
	@GOOS=windows GOARCH=amd64 go build -o ${windowsBinDir}/${clientBinary}.exe ./blue-client/.
	@GOOS=linux GOARCH=amd64 go build -o ${linuxBinDir}/${clientBinary} ./blue-client/.
	@GOOS=darwin GOARCH=amd64 go build -o ${darwinBinDir}/${clientBinary} ./blue-client/.
	@echo 'Build client success'

exec:
	@rm ./blue-client/exec.go
	@go generate ./script/gen_cli.go

commands:
	@go generate ./script/gen_cmds.go

run:
	@go run ./cmd/*

req:run
	@go run ./api/go/main.go

mkcmd:
	@echo "Enter command name: "; \
	read name; \
	echo "${cmdJson}" > ./commands/"$name.json"; \
	echo "Command created successfully"

scp:bin
	scp  ${linuxBinDir}/${serverBinary} $(testAddr)

testClu1:
	@go run ./cmd/* -c ./benchmark/testclu/blue-server1.json

testClu2:
	@go run ./cmd/* -c ./benchmark/testclu/blue-server2.json

t:
	@rm -rf ./storage
	@go run ./cmd/*