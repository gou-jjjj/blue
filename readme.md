### Blue 是什么？

blue 是一个用Golang语言实现的类似redis一样的缓存数据库，支持的特性包括：
* 持久化
* 多数据类型支持
* 切片集群
* 数据过期
* 数据淘汰
* 日志


### Blue 的使用

1. 直接运行Bin目录中对应平台的文件
2. 使用go run main.go运行
3. 使用make bin成后bin目录文件

### docker快速开始
#### 构建 Docker 镜像
```bash
docker build -t blue:0.1 .
```
#### 运行 Docker 容器
```bash
docker run -d \
  --name blue \
  -p 13140:13140 \
  blue:0.1
```
  
  
