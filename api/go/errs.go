package blue

import (
	"fmt"
)

func ErrArgu(s ...string) error {
	return fmt.Errorf("(error)  wrong number of arguments for '%s' command", s[0])
}

func ErrType(s ...string) error {
	return fmt.Errorf("(error)  unknown type '%s'", s[0])
}

func ErrCommandType(s ...string) error {
	return fmt.Errorf("(error)  unknown command type '%s'", s[0])
}

func ErrDataType(s ...string) error {
	return fmt.Errorf("(error)  unknown data type '%s'", s[0])
}

func ErrCommand(s ...string) error {
	return fmt.Errorf("(error)  unknown command '%s'", s[0])
}

func ErrSyntax(s ...string) error {
	return fmt.Errorf("(error)  syntax error '%s'", s[0])
}

func ErrCommandNil(s ...string) error {
	return fmt.Errorf("(error)  command is '%s'", "nil")
}

func ErrConnect(s ...string) error {
	return fmt.Errorf("(error)  connect error '%s'", s[0])
}

func ErrRead(s ...string) error {
	return fmt.Errorf("(error)  read command error '%s'", s[0])
}

func ErrInvalidResp(s ...string) error {
	return fmt.Errorf("(error)  %s", "invalid response")
}

/*
本科毕业设计（论文）
题目名称：基于GO语言实现的分布式缓存系统

使用的技术：GO语言、一致性哈希算法、时间轮算法、Lsm-Tree、TCP协议、分布式系统

功能包括：多数据类型(string、list、set、map), 数据过期, 数据持久化,分片集群, 功能配置，分级日志，容器部署，多语言客户端

根据以上内容，帮我完善下面目录结构

第1章 绪论
1.1 选题背景及研究意义
1.2 国内外研究现状
1.3 系统目标及内容
1.4 本章小结
第2章 缓存系统需求分析
2.1 系统可行性分析
2.1.1 技术可行性
2.1.2 经济可行性
2.2 功能需求分析
2.2.1 业务流程
2.2.2 功能列表及其说明
2.3
2.4性能需求
2.5 本章小结
3. 缓存系统概要设计
3.1 系统结构（软件结构）
3.2
3.3
3.4
3.5
3.6
3.7 本章小结
4. 缓存系统详细设计与实现
按模块，或可视组件来描述
5. 系统评估（case study）
6. 系统测试
6.1 功能测试
6.2 性能测试
结论
致谢
参考文献
*/
