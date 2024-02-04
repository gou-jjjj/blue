package bsp

func NewReq(handle Header) []byte {
	return AppendDone(NewHeader(handle, 1).Bytes())
}

func NewReqK(handle Header, key string) []byte {
	b := NewHeader(handle, 2).Bytes()
	b = append(AppendSplit(b), []byte(key)...)
	return AppendDone(b)

}

func NewReqKV(handle Header, key string, value string) []byte {
	b := AppendSplit(NewHeader(handle, 3).Bytes())
	b = AppendSplit(append(b, []byte(key)...))
	return AppendDone(append(b, []byte(value)...))
}

func NewReqKVs(handle Header, key string, value ...string) []byte {
	b := AppendSplit(NewHeader(handle, uint8(2+len(value))).Bytes())
	b = AppendSplit(append(b, []byte(key)...))
	for i := range value {
		b = AppendSplit(append(b, []byte(value[i])...))
	}

	return AppendDone(b[:len(b)-1])
}
