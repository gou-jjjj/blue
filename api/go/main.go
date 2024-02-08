package main

import "fmt"

func main() {
	c := NewClient(WithDefaultOpt())

	nset, err := c.Nset("hello", "789654123")
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

	nset, err = c.Del("hello")
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	nset, err = c.Nget("hello")
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	fmt.Printf("%v\n", nset)
}
