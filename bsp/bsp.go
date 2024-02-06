package bsp

import (
	"context"
	"fmt"
	"io"
	"sync"
)

const (
	Split byte = 0b0
	Done  byte = 0b1

	Expire byte = 0b110
	From   byte = 0b111
	End    byte = 0b1000
)

func AppendSplit(d []byte) []byte {
	return append(d, Split)
}

func AppendDone(d []byte) []byte {
	return append(d, Done)
}

type BspProtoInter interface {
	Key() string
	KeyBytes() []byte

	Values() [][]byte
	ValueBytes() []byte
	ValueStr() string
}

type BspProto struct {
	Header
	key   string
	value [][]byte
}

var bspPool = sync.Pool{
	New: func() interface{} {
		return &BspProto{}
	},
}

// todo: 添加handle
func (b BspProto) String() string {
	return fmt.Sprintf("%d %s %s", b.Header, b.key, b.value)
}

func BspProtos(ctx context.Context, r io.Reader) (chan *BspProto, chan *ErrResp) {
	protos := make(chan *BspProto)
	errs := make(chan *ErrResp)
	go parseReq(ctx, r, protos, errs)
	return protos, errs
}

func (b BspProto) ValueBytes() []byte {
	return b.value[0]
}

func (b BspProto) ValueStr() string {
	return string(b.value[0])
}

func (b BspProto) Key() string {
	return b.key
}

func (b BspProto) KeyBytes() []byte {
	return []byte(b.key)
}

func (b BspProto) Values() [][]byte {
	return b.value
}
