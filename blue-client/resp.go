package main

import (
	"blue/bsp"
	"strconv"
)

func NewReplyMessage(reply []byte) (string, error) {
	if bsp.Split != reply[len(reply)-1] {
		return "", ErrInvalidResp()
	}

	reply = reply[:len(reply)-1]

	s, ok := bsp.MessageMap[bsp.ReplyType(reply[0])]
	if !ok {
		return "", ErrInvalidResp()
	}

	switch s {
	case "number":
		return strconv.FormatUint(common.BytesToUint64(reply[1:]), 10), nil
	case "string", "list":
		return string(reply[1:]), nil
	}

	return s, nil
}
