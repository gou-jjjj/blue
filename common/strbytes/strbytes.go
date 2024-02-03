package strbytes

import (
	"fmt"
	"unsafe"
)

func Str2Bytes(s string) []byte {
	bytes := unsafe.StringData(s)
	fmt.Printf("%v\n", *bytes)

	return nil
}
