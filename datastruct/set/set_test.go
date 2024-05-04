package set

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	s := NewSet()
	fmt.Println(s.GetType())
	s.Add("111")
	s.Add("111")
	s.Add("1111")

	fmt.Printf("%v\n", s.GetType())

	fmt.Println(s.String())
}
