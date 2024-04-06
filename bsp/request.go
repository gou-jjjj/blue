// Package bsp 提供了构建请求数据的工具。
package bsp

import "blue/common/strbytes"

// RequestBuilder 是用于构建请求数据的结构体。
type RequestBuilder struct {
	data []byte // 存储构建的请求数据
}

// NewRequestBuilder 创建一个新的RequestBuilder实例。
// handle: 请求头的初始信息。
// 返回一个新的RequestBuilder实例，其数据部分初始化为给定请求头的字节表示。
func NewRequestBuilder(handle Header) *RequestBuilder {
	return &RequestBuilder{
		data: NewHeader(handle).Bytes(),
	}
}

// WithKey 向请求数据中添加一个键。
// key: 要添加的键。
// 返回修改后的RequestBuilder实例，以便进行链式调用。
func (rb *RequestBuilder) WithKey(key string) *RequestBuilder {
	rb.data = append(rb.data, []byte(key)...)
	return rb
}

// WithKeyNum 向请求数据中添加一个以数字形式表示的键。
// key: 要添加的键，字符串形式。
// 返回修改后的RequestBuilder实例，以便进行链式调用。
func (rb *RequestBuilder) WithKeyNum(key string) *RequestBuilder {
	by := strbytes.Str2Bytes(key)
	rb.data = append(rb.data, by...)
	return rb
}

// WithValueStr 向请求数据中添加一个字符串类型的值。
// value: 要添加的值。
// 返回修改后的RequestBuilder实例，以便进行链式调用。
func (rb *RequestBuilder) WithValueStr(value string) *RequestBuilder {
	rb.data = append(AppendSplit(rb.data), []byte(value)...)
	return rb
}

// WithValueNum 向请求数据中添加一个数值类型的值。
// value: 要添加的值，字符串形式。
// 返回修改后的RequestBuilder实例，以便进行链式调用。
func (rb *RequestBuilder) WithValueNum(value string) *RequestBuilder {
	by := strbytes.Str2Bytes(value)
	rb.data = append(AppendSplit(rb.data), by...)
	return rb
}

// WithValues 向请求数据中批量添加多个值。
// values: 要添加的值的切片。
// 返回修改后的RequestBuilder实例，以便进行链式调用。
func (rb *RequestBuilder) WithValues(values ...string) *RequestBuilder {
	for _, value := range values {
		rb.data = append(AppendSplit(rb.data), []byte(value)...)
	}
	return rb
}

// Build 构建并返回最终的请求数据。
// 返回构建完成的请求数据的字节切片。
func (rb *RequestBuilder) Build() []byte {
	return AppendDone(rb.data)
}
