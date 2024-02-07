package bsp

import "blue/common/strbytes"

type RequestBuilder struct {
	data []byte
}

func NewRequestBuilder(handle Header) *RequestBuilder {
	return &RequestBuilder{
		data: NewHeader(handle).Bytes(),
	}
}

func (rb *RequestBuilder) WithKey(key string) *RequestBuilder {
	rb.data = append(rb.data, []byte(key)...)
	return rb
}

func (rb *RequestBuilder) WithValueStr(value string) *RequestBuilder {
	rb.data = append(AppendSplit(rb.data), []byte(value)...)
	return rb
}

func (rb *RequestBuilder) WithValueNum(value string) *RequestBuilder {
	by := strbytes.Str2Bytes(value)
	rb.data = append(AppendSplit(rb.data), by...)
	return rb
}

func (rb *RequestBuilder) WithValues( values ...string) *RequestBuilder {
	for _, value := range values {
		rb.data = append(AppendSplit(rb.data), []byte(value)...)
	}
	return rb
}

func (rb *RequestBuilder) Build() []byte {
	return AppendDone(rb.data)
}