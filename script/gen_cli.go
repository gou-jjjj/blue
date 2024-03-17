//go:generate go run gen_cli.go

package main

import (
	"blue/commands"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"path/filepath"
)

func main() {
	// 读取 JSON 文件
	files := []string{} // 添加更多文件名

	err := filepath.Walk("../commands", func(path string, info fs.FileInfo, err error) error {
		if filepath.Ext(path) != ".json" {
			return nil
		}

		abs, err := filepath.Abs(path)
		if err != nil {
			return nil
		}

		files = append(files, abs)
		return nil
	})
	if err != nil {
		fmt.Println("Error walking directory:", err)
		return
	}

	var cmds []commands.Cmd
	for _, file := range files {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}

		// 解析 JSON 数据
		var fileCommands commands.Cmd
		err = json.Unmarshal(data, &fileCommands)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			return
		}

		cmds = append(cmds, fileCommands)
	}

}
