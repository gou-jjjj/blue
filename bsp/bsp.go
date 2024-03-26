package bsp

import (
	"context"
	"fmt"
	"io"
	"sync"
)

var bspPool = sync.Pool{
	New: func() interface{} {
		return &BspProto{}
	},
}

func NewBspProto() *BspProto {
	return bspPool.Get().(*BspProto)
}

func PutBspProto(b *BspProto) {
	b.Header = 0
	b.key = ""
	b.value = nil
	b.buf = nil
	bspPool.Put(b)
}

type BspProtoInter interface {
	SetHeader(Header)
	Key() string
	KeyBytes() []byte
	SetKey(string)

	Values() [][]byte
	ValueBytes() []byte
	ValueStr() string
	SetValue([]byte)
	SetValues([][]byte)
	Buf() []byte
	SetBuf([]byte)
}

type BspProto struct {
	Header
	key   string
	value [][]byte
	buf   []byte
}

func (b *BspProto) String() string {
	return fmt.Sprintf("Header: %+v, Key: %s, Value: %s", b.Header.HandleInfo(), b.key, b.value)
}

func BspProtos(ctx context.Context, r io.Reader) (chan *BspProto, chan *ErrResp) {
	protos := make(chan *BspProto)
	errs := make(chan *ErrResp)
	go parseReq(ctx, r, protos, errs)
	return protos, errs
}

func (b *BspProto) ValueBytes() []byte {
	return b.value[0]
}

func (b *BspProto) ValueStr() string {
	return string(b.value[0])
}

func (b *BspProto) SetHeader(h Header) {
	b.Header = h
}

func (b *BspProto) Key() string {
	return b.key
}

func (b *BspProto) SetKey(key string) {
	b.key = key
}

func (b *BspProto) KeyBytes() []byte {
	return []byte(b.key)
}

func (b *BspProto) Values() [][]byte {
	return b.value
}

func (b *BspProto) SetValue(value []byte) {
	b.value = [][]byte{value}
}

func (b *BspProto) SetValues(values [][]byte) {
	b.value = values
}

func (b *BspProto) Buf() []byte {
	return b.buf
}

func (b *BspProto) SetBuf(buf []byte) {
	b.buf = make([]byte, len(buf))
	copy(b.buf, buf)
}
