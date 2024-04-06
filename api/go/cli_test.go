package blue

import (
	"fmt"
	"testing"
	"time"
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
	c, err := NewClient(WithDefaultOpt(), func(config *Config) {
		config.Addr = "127.0.0.1:7890"
	})

	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 1e3; i++ {
		_, err = c.Set(fmt.Sprintf("aaaaaa%d", i), fmt.Sprintf("aaaaaa%d", i))
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Millisecond * 1000)
	}
}
