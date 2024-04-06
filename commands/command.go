package commands

// Cmd 结构体用于定义一个命令
type Cmd struct {
	Name      string   `json:"name"`      // 命令的名称
	Summary   string   `json:"summary"`   // 命令的简要说明
	Group     string   `json:"group"`     // 命令所属的组
	Arity     int      `json:"arity"`     // 命令的参数个数
	Key       string   `json:"key"`       // 命令的关键字
	Value     string   `json:"value"`     // 命令的值
	Arguments []string `json:"arguments"` // 命令的参数列表
}
