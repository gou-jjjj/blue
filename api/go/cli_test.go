package blue

import (
	"fmt"
	"testing"
)

func TestCli(t *testing.T) {
	c, _ := NewClient(WithDefaultOpt(), func(c *Config) {
		c.DefaultDB = true
		c.DB = 0
	})

	fmt.Println(c.Dbsize())
	fmt.Println(c.Type("DBSum"))
	fmt.Println(c.Nget("DBSum"))

}
