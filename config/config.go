package config

// 导入包
import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"

	print2 "blue/common/print"
	"blue/datastruct/list"
	"blue/datastruct/number"
	str "blue/datastruct/string"
)

// 配置相关的全局变量
var (
	blueConf BlueConf

	SvrCfg = blueConf.ServerConfig // 服务器配置

	CluCfg = blueConf.ClusterConfig // 集群配置

	LogCfg = blueConf.LogConfig // 日志配置

	CliCfg = blueConf.ClientConfig // 客户端配置

	StoCfg = blueConf.StorageConfig // 存储配置
)

// serverConfig 定义服务器配置结构体
type serverConfig struct {
	Ip         string   `json:"ip,omitempty"`   // 服务器IP
	Port       int      `json:"port,omitempty"` // 服务器端口
	TimeOut    int      `json:"time_out"`       // 超时时间
	DBSum      int      `json:"db_sum"`         // 数据库总数
	GuestToken []string `json:"guest_token"`    // 客户端令牌列表
	RootToken  []string `json:"root_token"`     // 管理员令牌列表
}

// SvrAddr 返回服务器的地址字符串
func (s *serverConfig) SvrAddr() string {
	return fmt.Sprintf("%s:%d", s.Ip, s.Port)
}

// clusterConfig 定义集群配置结构体
type clusterConfig struct {
	Cluster       string `json:"cluster,omitempty"`         // 集群名称
	ClusterAddr   string `json:"cluster_addr,omitempty"`    // 集群地址
	MyClusterAddr string `json:"my_cluster_addr,omitempty"` // 我的集群地址
	Ip            string `json:"ip,omitempty"`              // 本机IP
	Port          int    `json:"port,omitempty"`            // 本机端口
	TryTimes      int    `json:"try_times,omitempty"`       // 尝试次数
	DialTimeout   int    `json:"dial_timeout,omitempty"`    // 拨号超时时间
	Replicas      int    `json:"replicas,omitempty"`        // 副本数
	Consistent    int    `json:"consistent,omitempty"`      // 一致性
}

// logConfig 定义日志配置结构体
type logConfig struct {
	Output   string `json:"output"`              // 日志输出位置
	LogOut   string `json:"log_out,omitempty"`   // 日志输出文件
	LogLevel string `json:"log_level,omitempty"` // 日志级别
}

// clientConfig 定义客户端配置结构体
type clientConfig struct {
	ClientActive int `json:"client_active,omitempty"` // 活跃客户端数
	ClientLimit  int `json:"client_limit,omitempty"`  // 客户端限制数
}

// storageConfig 定义存储配置结构体
type storageConfig struct {
	StorageSet  []string `json:"storage_set,omitempty"`  // 存储集
	StoragePath string   `json:"storage_path,omitempty"` // 存储路径
}

// BlueConf 定义整体配置结构体
type BlueConf struct {
	ServerConfig  serverConfig  `json:"server_config"`  // 服务器配置
	ClusterConfig clusterConfig `json:"cluster_config"` // 集群配置
	LogConfig     logConfig     `json:"log_config"`     // 日志配置
	ClientConfig  clientConfig  `json:"client_config"`  // 客户端配置
	StorageConfig storageConfig `json:"storage_config"` // 存储配置
}

// Entries 将配置项转换为map形式返回
func (c *BlueConf) Entries() map[string]interface{} {
	v := reflect.ValueOf(c).Elem()
	t := v.Type()
	entries := make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		fieldValue := v.Field(i)
		fieldType := fieldValue.Type()

		if fieldValue.Kind() == reflect.Struct {
			for j := 0; j < fieldType.NumField(); j++ {

				if fieldValue.Field(j).CanInterface() {
					switch fieldValue.Field(j).Type().String() {
					case "string":
						entries[fieldType.Field(j).Name] = str.NewString(fieldValue.Field(j).String())
					case "[]string":
						quickList := list.NewQuickList()
						strs, ok := fieldValue.Field(j).Interface().([]string)
						if !ok {
							panic("type assertion failed")
						}
						for _, s := range strs {
							quickList.Add(s)
						}
						entries[fieldType.Field(j).Name] = quickList
					case "int":
						newNumber, err := number.NewNumber(fieldValue.Field(j).Int())
						if err != nil {
							panic(err)
						}
						entries[fieldType.Field(j).Name] = newNumber
					default:
						panic(fieldValue.Field(j).Type().String())
					}
				}
			}
		}
	}

	return entries
}

// 默认配置
var defaultConfig = BlueConf{
	ServerConfig: serverConfig{
		Ip:      "127.0.0.1", // 默认服务器IP
		Port:    8080,        // 默认服务器端口
		TimeOut: 10,          // 默认超时时间
		DBSum:   8,           // 默认数据库总数
	},
	LogConfig: logConfig{
		LogOut:   "./logfile/log.log", // 默认日志输出文件
		LogLevel: "info",              // 默认日志级别
	},
	ClientConfig: clientConfig{
		ClientActive: 10, // 默认活跃客户端数
		ClientLimit:  10, // 默认客户端限制数
	},
	StorageConfig: storageConfig{
		StoragePath: "./storage/data", // 默认存储路径
	},
}

// Init 从指定路径初始化配置
func Init(path string) map[string]interface{} {
	configFile, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	bytes := make([]byte, 0)
	reader := bufio.NewReader(configFile)
	for {
		readByte, err := reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		bytes = append(bytes, readByte)
	}

	err = json.Unmarshal(bytes, &blueConf)
	if err != nil {
		panic(err)
	}

	// 更新全局配置
	SvrCfg = blueConf.ServerConfig
	CluCfg = blueConf.ClusterConfig
	LogCfg = blueConf.LogConfig
	CliCfg = blueConf.ClientConfig
	StoCfg = blueConf.StorageConfig

	print2.ConfigInitSuccess() // 打印配置初始化成功信息

	return blueConf.Entries() // 返回配置项的map形式
}

// OpenCluster 判断是否开启集群
func OpenCluster() bool {
	if strings.ToLower(blueConf.ClusterConfig.Cluster) != "yes" {
		return false
	}
	return true
}

// OpenStorage 判断指定索引的存储是否开启
func OpenStorage(idx string) bool {
	for i := range StoCfg.StorageSet {
		if StoCfg.StorageSet[i] == idx {
			return true
		}
	}

	return false
}
