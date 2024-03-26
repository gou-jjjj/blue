package blue

import (
	"fmt"
	"testing"
)

func TestCli(t *testing.T) {
	c, _ := NewClient(WithDefaultOpt())

	fmt.Println(c.Get("a"))
	fmt.Println(c.Set("a", "b"))
}
