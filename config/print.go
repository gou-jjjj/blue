package config

import "fmt"

func LogInitSuccess() {
	fmt.Println("log init success ...")
}

func ConfigInitSuccess() {
	fmt.Println("config init success ...")
}

func ClusterInitSuccess() {
	fmt.Println("cluster init success ...")
}

func StorageInitSuccess() {
	fmt.Println("storage init success ...")
}

func ServerInitSuccess() {
	fmt.Println("Server init success ...")
}

func ErrPanic(err error, data ...string) {
	if err != nil {
		if len(data) > 0 {
			panic(fmt.Sprintf("%s: %s", err.Error(), data[0]))
		}
		panic(err)
	}
}
