package bsp

import (
	"fmt"
	"testing"
)

func TestHeader(t *testing.T) {
	header := NewHeader(StrSet, 2)

	fmt.Printf("%v\n", header.Handle())
	fmt.Printf("%v\n", header.Type())
	fmt.Printf("%v\n", header.Len())
	fmt.Printf("%v\n", header.TypeStr())
	fmt.Printf("%v\n", header.HandleStr())
	fmt.Printf("%v\n", header.ValueLenMax())
	fmt.Printf("%v\n", header.Bytes())
}
