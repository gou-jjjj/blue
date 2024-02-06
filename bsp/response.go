package bsp

import (
	"blue/common/strbytes"
	"fmt"
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
	if i.info != nil {
		return SufSplit(append([]byte{byte(i.i)}, i.info...))
	}

	return SufSplit([]byte{byte(i.i)})
}

func (i InfoReply) String() string {
	if i.info == nil || len(i.info) == 0 {
		return ReplyTypeMap[i.i]
	} else {
		return fmt.Sprintf("%s:%s", ReplyTypeMap[i.i], i.info)
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
	return SufSplit(append([]byte{byte(ReplyNumber)}, n.num...))
}

func (n NumResp) String() string {
	return strbytes.Bytes2Str(n.num)
}

func NewNum(num any) *NumResp {
	switch num.(type) {
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
		return &NumResp{
			n:   ReplyNumber,
			num: []byte{byte(0)},
		}
	}
}

// StrResp -------------------------------------
type StrResp struct {
	s   ReplyType
	msg []byte
}

func (s StrResp) Bytes() []byte {
	return SufSplit(append([]byte{byte(ReplyString)}, s.msg...))
}

func (s StrResp) String() string {
	return string(s.msg)
}

func NewStr(msg []byte) *StrResp {
	return &StrResp{
		s:   ReplyString,
		msg: msg,
	}
}

// ListResp -------------------------------------
type ListResp struct {
	l    ReplyType
	list [][]byte
}

func (l ListResp) Bytes() []byte {
	b := make([]byte, 0, len(l.list))
	for i := range l.list {
		b = append(b, SufSplit(l.list[i])...)
	}

	return append(append([]byte{byte(l.l)}, byte(len(l.list))), b...)
}

func (l ListResp) String() string {
	return fmt.Sprintf("%s", l.list)
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
	return string(r.msg)
}

func (r ErrResp) Bytes() []byte {
	return SufSplit(append([]byte{byte(r.e)}, r.msg...))
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

func SufSplit(b []byte) []byte {
	return append(b, Split)
}

func PreSplit(b []byte) []byte {
	return append([]byte{Split}, b...)
}
