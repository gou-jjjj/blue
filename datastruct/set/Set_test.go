package set

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	s := NewSet()

	s.Add("111")
	s.Add("111")
	s.Add("1111")

	fmt.Println(s.String())
}
