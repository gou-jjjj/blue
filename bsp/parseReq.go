package bsp

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
)

var split = []byte{Split}

func parse0(r io.Reader) (*BspProto, *ErrResp) {
	reader := bufio.NewReader(r)
	for {
		bs, err := reader.ReadBytes(Done)
		if err != nil {
			if normalErr(err) {
				return nil, RequestEnd
			}

			return nil, SyntaxErr
		}

		if len(bs) < 2 {
			return nil, SyntaxErr
		}
		fmt.Printf("bs: {%+b}\n", bs)
		res := NewBspProto()
		res.SetBuf(bs)
		res.SetHeader(NewHeader(Header(bs[0])))
		if res.Header == HandleErr {
			return nil, NewErr(ErrHeaderType)
		}

		arity := res.Header.HandleInfo().Arity
		if arity == 0 {
			if bs[1] != Done {
				return nil, NewErr(ErrSyntax)
			}
			return res, nil
		}

		bss := bytes.Split(bs[1:len(bs)-1], split)

		switch arity {
		case 1:
			if len(bss) != 1 {
				return nil, NewErr(ErrSyntax)
			}

			res.SetKey(string(bss[0]))

		case 2:
			if len(bss) != 2 {
				return nil, NewErr(ErrSyntax)
			}
			res.SetKey(string(bss[0]))
			res.SetValue(bss[1])

		case -1:
			if len(bss) == 0 {
				return nil, NewErr(ErrSyntax)
			}

			res.SetValues(bss)

		default:

		}
		return res, nil
	}
}

func parseReq(ctx context.Context, reader io.Reader, bp chan *BspProto, err chan *ErrResp) {
	defer func() {
		if er := recover(); er != nil {
			fmt.Printf("reco [%v]\n", er)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			proto, errResp := parse0(reader)
			if proto != nil {
				bp <- proto
				break
			}
			fmt.Printf("errResp [%v]\n", errResp)
			err <- errResp
			return
		}
	}
}

func normalErr(err error) bool {
	if strings.Contains(err.Error(), "An existing connection was forcibly closed by the remote host.") {
		return true
	}
	return false
}
