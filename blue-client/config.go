package main

var BC BlueConf

type BlueConf struct {
	Ip      string `json:"ip,omitempty"`
	Port    int    `json:"port,omitempty"`
	TimeOut uint16 `json:"time_out"`

	LogOut string `json:"log_out,omitempty"`
}

var defaultConf = BlueConf{
	Ip:      "127.0.0.1",
	Port:    8080,
	TimeOut: 60,
	LogOut:  "./log",
}
