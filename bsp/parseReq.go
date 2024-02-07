package bsp

import (
	"bufio"
	"bytes"
	"context"
	"io"
)

func parse0(r io.Reader) (*BspProto, *ErrResp) {
	reader := bufio.NewReader(r)
	for {
		bs, err := reader.ReadBytes(Done)
		if err != nil || len(bs) < 2 {
			if err != io.EOF {
				return nil, NewErr(ErrSyntax)
			}
			return nil, nil
		}

		res := BspPool.Get().(*BspProto)
		res.Header = NewHeader(Header(bs[0]))
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

			res.key = string(split[0])
			return res, nil
		case 2:

		case -1:

		default:

		}

	}
}

func parseReq(ctx context.Context, reader io.Reader, bp chan *BspProto, err chan *ErrResp) {
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

			err <- errResp
			return
		}
	}
}
