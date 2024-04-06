package bsp

import (
	"blue/common/strbytes"
	"fmt"
	"strconv"
)

// Reply 接口定义了回复的基本行为，包括转化为字节切片和字符串。
type Reply interface {
	Bytes() []byte
	String() string
}

// InfoReply 结构体实现了 Reply 接口，用于存储信息回复。
type InfoReply struct {
	i    ReplyType
	info []byte
}

// Bytes 将 InfoReply 转化为字节切片。
func (i InfoReply) Bytes() []byte {
	res := make([]byte, 0, len(i.info)+2)
	res = append(res, byte(i.i))
	res = append(res, i.info...)
	return AppendDone(res)
}

// String 将 InfoReply 转化为字符串。如果 info 为空，则只返回 ReplyType 的字符串表示。
func (i InfoReply) String() string {
	if i.info == nil || len(i.info) == 0 {
		return MessageMap[i.i]
	} else {
		return fmt.Sprintf("%s:%s", MessageMap[i.i], i.info)
	}
}

// NewInfo 创建一个新的 InfoReply 实例。
func NewInfo(i ReplyType, info ...[]byte) *InfoReply {
	if len(info) != 0 {
		return &InfoReply{
			i:    i,
			info: info[0],
		}
	}

	return &InfoReply{
		i: i,
	}
}

// NumResp 结构体实现了 Reply 接口，用于存储数字回复。
type NumResp struct {
	n   ReplyType
	num []byte
}

// Bytes 将 NumResp 转化为字节切片。
func (n NumResp) Bytes() []byte {
	res := make([]byte, 0, len(n.num)+2)
	res = append(res, byte(n.n))
	res = append(res, n.num...)
	return AppendDone(res)
}

// String 将 NumResp 转化为字符串。
func (n NumResp) String() string {
	return strbytes.Bytes2Str(n.num)
}

// NewNum 创建一个新的 NumResp 实例。
func NewNum(num any) *NumResp {
	switch num.(type) {
	case uint8:
		return &NumResp{
			n:   ReplyNumber,
			num: []byte{num.(uint8)},
		}
	case uint64:
		return &NumResp{
			n:   ReplyNumber,
			num: strbytes.Uint642Bytes(num.(uint64)),
		}
	case int64:
		return &NumResp{
			n:   ReplyNumber,
			num: strbytes.Uint642Bytes(uint64(num.(int64))),
		}
	case int:
		return &NumResp{
			n:   ReplyNumber,
			num: strbytes.Uint642Bytes(uint64(num.(int))),
		}
	case []byte:
		return &NumResp{
			n:   ReplyNumber,
			num: num.([]byte),
		}
	case string:
		return &NumResp{
			n:   ReplyNumber,
			num: []byte(num.(string)),
		}
	default:
		panic("wrong type")
		return &NumResp{}
	}
}

// StrResp 结构体实现了 Reply 接口，用于存储字符串回复。
type StrResp struct {
	s   ReplyType
	msg []byte
}

// Bytes 将 StrResp 转化为字节切片。
func (s StrResp) Bytes() []byte {
	res := make([]byte, 0, len(s.msg)+2)
	res = append(res, byte(s.s))
	res = append(res, s.msg...)
	return AppendDone(res)
}

// String 将 StrResp 转化为字符串。
func (s StrResp) String() string {
	return string(s.msg)
}

// NewStr 创建一个新的 StrResp 实例。
func NewStr(msg any) *StrResp {
	s := &StrResp{
		s: ReplyString,
	}
	switch msg.(type) {
	case string:
		s.msg = []byte(msg.(string))
	case []byte:
		s.msg = msg.([]byte)
	case uint8:
		s.msg = []byte(strconv.Itoa(int(msg.(uint8))))
	default:
		panic("wrong type")
	}

	return s
}

// ListResp 结构体实现了 Reply 接口，用于存储列表回复。
type ListResp struct {
	l    ReplyType
	list [][]byte
}

// Bytes 将 ListResp 转化为字节切片。
func (l ListResp) Bytes() []byte {
	b := make([]byte, 0, len(l.list))
	b = append(b, byte(l.l))
	for i := range l.list {
		b = append(b, AppendSplit(l.list[i])...)
	}

	return AppendDone(b)
}

// String 将 ListResp 转化为字符串。
func (l ListResp) String() string {
	return fmt.Sprintf("%v", l.list)
}

// NewList 创建一个新的 ListResp 实例。
func NewList(list ...[]byte) *ListResp {
	return &ListResp{
		l:    ReplyList,
		list: list,
	}
}

// ErrResp 结构体实现了 Reply 接口，用于存储错误回复。
type ErrResp struct {
	e   ReplyType
	msg []byte
}

// Error 实现了 error 接口，返回错误信息的字符串表示。
func (r ErrResp) Error() string {
	return r.String()
}

// Bytes 将 ErrResp 转化为字节切片。
func (r ErrResp) Bytes() []byte {
	res := make([]byte, 0, len(r.msg)+2)
	res = append(res, byte(r.e))
	res = append(res, r.msg...)
	return AppendDone(res)
}

// String 将 ErrResp 转化为字符串。如果 msg 为空，则只返回 ReplyType 的字符串表示。
func (r ErrResp) String() string {
	if r.msg == nil || len(r.msg) == 0 {
		return MessageMap[r.e]
	}
	return fmt.Sprintf("%s:%s", MessageMap[r.e], r.msg)
}

// NewErr 创建一个新的 ErrResp 实例。
func NewErr(e ReplyType, msg ...any) *ErrResp {
	re := &ErrResp{
		e: e,
	}

	if len(msg) != 0 {
		switch msg[0].(type) {
		case []byte:
			re.msg = msg[0].([]byte)
		case string:
			re.msg = []byte(msg[0].(string))
		}
	}

	return re
}

// ClusterReply 结构体存储集群回复的缓冲区。
type ClusterReply struct {
	buf []byte
}

// NewClusterReply 创建一个新的 ClusterReply 实例。
func NewClusterReply(buf []byte) *ClusterReply {
	return &ClusterReply{
		buf: buf,
	}
}

// Bytes 返回 ClusterReply 的字节切片。
func (c ClusterReply) Bytes() []byte {
	return c.buf
}

// String 返回 ClusterReply 的字符串表示。
func (c ClusterReply) String() string {
	return string(c.buf)
}

// 预定义的错误回复实例。
var (
	RequestEnd = NewErr(ErrEnd)
	SyntaxErr  = NewErr(ErrSyntax)
	HeaderType = NewErr(ErrHeaderType)
)
