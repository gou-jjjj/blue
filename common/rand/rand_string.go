package rand

import (
	"bytes"
	"math/rand"
)

const (
	randStr    = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	randStrLen = len(randStr)
)

func RandString(l int) string {
	if l <= 0 {
		l = 8
	}

	buf := bytes.Buffer{}
	for i := 0; i < l; i++ {
		buf.WriteByte(randStr[rand.Intn(randStrLen)])
	}
	return buf.String()
}
