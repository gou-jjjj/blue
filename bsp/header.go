package bsp

import "blue/commands"

type HeaderInter interface {
	Type() Header
	Handle() Header

	HandleInfo() commands.Cmd
	Bytes() []byte
}

const (
	TypeMask Header = 0b11100000

	TypeSystem Header = iota * (1 << 5)
	TypeDB
	TypeNumber
	TypeString
	TypeList
	TypeSet
	TypeJson
)

type Header uint8

const HandleErr Header = 255

func NewHeader(handle Header) Header {
	return handle
}

func (h Header) Type() Header {
	return h | TypeMask
}

func (h Header) Handle() Header {
	return h
}

func (h Header) Bytes() []byte {
	return []byte{byte(h)}
}

func (h Header) HandleInfo() commands.Cmd {
	return CommandsMap[h]
}
