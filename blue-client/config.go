package main

import (
	"bufio"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var BC = new(BlueConf)

type BlueConf struct {
	Ip       string `json:"ip,omitempty"`
	Port     int    `json:"port,omitempty"`
	TimeOut  int    `json:"time_out"`
	TryTimes int    `json:"try_times"`

	LogOut string `json:"log_out,omitempty"`
}

var defaultConf = BlueConf{
	Ip:      "127.0.0.1",
	Port:    8080,
	TimeOut: 60,
	LogOut:  "./log",
}

func init() {
	conf, err := os.Open("./blue-cli.conf")
	if err != nil {
		panic(err)
	}
	defer conf.Close()

	red := bufio.NewReader(conf)

	for {
		readString, err := red.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		if readString == "" {
			continue
		}

		parse(readString, BC)

	}
}

func parse(str string, conf *BlueConf) {
	str = strings.TrimSpace(str)

	if strings.HasPrefix(str, "#") || strings.HasPrefix(str, "//") {
		return
	}

	split := strings.Split(str, " ")
	confs := make([]string, 0)

	for i := range split {
		if split[i] == "" {
			continue
		}

		confs = append(confs, split[i])
		if len(confs) == 2 {
			break
		}

	}

	// 通过反射将 confs 赋值给conf
	elem := reflect.ValueOf(conf).Elem()

	if len(confs) == 2 {
		field := elem.FieldByName(confs[0])

		if field.IsValid() && field.CanSet() {
			switch field.Kind() {
			case reflect.String:
				field.SetString(confs[1])
			case reflect.Int:
				parseInt, err := strconv.Atoi(confs[1])
				if err == nil {
					field.SetInt(int64(parseInt))
				}
			case reflect.Uint16:
				parseUint, err := strconv.ParseUint(confs[1], 10, 16)
				if err == nil {
					field.SetUint(parseUint)
				}
			default:
				panic("unknown config")
			}
		}
	}

}
