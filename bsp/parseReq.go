package bsp

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
)

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

		split := bytes.Split(bs[1:len(bs)-1], []byte{Split})

		switch arity {
		case 1:
			if len(split) != 1 {
				return nil, NewErr(ErrSyntax)
			}

			res.SetKey(string(split[0]))

		case 2:
			if len(split) != 2 {
				return nil, NewErr(ErrSyntax)
			}
			res.SetKey(string(split[0]))
			res.SetValue(split[1])

		case -1:
			if len(split) == 0 {
				return nil, NewErr(ErrSyntax)
			}

			res.SetValues(split)

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
