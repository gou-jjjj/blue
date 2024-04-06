package bsp

import (
	"context"
	"fmt"
	"io"
	"sync"
)

// bspPool 用于 BSPProto 对象的复用，减少对象分配的开销。
var bspPool = sync.Pool{
	New: func() interface{} {
		return &BspProto{}
	},
}

// NewBspProto 从池中获取一个新的 BspProto 实例，用于减少内存分配。
func NewBspProto() *BspProto {
	return bspPool.Get().(*BspProto)
}

// PutBspProto 将使用完毕的 BspProto 实例放回池中，以便复用。
func PutBspProto(b *BspProto) {
	b.Header = 0
	b.key = ""
	b.value = nil
	b.buf = nil
	bspPool.Put(b)
}

// BspProtoInter 定义了 BspProto 的接口，提供操作 BSP 协议数据的方法。
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

// BspProto 实现了 BspProtoInter 接口，用于处理 BSP 协议的数据结构。
type BspProto struct {
	Header
	key   string
	value [][]byte
	buf   []byte
}

// String 将 BspProto 实例格式化为字符串，方便调试和日志记录。
func (b *BspProto) String() string {
	return fmt.Sprintf("Header: %v, Key: %s, Value: %s", b.Header.HandleInfo().Name, b.Key(), b.Values())
}

// BspProtos 从输入流中解析 BSP 协议数据，生成 BspProto 实例的通道。
func BspProtos(ctx context.Context, r io.Reader) (chan *BspProto, chan *ErrResp) {
	protos := make(chan *BspProto)
	errs := make(chan *ErrResp)
	go parseReq(ctx, r, protos, errs)
	return protos, errs
}

// ValueBytes 返回 BspProto 中的第一个值的字节序列。
func (b *BspProto) ValueBytes() []byte {
	return b.value[0]
}

// ValueStr 返回 BspProto 中的第一个值转换为字符串。
func (b *BspProto) ValueStr() string {
	return string(b.value[0])
}

// SetHeader 设置 BspProto 的头部信息。
func (b *BspProto) SetHeader(h Header) {
	b.Header = h
}

// Key 返回 BspProto 的键名。
func (b *BspProto) Key() string {
	return b.key
}

// SetKey 设置 BspProto 的键名。
func (b *BspProto) SetKey(key string) {
	b.key = key
}

// KeyBytes 返回 BspProto 键名的字节序列。
func (b *BspProto) KeyBytes() []byte {
	return []byte(b.key)
}

// Values 返回 BspProto 中的所有值。
func (b *BspProto) Values() [][]byte {
	return b.value
}

// SetValue 设置 BspProto 的单个值。
func (b *BspProto) SetValue(value []byte) {
	b.value = [][]byte{value}
}

// SetValues 设置 BspProto 的多个值。
func (b *BspProto) SetValues(values [][]byte) {
	b.value = values
}

// Buf 返回 BspProto 的缓冲区数据。
func (b *BspProto) Buf() []byte {
	return b.buf
}

// SetBuf 设置 BspProto 的缓冲区数据。
func (b *BspProto) SetBuf(buf []byte) {
	b.buf = make([]byte, len(buf))
	copy(b.buf, buf)
}
