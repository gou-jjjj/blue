// 定义一个包，名为bsp，用于处理头部信息和指令。
package bsp

// HeaderInter 接口定义了头部信息处理的方法集合。
type HeaderInter interface {
	// Type 方法返回头部的类型。
	Type() Header
	// Handle 方法返回头部的句柄。
	Handle() Header

	// HandleInfo 方法返回与头部关联的命令信息。
	HandleInfo() Cmd
	// Bytes 方法返回头部的字节表示。
	Bytes() []byte
}

// 定义头部类型的常量。
const (
	// TypeMask 用于掩码头部类型。
	TypeMask Header = 0b11100000

	// TypeSystem 表示系统类型的头部。
	TypeSystem Header = iota * (1 << 5)
	// TypeDB 表示数据库类型的头部。
	TypeDB
	// TypeNumber 表示数字类型的头部。
	TypeNumber
	// TypeString 表示字符串类型的头部。
	TypeString
	// TypeList 表示列表类型的头部。
	TypeList
	// TypeSet 表示集合类型的头部。
	TypeSet
	// TypeJson 表示JSON类型的头部。
	TypeJson
)

// Header 定义了头部信息的类型，是一个uint8的别名。
type Header uint8

// HandleErr 定义了一个错误句柄值。
const HandleErr Header = 255

// NewHeader 创建一个新的头部实例。
func NewHeader(handle Header) Header {
	return handle
}

// Type 返回头部的类型，通过与TypeMask进行与操作来获取。
func (h Header) Type() Header {
	return h & TypeMask
}

// Handle 返回头部的句柄。
func (h Header) Handle() Header {
	return h
}

// Bytes 返回头部的字节序列表示。
func (h Header) Bytes() []byte {
	return []byte{byte(h)}
}

// HandleInfo 返回与头部关联的命令信息。
func (h Header) HandleInfo() Cmd {
	return CommandsMap[h]
}
