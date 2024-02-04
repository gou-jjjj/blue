package bsp

import "blue/commands"

type HeaderInter interface {
	Type() Header
	Handle() Header
	Len() uint8

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

type Header uint16

const HandleErr Header = 255

func NewHeader(handle Header, Len uint8) Header {
	if handle >= cmdLen {
		return HandleErr
	}
	return Header(Len)<<8 | handle
}

func (h Header) Type() Header {
	return Header(uint8(h>>8)) | TypeMask
}

func (h Header) Handle() Header {
	return Header(uint8(h))
}

func (h Header) Len() uint8 {
	return uint8(h >> 8)
}

func (h Header) Bytes() []byte {
	return []byte{byte(h >> 8), byte(h)}
}

func (h Header) HandleInfo() commands.Cmd {
	return CommandsMap[h]
}
