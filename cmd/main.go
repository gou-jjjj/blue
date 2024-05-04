package main

import (
	"flag"
	"time"

	"blue/common/filename"
	print2 "blue/common/print"
	"blue/config"
	"blue/internal"
	"blue/log"
)

// 初始化程序，设置标题
func init() {
	print2.PrintTitle()
}

// 主函数，程序的入口点
func main() {
	flag.Parse() // 解析命令行参数

	// 初始化配置
	configDB := config.Init(*confPath)

	// 初始化日志系统
	log.Init(config.LogCfg.Output, config.LogCfg.LogLevel, config.LogCfg.LogOut)

	// 创建DB实例数组
	dbs := make([]*internal.DB, config.SvrCfg.DBSum+1)

	// 初始化db0
	dbs[0] = internal.NewDB(func(c *internal.DBConfig) {
		c.DataDictSize = 1024
		c.Index = 0
		c.InitData = configDB
	},
	)

	// 初始化db1到dbN（N为配置中的DB数量）
	for i := 1; i <= config.SvrCfg.DBSum; i++ {
		dbs[i] = internal.NewDB(func(c *internal.DBConfig) {
			c.DataDictSize = 1024
			c.Index = i
			// 设置每个DB的存储路径
			c.StorageOptions.DirPath = filename.StorageName(
				config.StoCfg.StoragePath, i)
		})
	}

	// 创建处理程序
	handler := internal.NewBlueServer(dbs...)

	// 初始化并启动服务器
	internal.NewServer(
		func(c *internal.Config) {
			c.Ip = config.SvrCfg.Ip
			c.Port = config.SvrCfg.Port
			c.ClientLimit = config.CliCfg.ClientLimit
			c.Timeout = time.Duration(config.SvrCfg.TimeOut) * time.Second
			c.HandlerFunc = handler
		},
	).Start()
}
