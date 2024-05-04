//go:generate go run ./gen_cmds.go
package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"strings"
)

var cmd = "type Cmd struct {\n\tName      string   `json:\"name\"`      // 命令的名称\n\tSummary   string   `json:\"summary\"`   // 命令的简要说明\n\tGroup     string   `json:\"group\"`     // 命令所属的组\n\tArity     int      `json:\"arity\"`     // 命令的参数个数\n\tKey       string   `json:\"key\"`       // 命令的关键字\n\tValue     string   `json:\"value\"`     // 命令的值\n\tArguments []string `json:\"arguments\"` // 命令的参数列表\n}\n\n"

// typeInfo 定义了一系列类型常量。
var typeInfo = `

const (
	TypeSystem Header = iota * (1 << 5)
	TypeDB
	TypeNumber
	TypeString
	TypeList
	TypeSet
	TypeJson
)

`

type Cmd struct {
	Name      string   `json:"name"`      // 命令的名称
	Summary   string   `json:"summary"`   // 命令的简要说明
	Group     string   `json:"group"`     // 命令所属的组
	Arity     int      `json:"arity"`     // 命令的参数个数
	Key       string   `json:"key"`       // 命令的关键字
	Value     string   `json:"value"`     // 命令的值
	Arguments []string `json:"arguments"` // 命令的参数列表
}

// typeMap 是一个字符串映射，用于将类型名称映射到其对应的类型常量。
var typeMap = map[string]string{
	"system": "TypeSystem",
	"db":     "TypeDB",
	"number": "TypeNumber",
	"string": "TypeString",
	"list":   "TypeList",
	"set":    "TypeSet",
	"json":   "TypeJson",
}

// main 函数是程序的入口点。
// 该函数主要完成以下任务：
// 1. 递归遍历指定目录下所有.json文件。
// 2. 读取并解析这些文件，将它们转换为命令结构体。
// 3. 根据这些命令结构体生成相应的Go代码。
func main() {
	// 初始化文件列表
	files := []string{} // 添加更多文件名

	// 遍历指定目录，寻找所有.json文件
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

	// 解析所有找到的.json文件
	var cmds []Cmd
	for _, file := range files {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}

		// 解析JSON数据
		var fileCommands Cmd
		err = json.Unmarshal(data, &fileCommands)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			return
		}

		cmds = append(cmds, fileCommands)
	}

	// 生成Go代码
	var code strings.Builder
	// 添加代码头部注释
	code.WriteString("// Code generated by go generate; DO NOT EDIT.\n")
	code.WriteString("// Code generated by go generate; DO NOT EDIT.\n")
	code.WriteString("// Code generated by go generate; DO NOT EDIT.\n\n")
	// 添加包声明
	code.WriteString("package bsp\n\n")

	code.WriteString(cmd)

	// 添加常量定义
	code.WriteString(fmt.Sprintf("const cmdLen = %d\n\n", len(cmds)))

	// 根据typeMap为每种类型生成相应的常量定义
	for k, v := range typeMap {
		writeType(&code, cmds, k, v)
	}

	// 生成HandleMap数组
	code.WriteString("var HandleMap = [...]string{\n")
	for _, cmd := range cmds {
		constName := strings.ToUpper(strings.ReplaceAll(cmd.Name, " ", "_"))
		code.WriteString(fmt.Sprintf("\t%s: \"%s\",\n", constName, constName))
	}
	code.WriteString("}\n\n")

	// 生成HandleMap2映射表
	code.WriteString("var HandleMap2 = map[string]Header{\n")
	for _, cmd := range cmds {
		constName := strings.ToUpper(strings.ReplaceAll(cmd.Name, " ", "_"))
		code.WriteString(fmt.Sprintf("\t\"%s\": %s,\n", constName, constName))
	}
	code.WriteString("}\n\n")

	// 生成CommandsMap数组
	code.WriteString("var CommandsMap = [...]Cmd{\n")
	for _, cmd := range cmds {
		constName := strings.ToUpper(strings.ReplaceAll(cmd.Name, " ", "_"))
		code.WriteString(
			fmt.Sprintf("\t%s: {Name:\"%s\",Summary: \"%s\", Group: \"%s\", Arity: %d, Key: \"%s\", Value: \"%s\", Arguments: %#v},\n",
				constName, constName, cmd.Summary, cmd.Group, cmd.Arity, cmd.Key, cmd.Value, cmd.Arguments))
	}
	code.WriteString("}\n\n")

	// 将生成的代码写入文件
	err = ioutil.WriteFile("../bsp/commands.go", []byte(code.String()), 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}

	fmt.Println("Commands generated successfully!")
}

// writeType 为指定类型的命令生成常量定义。
// code: 用于写入生成代码的字符串构建器。
// cmds: 所有命令的列表。
// typeName: 需要为其生成常量的命令类型。
// v: 该类型的基常量。
func writeType(code *strings.Builder, cmds []Cmd, typeName, v string) {
	if len(cmds) == 0 {
		return
	}

	// 为类型添加注释和常量定义
	code.WriteString(fmt.Sprintf("// %s -----------------------------\n", typeName))
	code.WriteString("const (\n")
	i := 1
	for _, cmd := range cmds {
		if cmd.Group == typeName {
			constName := strings.ToUpper(strings.ReplaceAll(cmd.Name, " ", "_"))
			code.WriteString(fmt.Sprintf("\t%s Header = %d + %s\n", constName, i, v))
			i++
		}
	}
	code.WriteString(")\n\n")
}
