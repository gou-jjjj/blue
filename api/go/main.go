package main

import "fmt"

func main() {
	c := NewClient(WithDefaultOpt())

	fmt.Println(c.Select("0"))

	fmt.Println(c.Nget("Port"))
}
