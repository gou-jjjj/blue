package blue

import (
	"fmt"
	"testing"
)

func handleErr(s string, err error) {
	if err != nil {
		panic(fmt.Sprintf("%s:%s", s, err))
	}
}

func TestAllCli(t *testing.T) {
	c, err := NewClient(WithDefaultOpt())
	if err != nil {
		t.Fatal(err)
	}

	handleErr(c.Set("a", "1"))
	handleErr(c.Set("b", "2"))
	handleErr(c.Set("c", "3"))
	handleErr(c.Get("c"))
	handleErr(c.Del("c"))

	handleErr(c.Expire("a", "10"))
	handleErr(c.Expire("b", "20"))
	handleErr(c.Type("b"))
	c.Close()
}

func TestCli(t *testing.T) {
	c, _ := NewClient(WithDefaultOpt(), func(c *Config) {
		c.DefaultDB = true
		c.DB = 0
	})

	fmt.Println(c.Dbsize())
	fmt.Println(c.Type("DBSum"))
	fmt.Println(c.Nget("DBSum"))
}
