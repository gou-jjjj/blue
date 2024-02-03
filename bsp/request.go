package bsp

import "bytes"

func NewBspDataReq(handle Header, key string, value ...[]byte) []byte {
	buf := bytes.Buffer{}

	header := NewHeader(handle, int8(len(value))+1)
	buf.Write(header.Bytes())
	buf.WriteString(key)
	buf.WriteByte(Split)

	for i := range value {
		buf.Write(value[i])
		buf.WriteByte(Split)
	}

	return buf.Bytes()
}

func NewBspDbReq() {

}

func NewBspSysReq(handle Header) []byte {
	return NewHeader(handle, 1).Bytes()
}
