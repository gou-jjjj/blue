package main

import "fmt"

func main() {
	c := NewClient(WithDefaultOpt())
	defer c.Close()
	nset, err := c.Nset("hello", "123")
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	fmt.Printf("%v\n", nset)

	nset, err = c.Nget("hello")
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	fmt.Printf("%v\n", nset)
}
