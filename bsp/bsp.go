package bsp

import (
	"bufio"
	"bytes"
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
	go parse(ctx, r, protos, errs)
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

func parse0(r io.Reader, bp chan *BspProto, err1 chan *ErrResp) {
	reader := bufio.NewReader(r)
	for {
		bs, err := reader.ReadBytes(Done)
		if err != nil || len(bs) < 3 {
			if err != io.EOF {
				err1 <- NewErr(ErrSyntax)
				return
			}
			break
		}

		res := bspPool.Get().(*BspProto)
		res.Header = NewHeader(Header(bs[0]), bs[1])
		if res.Header == HandleErr {
			err1 <- NewErr(ErrHeaderType)
			return
		}

		arity := res.Header.HandleInfo().Arity
		if arity == 0 {
			if bs[2] != Done {
				err1 <- NewErr(ErrSyntax)
				return
			}
			bp <- res
			break
		}

		if arity > 0 && bs[2] != Split {
			err1 <- NewErr(ErrSyntax)
			return
		}

		split := bytes.Split(bs[3:], []byte{Split})

		switch arity {
		case 1:
			if len(split) != 1 || split[0][len(split[0])-1] != Done {
				err1 <- NewErr(ErrSyntax)
				return
			}

			res.key = string(split[0][:len(split[0])-1])
			bp <- res
			break
		case 2:
			if len(split) != 2 || split[1][len(split[1])-1] != Done {
				err1 <- NewErr(ErrSyntax)
				return
			}
			res.key = string(split[0])
			res.value = [][]byte{split[1][:len(split[1])-1]}
			bp <- res
			break
		case -1:

		default:
			err1 <- NewErr(ErrSyntax)
			return
		}

	}
}

func parse(ctx context.Context, reader io.Reader, bp chan *BspProto, err chan *ErrResp) {
	defer func() {
		if err1 := recover(); err1 != nil {
			fmt.Printf("parse err:[%v]\n", err1)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			parse0(reader, bp, err)
		}
	}
}
