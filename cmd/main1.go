package main

import (
	"fmt"
)

// 主题接口
type Subject interface {
	Register(observer Observer)
	Unregister(observer Observer)
	Notify(message string)
}

// 观察者接口
type Observer interface {
	Update(message string)
}

// 具体主题
type ConcreteSubject struct {
	observers []Observer
}

func (cs *ConcreteSubject) Register(observer Observer) {
	cs.observers = append(cs.observers, observer)
}

func (cs *ConcreteSubject) Unregister(observer Observer) {
	for i, obs := range cs.observers {
		if obs == observer {
			cs.observers = append(cs.observers[:i], cs.observers[i+1:]...)
			break
		}
	}
}

func (cs *ConcreteSubject) Notify(message string) {
	for _, observer := range cs.observers {
		observer.Update(message)
	}
}

// 具体观察者
type ConcreteObserver struct {
	name string
}

func (co *ConcreteObserver) Update(message string) {
	fmt.Printf("%s 收到通知：%s\n", co.name, message)
}

func main() {
	// 创建主题
	subject := &ConcreteSubject{}

	// 创建观察者
	observer1 := &ConcreteObserver{name: "观察者1"}
	observer2 := &ConcreteObserver{name: "观察者2"}

	// 注册观察者
	subject.Register(observer1)
	subject.Register(observer2)

	// 发送通知
	subject.Notify("Hello, observers!")
	subject.Notify("Hello, observers,dier!")

	// 注销观察者
	subject.Unregister(observer1)

	// 再次发送通知
	subject.Notify("Hello again, observersas!")
}
