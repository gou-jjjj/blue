package main

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
)

var BC BlueConf

type BlueConf struct {
	Ip      string `json:"ip,omitempty"`
	Port    int    `json:"port,omitempty"`
	TimeOut uint16 `json:"time_out"`

	LogOut string `json:"log_out,omitempty"`
}

func init() {
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
}
