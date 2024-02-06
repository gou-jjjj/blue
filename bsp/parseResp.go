package bsp

import (
	"blue/common/strbytes"
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
)

func parse1(r io.Reader, bp chan string, err1 chan error) {
	reader := bufio.NewReader(r)
	for {
		bs, err := reader.ReadBytes(Done)
		if err != nil || len(bs) < 3 {
			if err != io.EOF {
				err1 <- errors.New("ERR replication error")
				return
			}
			break
		}
		message, err := NewReplyMessage(bs)
		if err != nil {
			err1 <- err
			return
		}
		bp <- message
	}
}

func NewReplyMessage(reply []byte) (string, error) {
	if reply == nil || len(reply) < 2 || reply[len(reply)-1] != Done {
		return "", NewErr(ErrReplication)
	}

	reply = reply[:len(reply)-1]

	switch reply[0] {
	case byte(ReplyNumber):
		return strbytes.Bytes2Str(reply[1:]), nil
	case byte(ReplyString):
		return string(reply[1:]), nil
	case byte(ReplyList):
		return fmt.Sprintf("%v", reply[1:]), nil
	}

	return MessageMap[reply[0]], nil
}

func parseResp(ctx context.Context, reader io.Reader, bp chan string, err chan error) {
	defer func() {
		if err1 := recover(); err1 != nil {
			fmt.Printf("parseResp err:[%v]\n", err1)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			parse1(reader, bp, err)
		}
	}
}
