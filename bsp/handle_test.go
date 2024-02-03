package bsp

import (
	"fmt"
	"testing"
)

func TestHeader_Handle(t *testing.T) {
	fmt.Printf("%#v\n", HandleMap[HandleMap2["str set"]])
	fmt.Printf("%#v\n", HandleMap[HandleMap2["num set"]])
	fmt.Printf("%#v\n", HandleMap[HandleMap2["num get"]])
	fmt.Printf("%#v\n", HandleMap[HandleMap2["str get"]])
}
