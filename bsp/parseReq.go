package bsp

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
)

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

func parseReq(ctx context.Context, reader io.Reader, bp chan *BspProto, err chan *ErrResp) {
	defer func() {
		if err1 := recover(); err1 != nil {
			fmt.Printf("parseReq err:[%v]\n", err1)
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
