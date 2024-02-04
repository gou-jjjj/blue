package rand

import (
	"bytes"
)

const (
	randStr    = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	randStrLen = uint32(len(randStr))
)

func RandString(l int) string {
	if l <= 0 {
		l = 8
	}

	buf := bytes.Buffer{}
	for i := 0; i < l; i++ {
		buf.WriteByte(randStr[Randu32()%randStrLen])
	}
	return buf.String()
}
