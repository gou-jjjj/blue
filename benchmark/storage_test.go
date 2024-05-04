package benchmark

import (
	blue "blue/api/go"
	"fmt"
	"strconv"
	"testing"
	"time"
)

var url = "0.0.0.0:13140"

func TestStorage(t *testing.T) {
	now := time.Now()
	defer func() {
		fmt.Println(time.Since(now))
	}()

	c, err := blue.NewClient(blue.WithDefaultOpt(), func(config *blue.Config) {
		config.Addr = url
	})
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	for i := 0; i < 1e3; i++ {
		_, err = c.Set(fmt.Sprintf("str%d", i), fmt.Sprintf("str:%d", i))
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("string done")
	for i := 0; i < 1e3; i++ {
		_, err := c.Lpush("list", fmt.Sprintf("%d", i))
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("list done")

	for i := 0; i < 1e3; i++ {
		_, err := c.Sadd("set", fmt.Sprintf("%d", i))
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("set done")
	for i := 0; i < 1e3; i++ {
		_, err := c.Nset(fmt.Sprintf("num%d", i), strconv.Itoa(i))
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("num done")
}

func TestDel(t *testing.T) {
	now := time.Now()
	defer func() {
		fmt.Println(time.Since(now))
	}()

	c, err := blue.NewClient(blue.WithDefaultOpt(), func(config *blue.Config) {
		config.Addr = url
	})
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	for i := 0; i < 1e3; i++ {
		_, err = c.Del(fmt.Sprintf("str%d", i))
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("string done")

	for i := 0; i < 1e3; i++ {
		_, err := c.Del(fmt.Sprintf("num%d", i))
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("num done")
}
