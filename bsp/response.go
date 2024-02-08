package bsp

import (
	"blue/common/strbytes"
	"fmt"
	"strconv"
)

type Reply interface {
	Bytes() []byte
	String() string
}

// InfoReply ----------------------------------
type InfoReply struct {
	i    ReplyType
	info []byte
}

func (i InfoReply) Bytes() []byte {
	res := make([]byte, 0, len(i.info)+2)
	res = append(res, byte(i.i))
	res = append(res, i.info...)
	return AppendDone(res)
}

func (i InfoReply) String() string {
	if i.info == nil || len(i.info) == 0 {
		return MessageMap[i.i]
	} else {
		return fmt.Sprintf("%s:%s", MessageMap[i.i], i.info)
	}
}

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

// NumResp -------------------------------------
type NumResp struct {
	n   ReplyType
	num []byte
}

func (n NumResp) Bytes() []byte {
	res := make([]byte, 0, len(n.num)+2)
	res = append(res, byte(n.n))
	res = append(res, n.num...)
	return AppendDone(res)
}

func (n NumResp) String() string {
	return strbytes.Bytes2Str(n.num)
}

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

// StrResp -------------------------------------
type StrResp struct {
	s   ReplyType
	msg []byte
}

func (s StrResp) Bytes() []byte {
	res := make([]byte, 0, len(s.msg)+2)
	res = append(res, byte(s.s))
	res = append(res, s.msg...)
	return AppendDone(res)
}

func (s StrResp) String() string {
	return string(s.msg)
}

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

// ListResp -------------------------------------
type ListResp struct {
	l    ReplyType
	list [][]byte
}

func (l ListResp) Bytes() []byte {
	b := make([]byte, 0, len(l.list))
	b = append(b, byte(l.l))
	for i := range l.list {
		b = append(b, AppendSplit(l.list[i])...)
	}

	return AppendDone(b)
}

func (l ListResp) String() string {
	return fmt.Sprintf("%v", l.list)
}

func NewList(list ...[]byte) *ListResp {
	return &ListResp{
		l:    ReplyList,
		list: list,
	}
}

// ErrResp -------------------------------------
type ErrResp struct {
	e   ReplyType
	msg []byte
}

func (r ErrResp) Error() string {
	return r.String()
}

func (r ErrResp) Bytes() []byte {
	res := make([]byte, 0, len(r.msg)+2)
	res = append(res, byte(r.e))
	res = append(res, r.msg...)
	return AppendDone(res)
}

func (r ErrResp) String() string {
	if r.msg == nil || len(r.msg) == 0 {
		return MessageMap[r.e]
	}
	return fmt.Sprintf("%s:%s", MessageMap[r.e], r.msg)
}

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

var (
	RequestEnd = NewErr(ErrEnd)
	SyntaxErr  = NewErr(ErrSyntax)
	HeaderType = NewErr(ErrHeaderType)
)
