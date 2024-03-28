package blue

import (
	"fmt"
	"testing"
)

func TestCli(t *testing.T) {
	c, _ := NewClient(WithDefaultOpt())
	fmt.Println(c.Sadd("a", "1"))
	fmt.Println(c.Sadd("a", "2"))
	fmt.Println(c.Sadd("a", "3"))
	fmt.Println(c.Sadd("a", "1"))

	fmt.Println(c.Sget("a"))
}
