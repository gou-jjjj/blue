package bsp

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"time"
)

const Split byte = '\n'

const (
	bytesTypeMask   = 0b11100000
	bytesHandleMask = 0b00011111
)

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

// todo: 添加handle
func (b *BspProto) String() string {
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

func parseHeader(reader io.Reader, b *BspProto) *ErrResp {
	header := make([]byte, 2)
	n, err := io.ReadFull(reader, header)
	if n == 0 && err == io.EOF {

		fmt.Println(n, err)

		time.Sleep(1000 * time.Second)
	}

	if err != nil {
		return NewErr(ErrSyntax)
	}

	b.Header = NewHeader(Header(header[0]), int8(header[1]))

	if ok := HandleMap[b.Handle()]; ok == "" {
		return NewErr(ErrHeaderType)
	}
	return nil
}

func parseBody(reader io.Reader, b *BspProto) *ErrResp {
	r := bufio.NewReader(reader)

	// read key
	bytes, err := r.ReadBytes(Split)
	if err != nil && err != io.EOF {
		return NewErr(ErrSyntax)
	}

	b.key = string(bytes[:len(bytes)-1])
	if len(b.key) == 0 {
		return NewErr(ErrSyntax)
	}

	// read value
	for i := b.Len(); i > 1; i-- {
		bytes, err = r.ReadBytes(Split)
		if err != nil && err != io.EOF {
			return NewErr(ErrSyntax)
		}

		b.value = append(b.value, bytes[:len(bytes)-1])
	}

	if len(b.value)+1 > int(b.ValueLenMax()) {
		return NewErr(ErrNumberArguments)
	}

	if len(b.value) == 0 {
		b.value = append(b.value, []byte{})
	}

	return nil
}

func parse(ctx context.Context, reader io.Reader, bp chan *BspProto, err chan *ErrResp) {
	defer func() {
		err1 := recover()
		fmt.Printf("re:[%v]\n", err1)
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			b := new(BspProto)
			if errs := parseHeader(reader, b); errs != nil {
				err <- errs
			}

			if errs := parseBody(reader, b); errs != nil {
				err <- errs
			}
			bp <- b
		}
	}
}
