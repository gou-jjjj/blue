package cluster

import (
	"fmt"
)

// 主题接口
type Subject interface {
	Register(addr ...string)
	Unregister(addr ...string)
	Offline(addr ...string)
	Online(addr ...string)
}

// 观察者接口
type Observer interface {
	Online(string)
	Offline(string)
}

// 具体观察者
type ConcreteObserver struct {
	name string
}

func (co *ConcreteObserver) Update(message string) {
	fmt.Printf("%s 收到通知：%s\n", co.name, message)
}
