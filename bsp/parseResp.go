package bsp

import (
	"blue/common/strbytes"
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
)

// parse1 从给定的io.Reader中解析数据，将解析成功的消息放入bp通道，错误放入err1通道。
func parse1(r io.Reader, bp chan string, err1 chan error) {
	reader := bufio.NewReader(r)
	for {
		bs, err := reader.ReadBytes(Done) // 尝试从reader读取到Done字符的数据。
		if err != nil || len(bs) < 3 {
			if err != io.EOF {
				err1 <- errors.New("ERR replication error") // 非EOF错误，发送错误信息。
				return
			}
			break // 遇到EOF，结束循环。
		}
		message, err := NewReplyMessage(bs) // 尝试构造回复消息。
		if err != nil {
			err1 <- err // 构造失败，发送错误信息。
			return
		}
		bp <- message // 构造成功，发送消息。
	}
}

// NewReplyMessage 根据给定的回复字节序列构造一个回复消息字符串。
// 返回消息字符串和可能的错误。
func NewReplyMessage(reply []byte) (string, error) {
	if reply == nil || len(reply) < 2 || reply[len(reply)-1] != Done {
		return "", NewErr(ErrReplication) // 验证失败，返回复制错误。
	}

	reply = reply[:len(reply)-1] // 移除末尾的Done字符。

	switch reply[0] { // 根据消息类型处理消息体。
	case byte(ReplyNumber):
		return strbytes.Bytes2Str(reply[1:]), nil
	case byte(ReplyString):
		return string(reply[1:]), nil
	case byte(ReplyList):
		return fmt.Sprintf("%v", reply[1:]), nil
	}

	return MessageMap[reply[0]], nil // 从MessageMap中查找并返回消息。
}

// parseResp 在给定的上下文中解析响应，将解析成功的消息放入bp通道，错误放入err通道。
func parseResp(ctx context.Context, reader io.Reader, bp chan string, err chan error) {
	for {
		select {
		case <-ctx.Done(): // 检查上下文是否取消。
			return
		default:
			parse1(reader, bp, err) // 默认情况下继续解析。
		}
	}
}
