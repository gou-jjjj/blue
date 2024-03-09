package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"

	"blue/datastruct/list"
	"blue/datastruct/number"
	str "blue/datastruct/string"
	"blue/internal"
)

var BC BlueConf

type serverConfig struct {
	Ip           string   `json:"ip,omitempty"`
	Port         int      `json:"port,omitempty"`
	TimeOut      int      `json:"time_out"`
	DBSum        int      `json:"db_sum"`
	GuestSession []string `json:"guest_session"`
	AdminSession []string `json:"admin_session"`
	RootSession  []string `json:"root_session"`
}

type clusterConfig struct {
	Cluster     int    `json:"cluster,omitempty"`
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
	ClientLive  int `json:"client_live,omitempty"`
	ClientLimit int `json:"client_limit,omitempty"`
}

type storageConfig struct {
	StoragePath string `json:"storage_path,omitempty"`
}

type BlueConf struct {
	ServerConfig serverConfig  `json:"server_config"`
	Cluster      clusterConfig `json:"cluster_config"`
	LogConfig    logConfig     `json:"log_config"`
	ClientConfig clientConfig  `json:"client_config"`
	Storage      storageConfig `json:"storage_config"`
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

func (c *BlueConf) String() string {
	s := fmt.Sprintf(
		`Ip: %v
Port: %v
TimeOut: %v
DBSum: %v
LogOut: %v
LogLevel: %v
ClientLive: %v
ClientLimit: %v
StoragePath: %v
`,
		c.ServerConfig.Ip, c.ServerConfig.Port, c.ServerConfig.TimeOut, c.ServerConfig.DBSum,
		c.LogConfig.LogOut, c.LogConfig.LogLevel,
		c.ClientConfig.ClientLive, c.ClientConfig.ClientLimit,
		c.Storage.StoragePath)
	return s
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
		LogLevel: "Info",
	},
	ClientConfig: clientConfig{
		ClientLive:  10,
		ClientLimit: 10,
	},
	Storage: storageConfig{
		StoragePath: "./storage/data",
	},
}

func InitConfig() *internal.DB {
	configFile, err := os.Open("./config.json")
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

	err = json.Unmarshal(bytes, &BC)
	if err != nil {
		panic(err)
	}

	log.Printf("log init success ...")

	return internal.NewDB(func(c *internal.DBConfig) {
		c.SetStorage = false
		c.DataDictSize = 1024
		c.Index = 0
		c.InitData = BC.Entries()
	})
}
