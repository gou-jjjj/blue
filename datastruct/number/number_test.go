package number

import (
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
