package config

import (
	"bufio"
	"encoding/json"
	"io"
	"net"
	"os"
	"reflect"

	"blue/datastruct/list"
	"blue/datastruct/number"
	str "blue/datastruct/string"
)

var (
	blueConf BlueConf

	SvrCfg = blueConf.ServerConfig

	CluCfg = blueConf.ClusterConfig

	LogCfg = blueConf.LogConfig

	CliCfg = blueConf.ClientConfig

	StoCfg = blueConf.StorageConfig
)

type serverConfig struct {
	Ip         string   `json:"ip,omitempty"`
	Port       int      `json:"port,omitempty"`
	TimeOut    int      `json:"time_out"`
	DBSum      int      `json:"db_sum"`
	GuestToken []string `json:"guest_token"`
	RootToken  []string `json:"root_token"`
}

type clusterConfig struct {
	ClusterAddr string `json:"cluster,omitempty"`
	Ip          string `json:"ip,omitempty"`
	Port        int    `json:"port,omitempty"`
	TryTimes    int    `json:"try_times,omitempty"`
	DialTimeout int    `json:"dial_timeout,omitempty"`
	Replicas    int    `json:"replicas,omitempty"`
	Consistent  int    `json:"consistent,omitempty"`
}

type logConfig struct {
	LogOut   string `json:"log_out,omitempty"`
	LogLevel string `json:"log_level,omitempty"`
}

type clientConfig struct {
	ClientActive int `json:"client_active,omitempty"`
	ClientLimit  int `json:"client_limit,omitempty"`
}

type storageConfig struct {
	StorageSet  []string `json:"storage_set,omitempty"`
	StoragePath string   `json:"storage_path,omitempty"`
}

type BlueConf struct {
	ServerConfig  serverConfig  `json:"server_config"`
	ClusterConfig clusterConfig `json:"cluster_config"`
	LogConfig     logConfig     `json:"log_config"`
	ClientConfig  clientConfig  `json:"client_config"`
	StorageConfig storageConfig `json:"storage_config"`
}

func (c *BlueConf) Entries() map[string]interface{} {
	v := reflect.ValueOf(c).Elem() // 确保c是指针类型，并获取所指向的值
	t := v.Type()
	entries := make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		fieldValue := v.Field(i) // 获取字段值
		fieldType := fieldValue.Type()

		if fieldValue.Kind() == reflect.Struct {
			for j := 0; j < fieldType.NumField(); j++ {

				if fieldValue.Field(j).CanInterface() { // 确保字段值可以被接口访问
					switch fieldValue.Field(j).Type().String() {
					case "string":
						entries[fieldType.Field(j).Name] = str.NewString(fieldValue.Field(j).String())
					case "[]string":
						quickList := list.NewQuickList()
						strings, ok := fieldValue.Field(j).Interface().([]string)
						if !ok {
							panic("type assertion failed")
						}
						for _, s := range strings {
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

var defaultConfig = BlueConf{
	ServerConfig: serverConfig{
		Ip:      "127.0.0.1",
		Port:    8080,
		TimeOut: 10,
		DBSum:   8,
	},
	LogConfig: logConfig{
		LogOut:   "./logfile/log.log",
		LogLevel: "info",
	},
	ClientConfig: clientConfig{
		ClientActive: 10,
		ClientLimit:  10,
	},
	StorageConfig: storageConfig{
		StoragePath: "./storage/data",
	},
}

func InitConfig(path string) map[string]interface{} {
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

	SvrCfg = blueConf.ServerConfig
	CluCfg = blueConf.ClusterConfig
	LogCfg = blueConf.LogConfig
	CliCfg = blueConf.ClientConfig
	StoCfg = blueConf.StorageConfig

	ConfigInitSuccess()

	return blueConf.Entries()
}

func (c clusterConfig) OpenCluster() bool {
	if c.ClusterAddr == "" {
		return false
	}

	ip, _, err := net.SplitHostPort(c.ClusterAddr)
	if err != nil {
		return false
	}

	// 解析IP地址
	ipAddress := net.ParseIP(ip)
	if ipAddress == nil {
		return false
	}

	return true
}

func (c storageConfig) OpenStorage(idx string) bool {
	for i := range c.StorageSet {
		if c.StorageSet[i] == idx {
			return true
		}
	}

	return false
}
