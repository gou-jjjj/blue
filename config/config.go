package config

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"
)

var BC BlueConf

type serverConfig struct {
	Ip      string `json:"ip,omitempty"`
	Port    int    `json:"port,omitempty"`
	TimeOut int    `json:"time_out"`
	DBSum   int    `json:"db_sum"`
}

type logConfig struct {
	LogOut   string `json:"log_out,omitempty"`
	LogLevel string `json:"log_level,omitempty"`
}

type clientConfig struct {
	ClientLive  int   `json:"client_live,omitempty"`
	ClientLimit int32 `json:"client_limit,omitempty"`
}

type storageConfig struct {
	Path string `json:"path,omitempty"`
}

type BlueConf struct {
	ServerConfig serverConfig  `json:"server_config"`
	LogConfig    logConfig     `json:"log_config"`
	ClientConfig clientConfig  `json:"client_config"`
	Storage      storageConfig `json:"storage_config"`
}

var defaultConfig = BlueConf{
	ServerConfig: serverConfig{
		Ip:      "127.0.0.1",
		Port:    8080,
		TimeOut: 10,
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
		Path: "./storage/data",
	},
}

func InitConfig() {
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
}
