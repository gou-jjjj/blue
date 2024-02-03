package bsp

import (
	"bytes"
	"fmt"
	"testing"
)

func TestBspProtos2(t *testing.T) {
	b := NewBspReq(NumGet, "key")
	buf := bytes.NewBuffer(b)

	protos, err := BspProtos(buf)
	if err != nil {
		t.Fatal(err)
	}

	header := protos[0]

	fmt.Printf("%v\n", header.Handle())
	fmt.Printf("%v\n", header.Type())
	fmt.Printf("%v\n", header.Len())
	fmt.Printf("%v\n", header.TypeStr())
	fmt.Printf("%v\n", header.HandleStr())
	fmt.Printf("%v\n", header.ValueLenMax())
	fmt.Printf("%v\n", header.Bytes())

	fmt.Printf("%v\n", header.Key())
	fmt.Printf("%v\n", header.ValueStr())
}
