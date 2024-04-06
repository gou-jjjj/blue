package bsp

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"strings"
)

// 定义分隔符常量
var split = []byte{Split}

// parse0 用于解析输入流中的BSP协议数据，并返回解析结果或错误。
func parse0(r io.Reader) (*BspProto, *ErrResp) {
	reader := bufio.NewReader(r)
	for {
		// 读取直到遇到指定的结束标识
		bs, err := reader.ReadBytes(Done)
		if err != nil {
			// 判断错误类型，返回相应的错误响应
			if normalErr(err) {
				return nil, RequestEnd
			}

			return nil, SyntaxErr
		}

		// 检查读取的数据长度是否满足最小要求
		if len(bs) < 2 {
			return nil, SyntaxErr
		}

		// 初始化BspProto对象，并设置读取到的数据
		res := NewBspProto()
		res.SetBuf(bs)
		res.SetHeader(NewHeader(Header(bs[0])))
		// 检查协议头是否错误
		if res.Header == HandleErr {
			return nil, NewErr(ErrHeaderType)
		}

		// 根据协议头的处理信息，解析数据
		arity := res.Header.HandleInfo().Arity
		if arity == 0 {
			// 无参数情况的处理
			if bs[1] != Done {
				return nil, NewErr(ErrSyntax)
			}
			return res, nil
		}

		// 使用指定的分隔符分割数据
		bss := bytes.Split(bs[1:len(bs)-1], split)

		// 根据参数个数设置协议内容
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
			// 对于不支持的参数个数，不做处理
		}
		return res, nil
	}
}

// parseReq 在给定的上下文中异步解析请求，将解析成功的协议对象发送至bp通道，错误发送至err通道。
func parseReq(ctx context.Context, reader io.Reader, bp chan *BspProto, err chan *ErrResp) {
	defer func() {
		// 捕获并忽略panic，保证函数正常退出
		_ = recover()
	}()

	for {
		select {
		case <-ctx.Done():
			// 当上下文被取消或超时时退出
			return
		default:
			// 默认情况下继续解析协议
			proto, errResp := parse0(reader)
			if proto != nil {
				// 解析成功，发送协议对象
				bp <- proto
				break
			}
			// 解析失败，发送错误响应
			err <- errResp
			return
		}
	}
}

// normalErr 用于判断错误是否为常见的预期错误，如连接被远程主机强制关闭或EOF。
func normalErr(err error) bool {
	// 判断错误信息中是否包含特定字符串
	if strings.Contains(err.Error(), "An existing connection was forcibly closed by the remote host.") {
		return true
	}

	if strings.Contains(err.Error(), "EOF") {
		return true
	}
	return false
}
