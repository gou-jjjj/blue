package blue

import (
	"fmt"
	"testing"
	"time"
)

func TestCli(t *testing.T) {
	c, _ := NewClient(WithDefaultOpt())

	fmt.Println(c.Lpush("a1", "1"))
	fmt.Println(c.Lpush("a1", "2"))
	fmt.Println(c.Lpush("a1", "3"))

	fmt.Println(c.Lpop("a1"))
	time.Sleep(1 * time.Second)
	fmt.Println(c.Lpop("a1"))
	time.Sleep(1 * time.Second)
	fmt.Println(c.Lpop("a1"))
}
