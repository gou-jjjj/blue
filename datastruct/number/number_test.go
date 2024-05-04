package number

import (
	"blue/common/strbytes"
	"fmt"
	"testing"
)

func TestNumber(t *testing.T) {
	N, _ := NewNumber("512")

	fmt.Printf("%v\n", N.Value())
	fmt.Printf("%v\n", N.GetType())

	fmt.Printf("%v\n", N.Get())
	N.Set(1000)

	fmt.Printf("%v\n", N.Add(24))
	fmt.Printf("%v\n", N.Sub(24))
}

func TestNumber2(t *testing.T) {
	fmt.Printf("%v\n", strbytes.CheckInt("123e31d"))
	fmt.Printf("%v\n", strbytes.CheckInt("31309049189493147837887"))

	by := strbytes.Str2Bytes("4182")

	number, err := NewNumber(by)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%v\n", number.Value())
}
