package blue

import (
	"fmt"
	"testing"
)

func TestCli(t *testing.T) {
	c := NewClient(WithDefaultOpt())
	fmt.Println(c.Select())
	fmt.Println(c.Select("43"))
	fmt.Println(c.Select("3"))
	fmt.Println(c.Select())
	fmt.Println(c.Select("43"))
	fmt.Println(c.Select("3"))
}
