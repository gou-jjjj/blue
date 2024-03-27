package blue

import (
	"fmt"
	"testing"
)

func TestCli(t *testing.T) {
	c, _ := NewClient(WithDefaultOpt())
	fmt.Println(c.Select("0"))
	kvs, err := c.Kvs()
	if err != nil {
		panic(err)
	}

	fmt.Println(kvs)
}
